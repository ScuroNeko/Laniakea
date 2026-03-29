package laniakea

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"git.scuroneko.dev/scuroneko/slog"
)

func TestExecRunnersRunsOnetimeSyncRunner(t *testing.T) {
	var calls atomic.Int32
	bot := &Bot[NoDB]{
		logger: slog.CreateLogger(),
		runners: []Runner[NoDB]{
			NewRunner("sync-once", func(*Bot[NoDB]) error {
				calls.Add(1)
				return nil
			}).Onetime(true).Async(false),
		},
	}

	bot.ExecRunners(context.Background())

	if got := calls.Load(); got != 1 {
		t.Fatalf("unexpected sync runner call count: %d", got)
	}
}

func TestExecRunnersStopsBackgroundRunnerOnCancel(t *testing.T) {
	var calls atomic.Int32
	triggered := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())

	bot := &Bot[NoDB]{
		logger: slog.CreateLogger(),
		runners: []Runner[NoDB]{
			NewRunner("background", func(*Bot[NoDB]) error {
				if calls.Add(1) == 1 {
					triggered <- struct{}{}
				}
				return nil
			}).Timeout(5 * time.Millisecond),
		},
	}

	bot.ExecRunners(ctx)

	select {
	case <-triggered:
	case <-time.After(time.Second):
		t.Fatal("background runner did not execute")
	}

	cancel()
	bot.runnerBgWG.Wait()

	if calls.Load() == 0 {
		t.Fatal("expected background runner to be called at least once")
	}
}
