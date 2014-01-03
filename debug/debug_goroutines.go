package debug

import (
	"bytes"
	"regexp"
)

var goroutineRx *regexp.Regexp
var runtimeRx *regexp.Regexp
var runtimeRx2 *regexp.Regexp
var mainRx *regexp.Regexp
var emptyLineRx *regexp.Regexp

func init() {
	var err error
	goroutineRx, err = regexp.Compile("^goroutine.*")
	PanicOnError(err)
	runtimeRx, err = regexp.Compile("^runtime.*")
	PanicOnError(err)
	runtimeRx2, err = regexp.Compile(".*pkg.runtime.*")
	PanicOnError(err)
	mainRx, err = regexp.Compile("^main.main()")
	PanicOnError(err)
	emptyLineRx, err = regexp.Compile("^[ \t\r\n]*$")
	PanicOnError(err)
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func OutstandingGoRoutines(trace string) (int, string) {
	foundMatchingLine := false
	count := 0
	ret := ""
	buffer := bytes.NewBuffer([]byte(trace))

	for {
		str, err := buffer.ReadString('\n')
		if err != nil && str == "" {
			break
		}
		if goroutineRx.MatchString(str) {
			continue
		}
		if runtimeRx.MatchString(str) {
			continue
		}
		if runtimeRx2.MatchString(str) {
			continue
		}
		if mainRx.MatchString(str) {
			continue
		}

		// If we foundMatchingLine and we're at the end of the section then
		// count it
		if emptyLineRx.MatchString(str) && foundMatchingLine {
			count += 1
			foundMatchingLine = false
			continue
		}

		// We found a new match so save it
		ret += "Outstanding: " + str
		foundMatchingLine = true
	}

	// We're at the last line so the section has effectively ended
	if foundMatchingLine {
		count += 1
	}

	return count, ret
}
