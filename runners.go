package laniakea

import (
	"time"
)

type RunnerFn func(*Bot) error
type RunnerBuilder struct {
	name    string
	onetime bool
	async   bool
	timeout time.Duration
	fn      RunnerFn
}
type Runner struct {
	Name    string
	Onetime bool
	Async   bool
	Timeout time.Duration
	Fn      RunnerFn
}

func NewRunner(name string, fn RunnerFn) *RunnerBuilder {
	return &RunnerBuilder{
		name: name, fn: fn, async: true,
	}
}
func (b *RunnerBuilder) Onetime(onetime bool) *RunnerBuilder {
	b.onetime = onetime
	return b
}
func (b *RunnerBuilder) Async(async bool) *RunnerBuilder {
	b.async = async
	return b
}
func (b *RunnerBuilder) Timeout(timeout time.Duration) *RunnerBuilder {
	b.timeout = timeout
	return b
}
func (b *RunnerBuilder) Build() Runner {
	return Runner{
		Name: b.name, Onetime: b.onetime, Async: b.async, Fn: b.fn, Timeout: b.timeout,
	}
}

func (b *Bot) ExecRunners() {
	for _, runner := range b.runners {
		if !runner.Onetime && !runner.Async {
			b.logger.Warnf("Runner %s not onetime, but sync\n", runner.Name)
			continue
		}
		if !runner.Onetime && runner.Async && runner.Timeout == (time.Second*0) {
			b.logger.Warnf("Background runner \"%s\" should have timeout", runner.Name)
		}

		if runner.Async && runner.Onetime {
			go func() {
				err := runner.Fn(b)
				if err != nil {
					b.logger.Warnf("Runner %s failed: %s\n", runner.Name, err)
				}
			}()
		} else if !runner.Async && runner.Onetime {
			err := runner.Fn(b)
			if err != nil {
				b.logger.Warnf("Runner %s failed: %s\n", runner.Name, err)
			}
		} else if !runner.Onetime {
			go func() {
				for {
					err := runner.Fn(b)
					if err != nil {
						b.logger.Warnf("Runner %s failed: %s\n", runner.Name, err)
					}
					time.Sleep(runner.Timeout)
				}
			}()
		}
	}
}
