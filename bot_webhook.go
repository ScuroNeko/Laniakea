package laniakea

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"
)

// BotWebhookOpts configures Telegram webhook registration and the local HTTP server.
type BotWebhookOpts struct {
	Path          string
	LocalPort     int
	UseStatusPath bool

	URL                string
	Certificate        []byte
	IPAddress          string
	MaxConnections     int8
	AllowedUpdates     []tgapi.UpdateType
	DropPendingUpdates bool
	SecretToken        string
}

// NewBotWebhookOpts returns webhook options with the default path, local port, and max connections.
func NewBotWebhookOpts() *BotWebhookOpts {
	return &BotWebhookOpts{
		Path:           "/",
		LocalPort:      8080,
		MaxConnections: 40,
	}
}

// SetPath sets the local HTTP path that receives Telegram webhook requests.
func (opts *BotWebhookOpts) SetPath(path string) *BotWebhookOpts {
	opts.Path = path
	return opts
}

// SetLocalPort sets the local HTTP port used by the webhook server.
func (opts *BotWebhookOpts) SetLocalPort(port int) *BotWebhookOpts {
	opts.LocalPort = port
	return opts
}

// SetUseStatusPath enables or disables the optional /status endpoint.
// A non-empty SecretToken is required when this endpoint is enabled.
func (opts *BotWebhookOpts) SetUseStatusPath(use bool) *BotWebhookOpts {
	opts.UseStatusPath = use
	return opts
}

// SetURL sets the public base URL Telegram should call for incoming updates.
func (opts *BotWebhookOpts) SetURL(url string) *BotWebhookOpts {
	opts.URL = url
	return opts
}

// SetCertificate sets the self-signed webhook certificate bytes to upload.
func (opts *BotWebhookOpts) SetCertificate(certificate []byte) *BotWebhookOpts {
	opts.Certificate = certificate
	return opts
}

// MustLoadCertificate loads a webhook certificate from disk and panics on failure.
func (opts *BotWebhookOpts) MustLoadCertificate(filename string) *BotWebhookOpts {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	opts.Certificate, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return opts
}

// SetIPAddress sets the fixed IP address Telegram should use for webhook delivery.
func (opts *BotWebhookOpts) SetIPAddress(ip string) *BotWebhookOpts {
	opts.IPAddress = ip
	return opts
}

// SetMaxConnections sets Telegram's maximum number of simultaneous webhook connections.
func (opts *BotWebhookOpts) SetMaxConnections(max int8) *BotWebhookOpts {
	opts.MaxConnections = max
	return opts
}

// SetAllowedUpdates sets the Telegram update types that should be delivered to the webhook.
func (opts *BotWebhookOpts) SetAllowedUpdates(updates ...tgapi.UpdateType) *BotWebhookOpts {
	opts.AllowedUpdates = append([]tgapi.UpdateType(nil), updates...)
	return opts
}

// SetDropPendingUpdates configures whether Telegram should drop pending updates while setting the webhook.
func (opts *BotWebhookOpts) SetDropPendingUpdates(drop bool) *BotWebhookOpts {
	opts.DropPendingUpdates = drop
	return opts
}

// SetSecretToken sets the secret token expected in Telegram webhook requests.
// The same token is also required to access /status when that endpoint is enabled.
func (opts *BotWebhookOpts) SetSecretToken(secretToken string) *BotWebhookOpts {
	opts.SecretToken = secretToken
	return opts
}

// RunWebhookWithContext registers the webhook and serves incoming updates until ctx is canceled.
//
// The bot uses the same update queue, worker pool, runner startup, and single-use lifecycle
// guarantees as RunWithContext. When opts.AllowedUpdates is empty, the bot-level update types
// configured through SetUpdateTypes/AddUpdateType are used. When UseStatusPath is enabled,
// SecretToken must be non-empty so the operational endpoint is not left public.
//
// When two TLS files are provided, the method serves HTTPS using the existing key-then-cert
// argument order.
func (bot *Bot[T]) RunWebhookWithContext(ctx context.Context, opts *BotWebhookOpts, tlsFiles ...string) error {
	if opts == nil {
		return errors.New("nil BotWebhookOpts")
	}
	if len(bot.prefixes) == 0 {
		return ErrNoPrefixes
	}
	if len(bot.plugins) == 0 {
		return ErrNoPlugins
	}
	if opts.URL == "" {
		return errors.New("empty BotWebhookOpts.URL")
	}
	if opts.MaxConnections > 100 || opts.MaxConnections <= 0 {
		return errors.New("BotWebhookOpts.MaxConnections must between 1 and 100")
	}
	if err := validateWebhookPath(opts.Path, opts.UseStatusPath); err != nil {
		return err
	}
	if opts.UseStatusPath && opts.SecretToken == "" {
		return errors.New("BotWebhookOpts.SecretToken required when status path is enabled")
	}
	if err := validateWebhookTLSFiles(tlsFiles); err != nil {
		return err
	}

	if opts.Certificate != nil && bot.uploader == nil {
		return errors.New("bot uploader nil, but certificate set")
	}

	return bot.runWebhookRuntime(ctx, func(runCtx context.Context) error {
		if opts.SecretToken == "" {
			bot.webhookLogger.Warnln("Bot webhook secret token empty. It's VERY recommended to set secret.")
		}

		i, err := bot.api.GetWebhookInfoWithContext(runCtx)
		if err != nil {
			return err
		}
		if i.URL == "" {
			bot.webhookLogger.Warnln("API returned webhook info with empty URL. There may be a long-poll")
		} else {
			_, err = bot.api.DeleteWebhookWithContext(runCtx, tgapi.DeleteWebhook{})
			if err != nil {
				return err
			}
			bot.webhookLogger.Infof("Bot webhook deleted: %s", i.URL)
		}

		allowedUpdates := bot.webhookAllowedUpdates(opts)

		var ok bool
		if opts.Certificate != nil {
			ok, err = bot.uploader.SetWebhookWithContext(runCtx, tgapi.UploadSetWebhook{
				URL:                fmt.Sprintf("%s%s", opts.URL, opts.Path),
				IPAddress:          opts.IPAddress,
				MaxConnections:     opts.MaxConnections,
				AllowedUpdates:     allowedUpdates,
				DropPendingUpdates: opts.DropPendingUpdates,
				SecretToken:        opts.SecretToken,
			}, tgapi.NewUploaderFile("certificate", opts.Certificate))
		} else {
			ok, err = bot.api.SetWebhookWithContext(runCtx, tgapi.SetWebhook{
				URL:                fmt.Sprintf("%s%s", opts.URL, opts.Path),
				IPAddress:          opts.IPAddress,
				MaxConnections:     opts.MaxConnections,
				AllowedUpdates:     allowedUpdates,
				DropPendingUpdates: opts.DropPendingUpdates,
				SecretToken:        opts.SecretToken,
			})
		}
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("failed to set webhook")
		}

		if len(tlsFiles) == 2 {
			return bot.runWebhookTLS(runCtx, opts, tlsFiles[0], tlsFiles[1])
		}

		return bot.runWebhook(runCtx, opts)
	})
}

// RunWebhook starts the webhook runtime with a background context.
//
// It is shorthand for RunWebhookWithContext(context.Background(), opts, tlsFiles...).
func (bot *Bot[T]) RunWebhook(opts *BotWebhookOpts, tlsFiles ...string) error {
	return bot.RunWebhookWithContext(context.Background(), opts, tlsFiles...)
}

// CloseWebhook removes the current Telegram webhook registration.
//
// It is separate from Close, which only releases local resources.
// Call it before switching a deployment from webhook delivery to polling.
func (bot *Bot[T]) CloseWebhook() error {
	var e []error
	if bot.api == nil {
		e = append(e, errors.New("bot api nil"))
	} else {
		if _, err := bot.api.DeleteWebhook(tgapi.DeleteWebhook{}); err != nil {
			if bot.webhookLogger != nil {
				bot.webhookLogger.Errorf("Failed to close webhook: %s", err.Error())
			} else if bot.logger != nil {
				bot.logger.Errorf("Failed to close webhook: %s", err.Error())
			}
			e = append(e, err)
		}
	}
	if bot.webhookLogger != nil {
		if err := bot.webhookLogger.Close(); err != nil {
			e = append(e, err)
		}
		bot.webhookLogger = nil
	}
	return errors.Join(e...)
}

func (bot *Bot[T]) webhookAllowedUpdates(opts *BotWebhookOpts) []tgapi.UpdateType {
	if len(opts.AllowedUpdates) > 0 {
		return append([]tgapi.UpdateType(nil), opts.AllowedUpdates...)
	}
	return bot.GetUpdateTypes()
}

func (bot *Bot[T]) runWebhookRuntime(ctx context.Context, run func(context.Context) error) error {
	if err := bot.beginRun(); err != nil {
		return err
	}
	defer bot.finishRun()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if bot.webhookLogger == nil {
		bot.webhookLogger = utils.CreateLogger("WEBHOOK", bot.GetLoggerLevel(), bot.logFormat, bot.logFormatter)
	}
	bot.addTokenReplacer(bot.webhookLogger)
	bot.ExecRunners(runCtx)

	workersDone := make(chan struct{})
	go func() {
		bot.startUpdateWorkers(runCtx)
		close(workersDone)
	}()

	runErr := run(runCtx)
	cancel()
	close(bot.updateQueue)
	<-workersDone
	bot.runnerOnceWG.Wait()
	bot.runnerBgWG.Wait()

	return runErr
}

func updateHandler[T any](ctx context.Context, bot *Bot[T], secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if secret != "" && r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != secret {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		const maxWebhookBodySize = 256 << 10 // 256 KiB
		r.Body = http.MaxBytesReader(w, r.Body, maxWebhookBodySize)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(data) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var up tgapi.Update
		if err := json.Unmarshal(data, &up); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			bot.webhookLogger.Errorln(err)
			return
		}
		bot.webhookLogger.Debugf("UPDATE id=%d type=%s size=%d from=%s", up.UpdateID, up.Type, len(data), r.RemoteAddr)
		if err := bot.enqueueUpdate(ctx, up); err != nil {
			bot.webhookLogger.Errorln(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func statusHandler[T any](bot *Bot[T], opts *BotWebhookOpts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := ""
		if r.Header.Get("Authorization") != "" {
			auth = r.Header.Get("Authorization")
		} else if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != "" {
			auth = r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		}
		if auth != opts.SecretToken {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		i, err := bot.api.GetWebhookInfoWithContext(r.Context())
		if err != nil {
			bot.webhookLogger.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			bot.webhookLogger.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(w, string(data)); err != nil {
			bot.webhookLogger.Errorln(err)
		}
	}
}

func (bot *Bot[T]) newWebhookMux(ctx context.Context, opts *BotWebhookOpts) *http.ServeMux {
	r := http.NewServeMux()
	if opts.UseStatusPath {
		r.HandleFunc("/status", statusHandler(bot, opts))
	}
	r.HandleFunc(opts.Path, updateHandler(ctx, bot, opts.SecretToken))
	return r
}
func (bot *Bot[T]) runWebhook(ctx context.Context, opts *BotWebhookOpts) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.LocalPort),
		Handler: bot.newWebhookMux(ctx, opts),
	}
	errCh := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	bot.webhookLogger.Infoln(fmt.Sprintf("Bot Webhook started at %s; waiting for updates at %s", srv.Addr, opts.URL))

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return <-errCh

	case err := <-errCh:
		return err
	}
}
func (bot *Bot[T]) runWebhookTLS(ctx context.Context, opts *BotWebhookOpts, key, cert string) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.LocalPort),
		Handler: bot.newWebhookMux(ctx, opts),
	}
	errCh := make(chan error, 1)

	go func() {
		err := srv.ListenAndServeTLS(cert, key)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	bot.webhookLogger.Infoln(fmt.Sprintf("Bot webhook started with TLS(%s, %s) at %s; waiting for updates at %s", key, cert, srv.Addr, opts.URL))

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return <-errCh

	case err := <-errCh:
		return err
	}
}
func validateWebhookPath(path string, useStatusPath bool) error {
	if path == "" {
		return errors.New("empty BotWebhookOpts.Path")
	}
	if !strings.HasPrefix(path, "/") {
		return errors.New("BotWebhookOpts.Path must start with '/'")
	}
	if strings.Contains(path, "?") || strings.Contains(path, "#") {
		return errors.New("BotWebhookOpts.Path must not contain query or fragment")
	}
	if useStatusPath && path == "/status" {
		return errors.New("BotWebhookOpts.Path must not be '/status' when status path is enabled")
	}
	return nil
}

func validateWebhookTLSFiles(tlsFiles []string) error {
	switch len(tlsFiles) {
	case 0, 2:
		return nil
	case 1:
		return errors.New("you must specify both private and public keys")
	default:
		return errors.New("too many files; you must specify only private and public keys")
	}
}
