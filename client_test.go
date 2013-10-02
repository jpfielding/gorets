/**
	TODO
 */
package gorets

import (
	"fmt"
	"testing"
)

func AssertEquals(t *testing.T, msg, expected, actual string) {
	if actual != expected {
		t.Errorf("%s: %s != %s", msg, expected, actual)
	}
}
func AssertEqualsTrue(t *testing.T, msg string, expected bool) {
	if !expected {
		t.Errorf("%s", msg)
	}
}
func AssertEqualsInt(t *testing.T, msg string, expected, actual int) {
	if actual != expected {
		t.Errorf("%s: %d != %d", msg, expected, actual)
	}
}

/* http://golang.org/pkg/testing/ */
func TestXxx(t *testing.T) {
	t.Log("bitches love go")
}

func BenchmarkXxx(b *testing.B) {
	b.Log("bitches are fast")
}

func ExampleHello() {
	fmt.Println("hello example")
}


