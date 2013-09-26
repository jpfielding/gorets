/**
	http://en.wikipedia.org/wiki/Digest_access_authentication
 */
package gorets

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
	"strconv"
	"time"
)

/**
	this method wraps up the functionality of creating a digest response from a challenge
 */
func DigestResponse(challenge, username, password, method, uri string) string {
	cType, args := parseChallenge(challenge)
	return challengeResponse(cType, args, username, password, method, uri)
}

func challengeResponse(cType string, challenge map[string]string, username, password, method, uri string) string {
	params := []string{
				"username=\"" +username+"\"",
				"realm=\""+challenge["realm"]+"\"",
				"nonce=\""+challenge["nonce"]+"\"",
				"uri=\""+uri+"\"",
				"response=\""+digest(challenge, username, password, method, uri)+"\"",
	}

	alg,hasAlg := challenge["algorithm"]
	if hasAlg {
		params = append(params, "algorithm=\""+alg+"\"")
	}

	opaque, hasOpq := challenge["opaque"]
	if hasOpq {
		params = append(params, "opaque=\""+ opaque +"\"")
	}

	return cType +" "+ strings.Join(params, ", ")
}

func createCnonce() string {
	hasher := md5.New()
	asString := strconv.FormatInt(time.Now().Unix(),10)
	io.WriteString(hasher, asString)
	return hex.EncodeToString(hasher.Sum(nil))
}


// TODO convert this into a Digest struct
func parseChallenge(challenge string) (string,map[string]string) {
	pieces := strings.SplitAfterN(challenge, " ", 2)
	cType, challenge := strings.Trim(pieces[0]," "), pieces[1]
	parts := map[string]string{}
	for _,part := range strings.Split(challenge, ", ") {
		split := strings.Split(part,"=")
		parts[split[0]] = strings.TrimSuffix(strings.TrimPrefix(split[1],"\""),"\"")
	}
	_,hasAlgorithm := parts["algorithm"]
	if !hasAlgorithm {
		parts["algorithm"] = "MD5"
	}

	return cType, parts
}


func digest(challenge map[string]string, username, password, method, uri string) string {
	realm := challenge["realm"]
	nonce := challenge["nonce"]
	//qop := challenge["qop"]
	algorithm := challenge["algorithm"]
	var cnonce string

	a1 := username +":"+ realm +":"+ password

	hasher := md5.New()

	io.WriteString(hasher, a1)
	ha1 := hex.EncodeToString(hasher.Sum(nil))

	// md5-sess:
	// ha1 = md5(a1) = md5(md5(username:realm:password):nonce:cnonce)
	if "md5-sess" == strings.ToLower(algorithm) {
		cnonce = createCnonce()
		hasher.Reset()
		md5sess := ha1 +":"+ nonce +":"+ cnonce
		io.WriteString(hasher, md5sess)
		a1 = hex.EncodeToString(hasher.Sum(nil))
	}

	// TODO QOP VARIANTS
	a2 := method +":"+ uri

	hasher.Reset()
	io.WriteString(hasher, a2)
	ha2 := hex.EncodeToString(hasher.Sum(nil))

	response := ha1 +":"+ nonce +":"+ ha2

	hasher.Reset()
	io.WriteString(hasher, response)
	digest := hex.EncodeToString(hasher.Sum(nil))
	return digest
}
