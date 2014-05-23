/**
http://en.wikipedia.org/wiki/Digest_access_authentication
*/

package gorets_client

import (
	"testing"
)

func TestDigest(t *testing.T) {
	username, password := "user", "passwd"

	wwwAuthenticate := `Digest realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", opaque="6e6f742075736564"`

	expected := `Digest username="user", realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", uri="/platinum/login", response="5f8d366fb430e9b395a84dba52247a35", algorithm="MD5", opaque="6e6f742075736564"`

	t.Logf("CHALLENGE: %s", wwwAuthenticate)

	method, uri := "GET", "/platinum/login"
	cType, challenge := parseChallenge(wwwAuthenticate)
	authorization := challengeResponse(cType, challenge, username, password, method, uri, "", "")

	t.Logf("EXPECTED : %s", expected)
	t.Logf("ACTUAL   : %s", authorization)

	_, e := parseChallenge(expected)
	_, a := parseChallenge(authorization)

	for k, v := range e {
		if a[k] != v {
			AssertEquals(t, k, v, a[k])
		}
	}
}

func TestDigestQopAuth(t *testing.T) {
	username, password := "user", "passwd"

	wwwAuthenticate := `Digest realm="rets@ris.rets.interealty.com",nonce="cbbaac704d6762fa2e4c37bb84ffddb8",opaque="77974390-589f-4ec6-b67c-4a9fc096c03f",qop="auth"`
	expected := `Digest username="user", realm="rets@ris.rets.interealty.com", nonce="cbbaac704d6762fa2e4c37bb84ffddb8", uri="/rets/login", response="3b2ce903001cf00a264c69674d691526", qop=auth, nc=00000001, cnonce="6ac2c2eee85f5c33", algorithm="MD5", opaque="77974390-589f-4ec6-b67c-4a9fc096c03f"`

	t.Logf("CHALLENGE: %s", wwwAuthenticate)

	method, uri := "GET", "/rets/login"
	cType, challenge := parseChallenge(wwwAuthenticate)
	nc := "00000001"
	cnonce := "6ac2c2eee85f5c33"
	authorization := challengeResponse(cType, challenge, username, password, method, uri, cnonce, nc)

	t.Logf("EXPECTED : %s", expected)
	t.Logf("ACTUAL   : %s", authorization)

	_, e := parseChallenge(expected)
	_, a := parseChallenge(authorization)

	for k, v := range e {
		if a[k] != v {
			AssertEquals(t, k, v, a[k])
		}
	}
}
