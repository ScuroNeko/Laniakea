package laniakea

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"git.scuroneko.dev/scuroneko/slog"
)

type runnerObserver struct {
	recordingObserver
}

func TestExecRunnersRunsOnetimeSyncRunner(t *testing.T) {
	var calls atomic.Int32
	bot := &Bot[NoData]{
		logger: slog.CreateLogger(),
		runners: []Runner[NoData]{
			NewRunner("sync-once", func(*Bot[NoData]) error {
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

	bot := &Bot[NoData]{
		logger: slog.CreateLogger(),
		runners: []Runner[NoData]{
			NewRunner("background", func(*Bot[NoData]) error {
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

func TestExecRunnersEmitObserverEvents(t *testing.T) {
	observer := &runnerObserver{}
	wantErr := errors.New("runner failed")

	bot := &Bot[NoData]{
		logger:   slog.CreateLogger(),
		observer: observer,
		runners: []Runner[NoData]{
			NewRunner("sync-once", func(*Bot[NoData]) error {
				return wantErr
			}).Onetime(true).Async(false),
		},
	}

	bot.ExecRunners(context.Background())

	if len(observer.runners) != 1 {
		t.Fatalf("expected one runner-finished event, got %d", len(observer.runners))
	}
	if got := observer.runners[0]; got.Name != "sync-once" || !errors.Is(got.Err, wantErr) {
		t.Fatalf("unexpected runner-finished event: %#v", got)
	}
	if len(observer.errors) != 1 {
		t.Fatalf("expected one error event, got %d", len(observer.errors))
	}
	if got := observer.errors[0]; got.HandlerKind != HandlerRunnerKind || got.HandlerName != "sync-once" || !errors.Is(got.Err, wantErr) {
		t.Fatalf("unexpected runner error event: %#v", got)
	}
}
