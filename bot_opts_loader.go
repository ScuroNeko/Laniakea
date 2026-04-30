package laniakea

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"
)

// ConfigVersion is the current version of the built-in JSON BotOpts file format.
const ConfigVersion = 1

// ErrConfigVersionMismatch reports that a config file declares a newer version
// than this library knows how to decode.
var ErrConfigVersionMismatch = fmt.Errorf("config version mismatch: expected %d", ConfigVersion)

type botOptsFileJSONLogger struct {
	LoggerBasePath   string          `json:"base_path"`
	UseRequestLogger bool            `json:"use_request_logger"`
	WriteToFile      bool            `json:"write_to_file"`
	LogFormat        utils.LogFormat `json:"log_format"`
}
type botOptsFileJSONAPI struct {
	UseTestServer  bool   `json:"use_test_server"`
	APIURL         string `json:"url"`
	RateLimit      int    `json:"rate_limit"`
	DropRLOverflow bool   `json:"drop_overflow"`
}

// BotOptsFileJSON is the JSON file representation of BotOpts.
type BotOptsFileJSON struct {
	Version           int                   `json:"version"`
	Token             string                `json:"token"`
	UpdateTypes       []tgapi.UpdateType    `json:"update_types"`
	Debug             bool                  `json:"debug"`
	ErrorTemplate     string                `json:"error_template"`
	Prefixes          []string              `json:"prefixes"`
	Logger            botOptsFileJSONLogger `json:"logger"`
	API               botOptsFileJSONAPI    `json:"api"`
	StrictPayloadType bool                  `json:"strict_payload_type"`
	MaxWorkers        int                   `json:"max_workers"`
}

// BotOptsFileJSONCodec encodes and decodes BotOpts using BotOptsFileJSON.
type BotOptsFileJSONCodec struct{}

// FromBytes decodes BotOpts from JSON file bytes.
func (codec BotOptsFileJSONCodec) FromBytes(data []byte) (*BotOpts, error) {
	fileOpts := new(BotOptsFileJSON)
	err := json.Unmarshal(data, fileOpts)
	if err != nil {
		return nil, err
	}
	if fileOpts.Version > ConfigVersion {
		return nil, ErrConfigVersionMismatch
	}
	opts := &BotOpts{
		Token:         fileOpts.Token,
		UpdateTypes:   fileOpts.UpdateTypes,
		Debug:         fileOpts.Debug,
		ErrorTemplate: fileOpts.ErrorTemplate,
		Prefixes:      fileOpts.Prefixes,

		LoggerBasePath:   fileOpts.Logger.LoggerBasePath,
		UseRequestLogger: fileOpts.Logger.UseRequestLogger,
		WriteToFile:      fileOpts.Logger.WriteToFile,
		LogFormat:        fileOpts.Logger.LogFormat,

		UseTestServer:         fileOpts.API.UseTestServer,
		APIURL:                fileOpts.API.APIURL,
		RateLimit:             fileOpts.API.RateLimit,
		DropRateLimitOverflow: fileOpts.API.DropRLOverflow,

		StrictPayloadType: fileOpts.StrictPayloadType,
		MaxWorkers:        fileOpts.MaxWorkers,

		FileConfigVersion: fileOpts.Version,
	}
	return opts, nil
}

// ToBytes encodes BotOpts into JSON file bytes.
func (codec BotOptsFileJSONCodec) ToBytes(opts *BotOpts) ([]byte, error) {
	fileOpts := &BotOptsFileJSON{
		Version:       ConfigVersion,
		Token:         opts.Token,
		UpdateTypes:   opts.UpdateTypes,
		Debug:         opts.Debug,
		ErrorTemplate: opts.ErrorTemplate,
		Prefixes:      opts.Prefixes,
		Logger: botOptsFileJSONLogger{
			LoggerBasePath:   opts.LoggerBasePath,
			UseRequestLogger: opts.UseRequestLogger,
			WriteToFile:      opts.WriteToFile,
			LogFormat:        opts.LogFormat,
		},
		API: botOptsFileJSONAPI{
			UseTestServer:  opts.UseTestServer,
			APIURL:         opts.APIURL,
			RateLimit:      opts.RateLimit,
			DropRLOverflow: opts.DropRateLimitOverflow,
		},
		StrictPayloadType: opts.StrictPayloadType,
		MaxWorkers:        opts.MaxWorkers,
	}
	data, err := json.Marshal(fileOpts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Load reads BotOpts from a JSON config file.
func (codec BotOptsFileJSONCodec) Load(filename string) (*BotOpts, error) {
	return LoadBotOptsFile(codec, filename)
}

// Save writes BotOpts to a JSON config file.
func (codec BotOptsFileJSONCodec) Save(filename string, opts *BotOpts) error {
	return SaveBotOptsFile(codec, filename, opts)
}

var envParameterRegex = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// BotOptsFileCodec decodes and encodes BotOpts file formats.
type BotOptsFileCodec interface {
	FromBytes([]byte) (*BotOpts, error)
	ToBytes(*BotOpts) ([]byte, error)
	Load(filename string) (*BotOpts, error)
	Save(filename string, opts *BotOpts) error
}

// LoadBotOptsFile reads a config file, expands env placeholders, and decodes BotOpts.
func LoadBotOptsFile(codec BotOptsFileCodec, filename string) (*BotOpts, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	data = expandEnvPlaceholdersInFile(data)
	return codec.FromBytes(data)
}

// SaveBotOptsFile encodes BotOpts with codec and writes the result to filename.
func SaveBotOptsFile(codec BotOptsFileCodec, filename string, opts *BotOpts) error {
	data, err := codec.ToBytes(opts)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func expandEnvPlaceholdersInFile(data []byte) []byte {
	return envParameterRegex.ReplaceAllFunc(data, func(match []byte) []byte {
		group := envParameterRegex.FindSubmatch(match)
		if len(group) != 2 {
			return match
		}
		key := group[1]
		value := os.Getenv(string(key))
		return []byte(value)
	})
}
