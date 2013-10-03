/**

 */
package gorets

import (
	"testing"
)

func TestUaAuth(t *testing.T) {
	uauth := CalculateUaAuthHeader("user-agent","user-agent-pw", "request-id","session-id","rets-version")
	AssertEquals(t, "bad auth", "Digest 052c1af08431d3cc9717a37f9d6de169", uauth)
}
