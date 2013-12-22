package simulation

import (
	"runtime"
)

type fatalfer interface {
	Fatalf(format string, args ...interface{})
}

func UnexpectedError(err error, action string, t fatalfer) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("Got an unexpected error(%v) while %v in %v:%v.\n", err, action, file, line)
	}
}
