/**
	TODO
 */
package gorets

import (
	"fmt"
	"testing"
)

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


