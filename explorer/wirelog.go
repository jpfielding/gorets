package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// WireLogPageRequest ...
type WireLogPageRequest struct {
	ID     string `json:"id"`
	Offset int64  `json:"offset"`
}

// WireLogPage ...
type WireLogPage struct {
	ID     string `json:"id"`
	Offset int64  `json:"offset"`
	Length int    `json:"length"`
	Size   int64  `json:"size"`
	Chunk  string `json:"chunk"`
}

// DefaultUpgrader ...
var DefaultUpgrader = websocket.Upgrader{
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: true,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
}

// WireLogSocket provides access to wire logs as websocket data
func WireLogSocket(upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		for {
			// read in the next message
			messageType, msg, err := conn.NextReader()
			if err != nil {
				return
			}
			if messageType != websocket.TextMessage {
				log.Println("wrong message type")
				return
			}
			req := WireLogPageRequest{}
			json.NewDecoder(msg).Decode(&req)
			fmt.Printf("wirelog request: %v\n", req)
			sess := sessions.Open(req.ID)
			// TODO look into keeping the file open
			// load up the file we are streaming
			wirelog, err := os.Open(sess.Wirelog())
			// optionally set a starting point
			if req.Offset > 0 {
				wirelog.Seek(req.Offset, 0)
			}
			stat, _ := wirelog.Stat()
			chunk := make([]byte, upgrader.WriteBufferSize)
			len, err := wirelog.Read(chunk)
			if len == 0 || err == io.EOF {
				return
			}
			wirelog.Close()
			// figure out which part of which file to send back
			page := WireLogPage{
				ID:     req.ID,
				Offset: req.Offset,
				Length: len,
				Size:   stat.Size(),
				Chunk:  string(chunk[0:len]),
			}
			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(page)
			if err = conn.WriteMessage(messageType, buf.Bytes()); err != nil {
				return
			}
		}
	}
}
