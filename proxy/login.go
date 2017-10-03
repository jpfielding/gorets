package proxy

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/jpfielding/gorets/rets"
)

// Login manages de/multiplexing requests to RETS servers
func Login(prefix string, srcs Sources) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, prefix)
		parts := strings.Split(sub, "/")
		src := parts[0]
		usr := parts[1]
		if _, ok := srcs[src]; !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "source %s not found", src)
			return
		}
		if _, ok := srcs[src][usr]; !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "user %s not found", usr)
			return
		}
		session := srcs[src][usr]
		_, urls, err := session.Get()
		if err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(res, "source %s, user %s login failed", src, usr)
			return
		}
		// success, send the urls (modified to point to this server)
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, asXML(fmt.Sprintf("%s/%s/%s/login", prefix, src, usr), *urls))
	}
}

func asXML(prefix string, urls rets.CapabilityURLs) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("<RETS ReplyCode='%d' ReplyText='%s'>\n", urls.Response.Code, urls.Response.Text))
	fmt.Fprintln(&buf, "<RETS-RESPONSE>")
	// if theres a value, write it to the buffer
	ifPrint := func(name, value string) {
		if value == "" {
			return
		}
		fmt.Fprintf(&buf, "%s=%s\n", name, value)
	}
	// fmt.Fprint(&buf, "\n")
	ifPrint("Login", fmt.Sprintf("%s/login", prefix))
	ifPrint("Action", fmt.Sprintf("%s/action", prefix))
	ifPrint("Search", fmt.Sprintf("%s/search", prefix))
	ifPrint("Get", fmt.Sprintf("%s/get", prefix))
	ifPrint("GetObject", fmt.Sprintf("%s/getobject", prefix))
	ifPrint("Logout", fmt.Sprintf("%s/logout", prefix))
	ifPrint("GetMetadata", fmt.Sprintf("%s/getmetadata", prefix))
	for k, v := range urls.AdditionalURLs {
		ifPrint(k, v)
	}

	ifPrint("MemberName", urls.MemberName)
	ifPrint("User", urls.User)
	ifPrint("Broker", urls.Broker)
	ifPrint("MetadataVersion", urls.MetadataVersion)
	ifPrint("MinMetadataVersion", urls.MinMetadataVersion)
	ifPrint("OfficeList", strings.Join(urls.OfficeList, ","))
	if urls.TimeoutSeconds > 0 {
		ifPrint("TimeoutSeconds", fmt.Sprintf("%d", urls.TimeoutSeconds))
	}

	fmt.Fprintln(&buf, "</RETS-RESPONSE>")
	fmt.Fprintln(&buf, "</RETS>")

	return buf.String()
}
