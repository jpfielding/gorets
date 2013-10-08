/**
	http://en.wikipedia.org/wiki/Digest_access_authentication
 */
package gorets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"strconv"
	"time"
	"hash"
)

/**
	this method wraps up the functionality of creating a digest response from a challenge

	TODO - http://tools.ietf.org/html/rfc2617
	TODO - CLEAN THIS TURD UP!!!
	TODO - actually keep a legit nonce counter, possibly in a digest struct?!?!?
 */
func DigestResponse(challenge, username, password, method, uri string) string {
	cType, args := parseChallenge(challenge)
	return challengeResponse(cType, args, username, password, method, uri, createCnonce(), "00000001")
}

func challengeResponse(cType string, challenge map[string]string, username, password, method, uri, cnonce, nc string) string {
	params := make([]string, len(challenge)+5)
	response := digest(challenge, username, password, method, uri, cnonce, nc)
	params[0] = fmt.Sprintf("username=\"%s\"", username)
	params[1] = fmt.Sprintf("uri=\"%s\"", uri)
	params[2] = fmt.Sprintf("response=\"%s\"", response)
	params[3] = fmt.Sprintf("cnonce=\"%s\"", cnonce)
	params[4] = fmt.Sprintf("nc=\"%s\"", nc)
	i := 5
	for k,v := range challenge {
		params[i] = fmt.Sprintf("%s=\"%s\"", k, v)
		i++
	}

	return cType +" "+ strings.Join(params, ", ")
}

func createCnonce() string {
	asString := strconv.FormatInt(time.Now().Unix(),10)
	return md5ThenHex(md5.New(), asString)
}


// TODO convert this into a Digest struct
func parseChallenge(challenge string) (string,map[string]string) {
	pieces := strings.SplitAfterN(challenge, " ", 2)
	cType, challenge := strings.TrimSpace(pieces[0]), pieces[1]
	parts := map[string]string{}
	for _,part := range strings.Split(challenge, ",") {
		part = strings.TrimSpace(part)
		split := strings.Split(part,"=")
		parts[split[0]] = strings.TrimSuffix(strings.TrimPrefix(split[1],"\""),"\"")
	}
	_,hasAlgorithm := parts["algorithm"]
	if !hasAlgorithm {
		parts["algorithm"] = "MD5"
	}

	return cType, parts
}


func digest(challenge map[string]string, username, password, method, uri, cnonce, nc string) string {
	realm := challenge["realm"]
	nonce := challenge["nonce"]
	qop := challenge["qop"]
	algorithm := challenge["algorithm"]

	hasher := md5.New()

	a1 := strings.Join([]string{username, realm, password},":")
	ha1 := md5ThenHex(hasher, a1)

	a2 := method +":"+ uri
	ha2 := md5ThenHex(hasher, a2)

	// md5-sess:
	// ha1 = md5(a1) = md5(md5(username:realm:password):nonce:cnonce)
	if "md5-sess" == strings.ToLower(algorithm) {
		md5sess := strings.Join([]string{ha1, nonce, cnonce}, ":")
		ha1 = md5ThenHex(hasher, md5sess)
	}

	switch qop {
	case "auth":
		response := strings.Join([]string{ha1, nonce, nc, cnonce, qop, ha2},":")
		return md5ThenHex(hasher, response)
	case "auth-int":
		// TODO - requires hash of entity body
		return "qop: auth-int not yet supported"
	}
	response := strings.Join([]string{ha1,nonce,ha2},":")
	return md5ThenHex(hasher, response)
}

func md5ThenHex(hasher hash.Hash, value string) string {
	hasher.Reset()
	io.WriteString(hasher, value)
	return hex.EncodeToString(hasher.Sum(nil))
}
