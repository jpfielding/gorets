/**
http://en.wikipedia.org/wiki/Digest_access_authentication
*/

package client

import (
	"testing"
	"crypto/md5"
)

func TestParseChallenge(t *testing.T) {
	wwwAuthenticate := `Digest realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", opaque="6e6f742075736564"`

	digest, err := NewDigest(wwwAuthenticate)
	if err != nil {
		t.Fail()
	}

	expected := &Digest {
		Realm: "users@mris.com",
		Nonce: "31333739363738363932323632201e00230a639db77779b354d601ee5d2e",
		Algorithm: "MD5",
    	Opaque: "6e6f742075736564",
    	Qop: "",
    	NonceCount: 1,
	}

	equals(t, expected, digest)
}

func TestCreateHa1Md5(t *testing.T) {
	username := "user"
	password := "passwd"
	cnonce := "6ac2c2eee85f5c33"

	expected := "5c0da895491be93b455ebb56f7ae0a9f"
	actual := createNoQopDigest().createHa1(username, password, cnonce, md5.New())

	equals(t, expected, actual)
}

func TestCreateHa1Md5Sess(t *testing.T) {
	username := "user"
	password := "passwd"
	cnonce := "6ac2c2eee85f5c33"

	wwwAuthenticate := `Digest realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", opaque="6e6f742075736564, algorithm="MD5-sess"`

	digest, err := NewDigest(wwwAuthenticate)
	if err != nil {
		t.Fail()
	}

	expected := "f1843845124dcba66fef064f5aa7a782"
	actual := digest.createHa1(username, password, cnonce, md5.New())

	equals(t, expected, actual)
}

func TestCreateHa2(t *testing.T) {
	method, uri := "GET", "/platinum/login"

	expected := "1b11c4ebed4a67753078be8020ea9d19"
	actual := createNoQopDigest().createHa2(method, uri, md5.New())

	equals(t, expected, actual)
}

func TestCreateResponseNoQop(t *testing.T) {
	nc := "00000001"
	cnonce := "6ac2c2eee85f5c33"
	ha1 := "5c0da895491be93b455ebb56f7ae0a9f"
	ha2 := "1b11c4ebed4a67753078be8020ea9d19"

	expected := "5f8d366fb430e9b395a84dba52247a35"
	actual := createNoQopDigest().createResponse(ha1, ha2, nc, cnonce, md5.New())

	equals(t, expected, actual)
}

func TestCreateResponseQopAuth(t *testing.T) {
	nc := "00000001"
	cnonce := "6ac2c2eee85f5c33"
	ha1 := "5c0da895491be93b455ebb56f7ae0a9f"
	ha2 := "1b11c4ebed4a67753078be8020ea9d19"

	expected := "28552064c4cde9a3af7610e7ae286d50"
	actual := createAuthDigest().createResponse(ha1, ha2, nc, cnonce, md5.New())

	equals(t, expected, actual)
}

func TestDigest(t *testing.T) {
	username, password := "user", "passwd"

	wwwAuthenticate := `Digest realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", opaque="6e6f742075736564"`

	expected := `Digest username="user", realm="users@mris.com", nonce="31333739363738363932323632201e00230a639db77779b354d601ee5d2e", uri="/platinum/login", response="5f8d366fb430e9b395a84dba52247a35", algorithm="MD5", opaque="6e6f742075736564"`

	t.Logf("CHALLENGE: %s", wwwAuthenticate)

	method, uri := "GET", "/platinum/login"
	
	actual := createNoQopDigest().CreateDigestResponse(username, password, method, uri)

	equals(t, expected, actual)
}

func TestDigestQopAuth(t *testing.T) {
	username, password := "user", "passwd"

	wwwAuthenticate := `Digest realm="rets@ris.rets.interealty.com",nonce="cbbaac704d6762fa2e4c37bb84ffddb8",opaque="77974390-589f-4ec6-b67c-4a9fc096c03f",qop="auth"`
	expected := `Digest username="user", realm="rets@ris.rets.interealty.com", nonce="cbbaac704d6762fa2e4c37bb84ffddb8", uri="/rets/login", response="3b2ce903001cf00a264c69674d691526", qop=auth, nc=00000001, cnonce="6ac2c2eee85f5c33", algorithm="MD5", opaque="77974390-589f-4ec6-b67c-4a9fc096c03f"`

	t.Logf("CHALLENGE: %s", wwwAuthenticate)

	method, uri := "GET", "/rets/login"
	cnonce := "6ac2c2eee85f5c33"

	dig, err := NewDigest(wwwAuthenticate)
	if err != nil {
		t.Fail()
	}
	actual :=  dig.computeAuthorization(username, password, method, uri, cnonce)
	
	equals(t, expected, actual)
}

func createNoQopDigest() *Digest {
	return &Digest{ 
		Realm: "users@mris.com",
		Nonce: "31333739363738363932323632201e00230a639db77779b354d601ee5d2e",
		Algorithm: "MD5",
    	Opaque: "6e6f742075736564",
    	Qop: "",
    	NonceCount: 1,
    }
}

func createAuthDigest() *Digest {
	return &Digest{ 
		Realm: "users@mris.com",
		Nonce: "31333739363738363932323632201e00230a639db77779b354d601ee5d2e",
		Algorithm: "MD5",
    	Opaque: "6e6f742075736564",
    	Qop: "auth",
    	NonceCount: 1,
    }
}
