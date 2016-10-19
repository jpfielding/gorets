package explorer

import "net/http"

// NewCors ...
func NewCors(origin string) Cors {
	return Cors{
		Origin:  origin, // TODO make this a regex and if it matches the given origin set the origin provided
		Methods: "POST, GET, OPTIONS, PUT, DELETE",
		Headers: "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}
}

// Cors helper to deal with headers
type Cors struct {
	Origin  string
	Methods string
	Headers string
}

// Wrap support
func (c Cors) Wrap(wrapped func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", c.Methods)
			w.Header().Set("Access-Control-Allow-Headers", c.Headers)
		}
		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}
		wrapped(w, r)
	}
}
