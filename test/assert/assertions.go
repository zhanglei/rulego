package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

//CallerInfo This function is inspired by:
//https://github.com/stretchr/testify/blob/master/assert/assertions.go
func CallerInfo() []string {
	var pc uintptr
	var ok bool
	var file string
	var line int
	var name string

	callers := []string{}
	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case, see #180
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()

		// testing.tRunner is the standard library function that calls
		// tests. Subtests are called directly by tRunner, without going through
		// the Test/Benchmark/Example function that contains the t.Run calls, so
		// with subtests we should break when we hit tRunner, without adding it
		// to the list of callers.
		if name == "testing.tRunner" {
			break
		}

		parts := strings.Split(file, "/")
		if len(parts) > 1 {
			filename := parts[len(parts)-1]
			dir := parts[len(parts)-2]
			if (dir != "assert" && dir != "mock" && dir != "require") || filename == "mock_test.go" {
				callers = append(callers, fmt.Sprintf("%s:%d", file, line))
			}
		}

	}

	return callers
}

func Equal(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%v != %v\n Error Trace:   %s", a, b, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}

func True(t *testing.T, value bool) {
	if !value {
		t.Errorf("%v should be true\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func False(t *testing.T, value bool) {
	if value {
		t.Errorf("%v should be false\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func NotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Errorf("%v should be not nil\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func Nil(t *testing.T, value interface{}) {
	if value != nil {
		t.Errorf("%v should be nil\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
