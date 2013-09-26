/**
 * Created with IntelliJ IDEA.
 * User: jp
 * Date: 9/20/13
 * Time: 11:57 AM
 * To change this template use File | Settings | File Templates.
 */
package main

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

