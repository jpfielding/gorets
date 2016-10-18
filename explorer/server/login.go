package server

import "net/http"

// Login hash the user/pwd into a tmp space for writing data locally
func Login() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
