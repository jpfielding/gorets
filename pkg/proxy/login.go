package proxy

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/jpfielding/gorets/pkg/rets"
)

// Login manages de/multiplexing requests to RETS servers
func Login(ops map[string]string, srcs Sources) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, ops["Login"])
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
		fmt.Fprintf(res, asXML(ops, src, usr, *urls))
	}
}

func asXML(ops map[string]string, src, usr string, urls rets.CapabilityURLs) string {
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
	// list our currently supported operations
	for o, prefix := range ops {
		prefix = strings.TrimSuffix(prefix, "/")
		ifPrint(o, fmt.Sprintf("%s/%s/%s", prefix, src, usr))
	}
	// TODO the source additional urls??
	// for k, v := range urls.AdditionalURLs {
	// 	ifPrint(k, v)
	// }

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
