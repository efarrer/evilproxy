package debug

import (
	"errors"
	"testing"
)

func TestPanicErrorDoesntPanicOnNil(t *testing.T) {
	defer func() {
		if nil != recover() {
			t.Fatalf("PanicOnError shouldn't panic on nil\n")
		}
	}()
	PanicOnError(nil)
}

func TestPanicErrorDoesPanicWithError(t *testing.T) {
	err := errors.New("test")
	defer func() {
		actualErr := recover()
		if err != actualErr {
			t.Fatalf("PanicOnError should have paniced with %v, but got %v instead\n", err, actualErr)
		}
	}()
	PanicOnError(err)
}

func TestOutstandingGoRoutinesReturnsZeroOnEmptyTrace(t *testing.T) {
	if cnt, str := OutstandingGoRoutines(""); cnt != 0 || str != "" {
		t.Fatalf("Expected no goroutines got %v,%v\n", cnt, str)
	}
}

func TestOutstandingGoRoutinesSkipsGoroutineLines(t *testing.T) {
	if cnt, str := OutstandingGoRoutines("goroutine 1 [running]:\n"); cnt != 0 || str != "" {
		t.Fatalf("Expected no goroutines got %v,%v\n", cnt, str)
	}
}

func TestOutstandingGoRoutinesSkipsRuntimeLines(t *testing.T) {
	if cnt, str := OutstandingGoRoutines("runtime/pprof.(*Profile).WriteTo\n/usr/local/go/src/pkg/runtime"); cnt != 0 || str != "" {
		t.Fatalf("Expected no runtime got %v,%v\n", cnt, str)
	}
}

func TestOutstandingGoRoutinesSkipsMainMainLines(t *testing.T) {
	if cnt, str := OutstandingGoRoutines("main.main()\n"); cnt != 0 || str != "" {
		t.Fatalf("Expected no runtime got %v,%v\n", cnt, str)
	}
}

func TestOutstandingGoRoutinesCountsOtherLines(t *testing.T) {
	if cnt, str := OutstandingGoRoutines("    /home/efarrer/Private/gocode/src/evilproxy/evilproxy.go:82 +0x5ff\n"); cnt != 1 || str == "" {
		t.Fatalf("Expected one runtime got %v,%v\n", cnt, str)
	}
}

func TestOnlyOneGoRoutineIsCountedPerSection(t *testing.T) {
	input := "    /home/efarrer/Private/gocode/src/evilproxy/evilproxy.go:82 +0x5ff\n" +
		"    /home/efarrer/Private/gocode/src/evilproxy/evilproxy.go:82 +0x5ff\n"

	if cnt, str := OutstandingGoRoutines(input); cnt != 1 || str == "" {
		t.Fatalf("Expected one runtime got %v,%v\n", cnt, str)
	}
}

func TestEachSectionIsANewGoRoutine(t *testing.T) {
	input := "    /home/efarrer/Private/gocode/src/evilproxy/evilproxy.go:82 +0x5ff\n" +
		"\n" +
		"    /home/efarrer/Private/gocode/src/evilproxy/evilproxy.go:82 +0x5ff\n"

	if cnt, str := OutstandingGoRoutines(input); cnt != 2 || str == "" {
		t.Fatalf("Expected two runtime got %v,%v\n", cnt, str)
	}
}

func TestRealWorldOutput(t *testing.T) {
	input := "goroutine 1 [running]:\n" +
		"runtime/pprof.writeGoroutineStacks(0x7fa3b90837c8, 0xc210050540, 0x7fa3b9079000, 0xd60000c21000a190)\n" +
		"    /usr/local/go/src/pkg/runtime/pprof/pprof.go:511 +0x7c\n" +
		"runtime/pprof.writeGoroutine(0x7fa3b90837c8, 0xc210050540, 0x2, 0xd6c89f9e00000000, 0x1b90837c8)\n" +
		"    /usr/local/go/src/pkg/runtime/pprof/pprof.go:500 +0x3c\n" +
		"runtime/pprof.(*Profile).WriteTo(0x680500, 0x7fa3b90837c8, 0xc210050540, 0x2, 0x1, ...)\n" +
		"    /usr/local/go/src/pkg/runtime/pprof/pprof.go:229 +0xb4\n" +
		"main.main()\n" +
		"    /home/efarrer/Private/gocode/src/evilproxy/evilproxy.go:82 +0x5ff\n"

	if cnt, str := OutstandingGoRoutines(input); cnt != 1 || str == "" {
		t.Fatalf("Expected one runtime got %v,%v\n", cnt, str)
	}
}
