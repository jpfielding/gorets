package gorets_client

import (
	"testing"
)

func TestUaAuth(t *testing.T) {
	uauth := CalculateUaAuthHeader("user-agent", "user-agent-pw", "request-id", "session-id", "rets-version")
	assert(t, "Digest 052c1af08431d3cc9717a37f9d6de169" == uauth, "bad auth")

	uauth2 := CalculateUaAuthHeader("user-agent", "user-agent-pw", "", "", "rets-version")
	assert(t, "Digest 73cc7ccfe417292b1155c5ccee7fbdab" == uauth2, "bad auth")
}
