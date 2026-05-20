package laniakea

import (
	"context"
	"time"
)

// RunnerFn is the function type for a runner. It receives a pointer to
// the Bot and returns an error if execution fails.
type RunnerFn[T AppData] func(*Bot[T]) error

// Runner represents a configurable background or one-time task to be
// executed by a Bot.
//
// Runners are configured using builder methods Async and Every. Once the
// bot's runtime has started executing the runner, it should not be modified.
//
// Execution semantics:
//   - every=0, async=true:  Run once in a goroutine (non-blocking, default).
//   - every=0, async=false: Run once synchronously (blocks runtime startup).
//   - every>0, async=true:  Run repeatedly in a goroutine with the given interval.
//   - every>0, async=false: Invalid configuration — skipped with a warning.
type Runner[T AppData] struct {
	name  string        // Human-readable name for logging
	async bool          // If true, runs in a goroutine; else, runs synchronously
	every time.Duration // Interval between periodic executions; zero means one-shot
	fn    RunnerFn[T]   // The function to execute
}

// NewRunner creates a new Runner with the given name and function.
//
// The default configuration is async=true and every=0, i.e. a one-shot
// goroutine that fires once when the bot runtime starts. Use Async and Every
// to customize this. Do not call builder methods concurrently or after the
// bot runtime has begun executing runners.
func NewRunner[T AppData](name string, fn RunnerFn[T]) Runner[T] {
	return Runner[T]{
		name:  name,
		fn:    fn,
		async: true,
		every: 0,
	}
}

// Async sets whether the runner executes synchronously or asynchronously.
// If true, the runner runs in a goroutine (non-blocking).
// If false, the runner blocks the caller during execution.
//
// Note: periodic runners (Every > 0) require async=true and are skipped with
// a warning when async=false.
func (r Runner[T]) Async(async bool) Runner[T] {
	r.async = async
	return r
}

// Every sets the interval between repeated executions of a periodic runner.
//
// A zero value (the default) keeps the runner one-shot. A positive value
// schedules the runner to fire repeatedly with the given interval and
// requires async=true; periodic sync runners are skipped with a warning.
func (r Runner[T]) Every(timeout time.Duration) Runner[T] {
	r.every = timeout
	return r
}

// ExecRunners executes all runners registered on the Bot with context-based lifecycle management.
//
// Execution semantics by configuration:
//   - every=0, async=true:  Runs once in a goroutine (fire and forget).
//   - every=0, async=false: Runs once synchronously; warns if slower than 2 seconds.
//   - every>0, async=true:  Runs in a loop with the configured interval until ctx.Done().
//   - every>0, async=false: Skipped with a warning (invalid configuration).
//
// Background runners listen for ctx.Done() and gracefully shut down when the context is canceled.
//
// This method is typically called once during bot startup from RunWithContext or
// RunWebhookWithContext.
func (bot *Bot[T]) ExecRunners(ctx context.Context) {
	bot.logger.Infoln("Executing runners...")
	for _, runner := range bot.runners {
		if runner.every > 0 && !runner.async {
			bot.logger.Warnf("Runner %q is periodic but sync; skipping (use Async(true))\n", runner.name)
			continue
		}

		if runner.every == 0 && runner.async {
			// One-time async: fire and forget
			bot.runnerOnceWG.Add(1)
			go func(r Runner[T]) {
				defer bot.runnerOnceWG.Done()
				startedAt := time.Now()
				err := r.fn(bot)
				bot.safeEmitEvent(ctx, RunnerFinishedEvent{
					Name:     r.name,
					Duration: time.Since(startedAt),
					Err:      err,
				})
				if err != nil {
					bot.safeEmitEvent(ctx, ErrorEvent{
						Plugin:      "bot",
						HandlerKind: HandlerRunnerKind,
						HandlerName: r.name,
						Err:         err,
						UserFacing:  false,
					})
					bot.logger.Warnf("Runner %s failed: %s\n", r.name, err)
				}
			}(runner)
		} else if runner.every == 0 && !runner.async {
			// One-time sync: block until done
			t := time.Now()
			err := runner.fn(bot)
			elapsed := time.Since(t)
			bot.safeEmitEvent(ctx, RunnerFinishedEvent{
				Name:     runner.name,
				Duration: elapsed,
				Err:      err,
			})
			if err != nil {
				bot.safeEmitEvent(ctx, ErrorEvent{
					Plugin:      "bot",
					HandlerKind: HandlerRunnerKind,
					HandlerName: runner.name,
					Err:         err,
					UserFacing:  false,
				})
				bot.logger.Warnf("Runner %s failed: %s\n", runner.name, err)
			}
			if elapsed > time.Second*2 {
				bot.logger.Warnf("Runner %s too slow. Elapsed time %v >= 2s\n", runner.name, elapsed)
			}
		} else if runner.every > 0 && runner.async {
			// Background loop: periodic execution with graceful shutdown
			bot.runnerBgWG.Add(1)
			go func(r Runner[T]) {
				defer bot.runnerBgWG.Done()
				ticker := time.NewTicker(r.every)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
					}
					// When both ctx.Done() and ticker.C are ready at the same
					// time, Go's select picks one at random. Re-check ctx so a
					// late tick after cancellation does not fire one extra
					// invocation past shutdown.
					if ctx.Err() != nil {
						return
					}
					startedAt := time.Now()
					err := r.fn(bot)
					bot.safeEmitEvent(ctx, RunnerFinishedEvent{
						Name:     r.name,
						Duration: time.Since(startedAt),
						Err:      err,
					})
					if err != nil {
						bot.safeEmitEvent(ctx, ErrorEvent{
							Plugin:      "bot",
							HandlerKind: HandlerRunnerKind,
							HandlerName: r.name,
							Err:         err,
							UserFacing:  false,
						})
						bot.logger.Warnf("Runner %s failed: %s\n", r.name, err)
					}
				}
			}(runner)
		}
	}
}
