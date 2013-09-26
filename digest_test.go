/**
	http://en.wikipedia.org/wiki/Digest_access_authentication
*/

package gorets

import (
	"testing"
)


func TestChallengeDigest(t *testing.T) {
	username, password := "2006973", "TWD1"

	wwwAuthenticate := `Digest realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", opaque="6e6f742075736564"`

	expected := `Digest username="2006973", realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", uri="/platinum/login", response="3c3cff0b1761f96f1627bc2c253ec746", algorithm="MD5", opaque="6e6f742075736564"`

	t.Logf("CHALLENGE: %s", wwwAuthenticate)

	method, uri := "GET", "/platinum/login"
	cType, challenge := parseChallenge(wwwAuthenticate)
	authorization := challengeResponse(cType, challenge, username, password, method, uri)

	t.Logf("EXPECTED : %s", expected)
	t.Logf("ACTUAL   : %s", authorization)

	if expected != authorization {
		t.Errorf("%s != %s", expected, authorization)
	}
}

