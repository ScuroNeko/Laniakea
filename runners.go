package laniakea

import (
	"time"
)

type RunnerFn[T DbContext] func(*Bot[T]) error
type Runner[T DbContext] struct {
	name    string
	onetime bool
	async   bool
	timeout time.Duration
	fn      RunnerFn[T]
}

func NewRunner[T DbContext](name string, fn RunnerFn[T]) *Runner[T] {
	return &Runner[T]{
		name: name, fn: fn, async: true,
	}
}
func (b *Runner[T]) Onetime(onetime bool) *Runner[T] {
	b.onetime = onetime
	return b
}
func (b *Runner[T]) Async(async bool) *Runner[T] {
	b.async = async
	return b
}
func (b *Runner[T]) Timeout(timeout time.Duration) *Runner[T] {
	b.timeout = timeout
	return b
}

func (bot *Bot[T]) ExecRunners() {
	bot.logger.Infoln("Executing runners...")
	for _, runner := range bot.runners {
		if !runner.onetime && !runner.async {
			bot.logger.Warnf("Runner %s not onetime, but sync\n", runner.name)
			continue
		}
		if !runner.onetime && runner.async && runner.timeout == (time.Second*0) {
			bot.logger.Warnf("Background runner \"%s\" should have timeout", runner.name)
		}

		if runner.async && runner.onetime {
			go func() {
				err := runner.fn(bot)
				if err != nil {
					bot.logger.Warnf("Runner %s failed: %s\n", runner.name, err)
				}
			}()
		} else if !runner.async && runner.onetime {
			t := time.Now()
			err := runner.fn(bot)
			if err != nil {
				bot.logger.Warnf("Runner %s failed: %s\n", runner.name, err)
			}
			elapsed := time.Since(t)
			if elapsed > time.Second*2 {
				bot.logger.Warnf("Runner %s too slow. Elapsed time %s>=2s", runner.name, elapsed)
			}
		} else if !runner.onetime {
			go func() {
				for {
					err := runner.fn(bot)
					if err != nil {
						bot.logger.Warnf("Runner %s failed: %s\n", runner.name, err)
					}
					time.Sleep(runner.timeout)
				}
			}()
		}
	}
}
