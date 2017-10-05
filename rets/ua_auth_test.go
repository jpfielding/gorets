package rets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUaAuth(t *testing.T) {
	uaauth := UserAgentAuthentication{
		UserAgent:         "user-agent",
		UserAgentPassword: "user-agent-pw",
	}
	header1 := uaauth.header("request-id", "session-id", "rets-version")
	assert.Equal(t, "Digest 052c1af08431d3cc9717a37f9d6de169", header1, "bad auth")

	header2 := uaauth.header("", "", "rets-version")
	assert.Equal(t, "Digest 73cc7ccfe417292b1155c5ccee7fbdab", header2, "bad auth")
}
