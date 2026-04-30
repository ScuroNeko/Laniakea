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
// Runners are configured using builder methods: Once(), Async(), Every().
// Once Execute() is called, the Runner should not be modified.
//
// Execution semantics:
//   - once=true, async=false: Run once synchronously (blocks).
//   - once=true, async=true:  Run once in a goroutine (non-blocking).
//   - once=false, async=true: Run repeatedly in a goroutine with timeout.
//   - once=false, async=false: Invalid configuration — ignored with warning.
type Runner[T AppData] struct {
	name  string        // Human-readable name for logging
	once  bool          // If true, runs once; if false, runs periodically
	async bool          // If true, runs in a goroutine; else, runs synchronously
	every time.Duration // Duration to wait between periodic executions (ignored if once=true)
	fn    RunnerFn[T]   // The function to execute
}

// NewRunner creates a new Runner with the given name and function.
// By default, the Runner is configured as async=true (non-blocking).
//
// Builder methods (Once, Async, Every) can be chained to customize behavior.
// DO NOT call builder methods concurrently or after Execute().
func NewRunner[T AppData](name string, fn RunnerFn[T]) Runner[T] {
	return Runner[T]{
		name:  name,
		fn:    fn,
		async: true, // Default: run asynchronously
		every: 0,    // Default: no timeout (ignored if once=true)
	}
}

// Once sets whether the runner executes once or repeatedly.
// If true, the runner runs only once.
// If false, the runner runs in a loop with the configured timeout.
func (r Runner[T]) Once(once bool) Runner[T] {
	r.once = once
	return r
}

// Async sets whether the runner executes synchronously or asynchronously.
// If true, the runner runs in a goroutine (non-blocking).
// If false, the runner blocks the caller during execution.
//
// Note: If once=false and async=false, the runner will be skipped with a warning.
func (r Runner[T]) Async(async bool) Runner[T] {
	r.async = async
	return r
}

// Every sets the duration to wait between repeated executions for
// non-once runners.
//
// If once=true, this value is ignored.
// If once=false and async=true, this timeout determines the sleep interval
// between loop iterations.
//
// A zero value (time.Duration(0)) is allowed but may trigger a warning
// if used with a background (non-once) async runner.
func (r Runner[T]) Every(timeout time.Duration) Runner[T] {
	r.every = timeout
	return r
}

// ExecRunners executes all runners registered on the Bot with context-based lifecycle management.
//
// It logs warnings for misconfigured runners:
//   - Sync, non-once runners are skipped (invalid configuration).
//   - Background (non-once, async) runners without a timeout trigger a warning.
//
// Execution logic:
//   - once + async: Runs once in a goroutine.
//   - once + sync:  Runs once synchronously; warns if slower than 2 seconds.
//   - !once + async: Runs in a loop with timeout between iterations until ctx.Done().
//   - !once + sync: Skipped with warning.
//
// Background runners listen for ctx.Done() and gracefully shut down when the context is canceled.
//
// This method is typically called once during bot startup from RunWithContext or
// RunWebhookWithContext.
func (bot *Bot[T]) ExecRunners(ctx context.Context) {
	bot.logger.Infoln("Executing runners...")
	for _, runner := range bot.runners {
		// Validate configuration
		if !runner.once && !runner.async {
			bot.logger.Warnf("Runner %s not once, but sync — skipping\n", runner.name)
			continue
		}
		if !runner.once && runner.async && runner.every == 0 {
			bot.logger.Warnf("Background runner \"%s\" has no timeout — skipping\n", runner.name)
			continue
		}

		if runner.once && runner.async {
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
		} else if runner.once && !runner.async {
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
		} else if !runner.once && runner.async {
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
				}
			}(runner)
		}
		// Note: !once && !async is already skipped above
	}
}
