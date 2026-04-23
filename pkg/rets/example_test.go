package rets

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func ExampleDefaultSession() {
	ctx := context.Background()

	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, `<RETS ReplyCode="0" ReplyText="Success">
<RETS-RESPONSE>
Search=/search
Logout=/logout
</RETS-RESPONSE>
</RETS>`)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, `<RETS ReplyCode="0" ReplyText="Logging out">
<RETS-RESPONSE>
SignOffMessage=Goodbye
</RETS-RESPONSE>
</RETS>`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	requester, err := DefaultSession(
		"username",  // RETS login username
		"password",  // RETS login password
		"MyApp/1.0", // User-Agent header value
		"",          // user-agent password (empty disables UA auth hashing)
		"RETS/1.8",  // RETS-Version header value
		nil,         // nil uses the default HTTP transport
	)
	if err != nil {
		panic(err)
	}

	urls, err := Login(ctx, requester, LoginRequest{URL: server.URL + "/login"})
	if err != nil {
		panic(err)
	}

	searchURL, err := url.Parse(urls.Search)
	if err != nil {
		panic(err)
	}
	fmt.Println("search path:", searchURL.Path)

	logout, err := Logout(ctx, requester, LogoutRequest{URL: urls.Logout})
	if err != nil {
		panic(err)
	}
	fmt.Println("logout reply:", logout.ReplyCode)

	// Output:
	// search path: /search
	// logout reply: 0
}
