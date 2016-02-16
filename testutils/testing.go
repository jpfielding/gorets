/**
	https://github.com/benbjohnson/testing
	http://golang.org/pkg/testing/

func TestXxx(t *testing.T) {
	t.Log("bitches love go")
}

func BenchmarkXxx(b *testing.B) {
	b.Log("bitches are fast")
}

func ExampleHello() {
	fmt.Println("hello bitch")
}

*/
package client

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		fatalf(tb, msg, v...)
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		fatalf(tb, "Unexpected error: %s", err.Error())
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// NotOk fails the test if an err is nil.
func NotOk(tb testing.TB, err error) {
	if err == nil {
		fatalf(tb, "Did not find expected error.")
	}
}

func fatalf(tb testing.TB, msg string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	fmt.Fprintf(os.Stderr, "%s:%d:%s\n\n", file, line, fmt.Sprintf(msg, v...))
	if tb == nil {
		os.Exit(1)
	}
	tb.FailNow()
}
