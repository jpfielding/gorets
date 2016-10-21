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

// Bytes ...
func (w WireLogPage) Bytes() []byte {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(w)
	return buf.Bytes()
}

// WirelogUpgrader ...
var WirelogUpgrader = websocket.Upgrader{
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: true,
	ReadBufferSize:    256,
	WriteBufferSize:   4096,
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
			// find the source for this message
			sess := sessions.Open(req.ID)
			// TODO look into keeping the file open
			wirelog, err := os.Open(sess.Wirelog())
			stat, _ := wirelog.Stat()
			// the noop message
			page := WireLogPage{
				ID:     req.ID,
				Offset: req.Offset,
				Length: 0,
				Size:   stat.Size(),
				Chunk:  "",
			}
			// offset over runs the file size
			if req.Offset < stat.Size() {
				// optionally set a starting point
				wirelog.Seek(req.Offset, 0)
				chunk := make([]byte, upgrader.WriteBufferSize-128)
				len, er := wirelog.Read(chunk)
				// if we read something... send it back
				if len != 0 && er != io.EOF {
					page.Length = len
					page.Chunk = string(chunk[0:len])
				}
			}
			wirelog.Close()
			if err = conn.WriteMessage(messageType, page.Bytes()); err != nil {
				return
			}
		}
	}
}
