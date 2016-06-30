package rets

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"strconv"
	"strings"
	"time"
)

// Digest http://en.wikipedia.org/wiki/Digest_access_authentication
type Digest struct {
	Realm      string
	Nonce      string
	Algorithm  string
	Opaque     string
	Qop        string
	NonceCount int
}

func parseChallenge(chall string) (*Digest, error) {
	chall = strings.Trim(chall, " \r\n")
	if !strings.HasPrefix(strings.ToLower(chall), "digest ") {
		return nil, errors.New("Invalid challenge. Challenge must begin with \"digest\"")
	}

	chall = strings.Trim(chall[7:], " \r\n")

	digest, err := parseDirectives(chall, &Digest{Algorithm: "MD5", NonceCount: 1})
	if err != nil {
		return nil, err
	}
	return digest, nil
}

func parseDirectives(directives string, d *Digest) (*Digest, error) {
	for _, e := range strings.Split(directives, ",") {
		directive := strings.SplitN(strings.Trim(e, " \r\n"), "=", 2)
		if len(directive) < 2 {
			return nil, errors.New("Invalid challenge")
		}
		switch directive[0] {
		case "algorithm":
			d.Algorithm = strings.Trim(directive[1], "\"")
		case "domain":
			break
		case "qop":
			d.Qop = strings.Trim(directive[1], "\"")
			if d.Qop == strings.ToLower("auth-int") {
				return nil, errors.New("auth-int is not supported")
			}
		case "nonce":
			d.Nonce = strings.Trim(directive[1], "\"")
		case "opaque":
			d.Opaque = strings.Trim(directive[1], "\"")
		case "realm":
			d.Realm = strings.Trim(directive[1], "\"")
		case "stale":
			break
		default:
			return nil, errors.New("Invalid challenge. Cannot parse directives.")
		}
	}
	return d, nil
}

// NewDigest ...
func NewDigest(chall string) (*Digest, error) {
	c, err := parseChallenge(chall)
	if err != nil {
		return nil, err
	}
	return &Digest{
		Realm:      c.Realm,
		Nonce:      c.Nonce,
		Algorithm:  c.Algorithm,
		Opaque:     c.Opaque,
		Qop:        c.Qop,
		NonceCount: 1,
	}, nil
}

func (d *Digest) createCnonce() string {
	asString := strconv.FormatInt(time.Now().Unix(), 10)
	return md5ThenHex(md5.New(), asString)
}

func (d *Digest) createHa1(username, password, cnonce string, hasher hash.Hash) string {
	//Assuming MD5 or unspecified
	a1 := strings.Join([]string{username, d.Realm, password}, ":")
	ha1 := md5ThenHex(hasher, a1)

	if "md5-sess" == strings.ToLower(d.Algorithm) {
		md5sess := strings.Join([]string{ha1, d.Nonce, cnonce}, ":")
		ha1 = md5ThenHex(hasher, md5sess)
	}

	return ha1
}

func (d *Digest) createHa2(method, uri string, hasher hash.Hash) string {
	//Assuming qop is auth or unspecified
	a2 := method + ":" + uri
	return md5ThenHex(hasher, a2)
}

func (d *Digest) createResponse(ha1, ha2, nc, cnonce string, hasher hash.Hash) string {
	//qop is unspecified
	response := strings.Join([]string{ha1, d.Nonce, ha2}, ":")

	//qop is auth
	if "auth" == strings.ToLower(d.Qop) {
		response = strings.Join([]string{ha1, d.Nonce, nc, cnonce, d.Qop, ha2}, ":")
	}
	return md5ThenHex(hasher, response)
}

func (d *Digest) computeAuthorization(username, password, method, uri, cnonce string) string {
	nc := fmt.Sprintf("%08x", d.NonceCount)
	ha1 := d.createHa1(username, password, cnonce, md5.New())
	ha2 := d.createHa2(method, uri, md5.New())
	response := d.createResponse(ha1, ha2, nc, cnonce, md5.New())

	sl := []string{fmt.Sprintf(`username="%s"`, username)}
	sl = append(sl, fmt.Sprintf(`realm="%s"`, d.Realm))
	sl = append(sl, fmt.Sprintf(`nonce="%s"`, d.Nonce))
	sl = append(sl, fmt.Sprintf(`uri="%s"`, uri))
	sl = append(sl, fmt.Sprintf(`response="%s"`, response))
	if d.Qop != "" {
		sl = append(sl, fmt.Sprintf("qop=%s", d.Qop))
		sl = append(sl, fmt.Sprintf("nc=%s", nc))
		sl = append(sl, fmt.Sprintf(`cnonce="%s"`, cnonce))
	}
	if d.Algorithm != "" {
		sl = append(sl, fmt.Sprintf(`algorithm="%s"`, d.Algorithm))
	}
	if d.Opaque != "" {
		sl = append(sl, fmt.Sprintf(`opaque="%s"`, d.Opaque))
	}

	d.NonceCount++
	return fmt.Sprintf("Digest %s", strings.Join(sl, ", "))
}

// CreateDigestResponse ...
func (d *Digest) CreateDigestResponse(username, password, method, uri string) string {
	return d.computeAuthorization(username, password, method, uri, d.createCnonce())
}

func md5ThenHex(hasher hash.Hash, value string) string {
	hasher.Reset()
	io.WriteString(hasher, value)
	return hex.EncodeToString(hasher.Sum(nil))
}
