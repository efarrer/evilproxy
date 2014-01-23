package main

import (
	"errors"
	"fmt"
	"github.com/efarrer/gofutures"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func stdParallel() gofutures.FutureGroup {
	return gofutures.ParallelFutureGroup(
		gofutures.NewConcatenateStringValueReducer(""),
		gofutures.NewConcatenateErrorReducer("\n"))
}

func stdSerial() gofutures.FutureGroup {
	return gofutures.SerializeFutureGroup(
		gofutures.NewConcatenateStringValueReducer(""),
		gofutures.NewConcatenateErrorReducer("\n"))
}

func stdParallelFuture(futures ...gofutures.Future) gofutures.Future {
	return stdParallel().AddFutures(futures...).Finalize()
}

func stdSerialFuture(futures ...gofutures.Future) gofutures.Future {
	return stdSerial().AddFutures(futures...).Finalize()
}

func createEchoFuture(msg string) gofutures.Future {
	return gofutures.NewFuncFuture(func() (interface{}, error) {
		return msg, nil
	})
}

func createCmdFuture(cmd string, arg ...string) gofutures.Future {
	return gofutures.NewFuncFuture(func() (interface{}, error) {
		cmd := exec.Command(cmd, arg...)
		bytes, err := cmd.CombinedOutput()
		return string(bytes), err
	})
}

func createDirsCmdFuture(dirs []string, cmd string, arg ...string) gofutures.Future {
	parallelGroup := stdParallel()

	for _, dir := range dirs {
		parallelGroup.AddFutures(createCmdFuture(cmd, append(arg, dir)...))
	}
	return parallelGroup.Finalize()
}

func createCoverVerifyFuture(coverPath string) gofutures.Future {
	return gofutures.NewFuncFuture(func() (interface{}, error) {
		// Make sure the cover file exists
		if _, err := os.Stat(coverPath); err != nil {
			if os.IsNotExist(err) {
				return "", nil
			}
		}

		cmd := createCmdFuture("go", "tool", "cover", "-func="+coverPath)
		err := cmd.Start()
		if err != nil {
			return "", err
		}
		matcherRe, err := regexp.Compile("100.0%")
		if err != nil {
			return "", err
		}
		output, err := cmd.Results()
		lines := strings.Split(output.(string), "\n")
		for _, line := range lines {
			if line != "" && !matcherRe.MatchString(line) {
				return line, errors.New(fmt.Sprintf("Unit test code coverage is too low.\n\t%v.", line))
			}
		}
		return "", nil
	})
}

func main() {
	gotoStartDirectory()

	dirs := getCurrentSubDirs()

	fmtFutureMsg := createEchoFuture("Format the code.\n")
	fmtFuture := createDirsCmdFuture(dirs, "go", "fmt")

	vetFutureMsg := createEchoFuture("Vet the code.\n")
	vetFuture := createDirsCmdFuture(dirs, "go", "vet")

	testFutureMsg := createEchoFuture("Test the code.\n")
	testFuture := createDirsCmdFuture(dirs, "go", "test", "-race")

	// Generate the path to the coverage file
	coverPath := func(dir string) string {
		return "./" + dir + ".out"
	}

	coverGroup := stdParallel()
	coverGroup.AddFutures(createEchoFuture("Generate code coverage.\n"))
	for _, dir := range dirs {
		coverGroup.AddFutures(createCmdFuture("go", "test", "-coverprofile", coverPath(dir), "./"+dir))
	}

	verifyGroup := stdParallel()
	verifyGroup.AddFutures(createEchoFuture("Verify code coverage.\n"))
	for _, dir := range dirs {
		verifyGroup.AddFutures(createCoverVerifyFuture(coverPath(dir)))
	}

	compileFutureMsg := createEchoFuture("Compile the code.\n")
	compileFuture := createCmdFuture("go", "build", "-race")

	serial := stdSerialFuture(
		fmtFutureMsg, fmtFuture,
		stdParallelFuture(
			vetFutureMsg, vetFuture,
			testFutureMsg, testFuture,
			stdSerialFuture(
				coverGroup.Finalize(),
				verifyGroup.Finalize())),
		compileFutureMsg, compileFuture)

	serial.Start()
	val, err := serial.Results()
	if err != nil {
		log.Fatalf("Error: %v\n\n%v\n", val, err)
	}
	fmt.Printf("%v\n", val)
}
