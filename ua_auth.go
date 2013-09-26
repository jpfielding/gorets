/**
	provides the basic mechanism for User Agent authentication for rets
 */
package gorets


import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
)


func CalculateUaAuthHeader(userAgent, userAgentPw, requestId, sessionId, retsVersion string) string {
	hasher := md5.New()

	io.WriteString(hasher, userAgent +":"+ userAgentPw)
	secretHash := hex.EncodeToString(hasher.Sum(nil))

	pieces := secretHash+":"+strings.TrimSpace(requestId)+":"+strings.TrimSpace(sessionId)+":"+retsVersion;

	hasher.Reset()
	io.WriteString(hasher, pieces)
	return "Digest "+ hex.EncodeToString(hasher.Sum(nil))
}
