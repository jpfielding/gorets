/**
 * Created with IntelliJ IDEA.
 * User: jp
 * Date: 9/20/13
 * Time: 11:55 AM
 * To change this template use File | Settings | File Templates.
 */
/*
Package main for the start of a rets client in go
*/
package main

import (
	// TODO figure out the right way to import locally
	"./auth"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)


func main () {
	username := flag.String("username", "", "Username for the RETS server")
	password := flag.String("password", "", "Password for the RETS server")
	loginUrl := flag.String("login-url", "", "Login URL for the RETS server")
	userAgent := flag.String("user-agent","Threewide/1.0","User agent for the RETS client")

	flag.Parse()


	// setup and custom http params here
	client := &http.Client{
	}

	req, err := http.NewRequest("GET", *loginUrl, nil)
	if err != nil {
		fmt.Println(err)
		// handle error
	}

	req.Header.Add("User-Agent", *userAgent)
	req.Header.Add("RETS-Version", "RETS/1.5")
	req.Header.Add("Accept", "*/*")

	fmt.Println("REQUEST:", req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("RESPONSE:",resp.Header)

	// digest auth
	if resp.StatusCode == 401 {
		challenge := resp.Header.Get("Www-Authenticate")
		if !strings.HasPrefix(strings.ToLower(challenge), "digest") {
			panic("recognized challenge: "+ challenge)
		}
		req.Header.Add("Authorization", auth.DigestResponse(challenge, *username, *password, req.Method, req.URL.Path))
		resp, err = client.Do(req)
	}

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("RESPONSE (AUTH):",resp)
	hah, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(hah))
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
	httputil.DumpResponse(resp,true)

//	fmt.Println(body)
}
