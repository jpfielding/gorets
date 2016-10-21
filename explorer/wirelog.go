package explorer

import (
	"bytes"
	"encoding/json"
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
		var sess *Session
		var wirelog *os.File
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
			// load up the file we are streaming
			if sess == nil {
				sess = sessions[req.ID]
				wirelog, err = os.Open(sess.Wirelog())
				defer wirelog.Close()
			}
			// optionally set a starting point
			if req.Offset > 0 {
				wirelog.Seek(req.Offset, 0)
			}
			stat, _ := wirelog.Stat()
			chunk := make([]byte, upgrader.WriteBufferSize/2)
			len, err := wirelog.Read(chunk)
			if len == 0 || err == io.EOF {
				return
			}
			// figure out which part of which file to send back
			page := WireLogPage{
				ID:     req.ID,
				Offset: req.Offset + int64(len),
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
