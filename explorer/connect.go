package explorer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Connect ...
// input: Connection
// output: rets.CapabilityURLS
func Connect(conns map[string]Connection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Connection
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Printf("params: %v\n", p)
		ctx := context.Background()
		_, err = p.Login(ctx)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		conns[p.ID] = p
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p.URLs)

		JSONStore("/tmp/gorets/connections.json", &conns)

	}
}
