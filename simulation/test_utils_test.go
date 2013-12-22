package simulation

import (
	"errors"
	"testing"
)

type fatalferWatcher struct {
	fatalfCalled bool
}

func (f *fatalferWatcher) Fatalf(format string, args ...interface{}) {
	f.fatalfCalled = true
}

func TestUnexpectedErrorDoesntCallFatalIfNoError(t *testing.T) {
	watcher := &fatalferWatcher{false}
	UnexpectedError(nil, "don't care", watcher)
	if watcher.fatalfCalled {
		t.Fatalf("Didn't expect Fatalf to be called with no error\n")
	}
}

func TestUnexpectedErrorCallsFatalIfError(t *testing.T) {
	watcher := &fatalferWatcher{false}
	UnexpectedError(errors.New("Any error"), "don't care", watcher)
	if !watcher.fatalfCalled {
		t.Fatalf("Expected Fatalf to be called with an error\n")
	}
}
