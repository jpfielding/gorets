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
		// dont fill up the whole buffer
		buffer := make([]byte, upgrader.WriteBufferSize-128)
		for {
			// read in the next message
			messageType, msg, err := conn.NextReader()
			if err != nil {
				conn.WriteMessage(messageType, []byte(fmt.Sprintf("{'error':'%v'}", err)))
				continue
			}
			if messageType != websocket.TextMessage {
				conn.WriteMessage(messageType, []byte("{'error': 'wrong message type'}"))
				continue
			}
			req := WireLogPageRequest{}
			json.NewDecoder(msg).Decode(&req)
			log.Printf("wirelog request: %v\n", req)
			// find the source for this message
			sess := sessions.Open(req.ID)
			if sess == nil {
				conn.WriteMessage(messageType, []byte("{'error': 'id not found'}"))
				continue
			}
			err = sess.ReadWirelog(func(f *os.File, err error) error {
				if err != nil {
					return err
				}
				stat, _ := f.Stat()
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
					// set a starting point
					f.Seek(req.Offset, 0)
					len, er := f.Read(buffer)
					// if we read something... send it back
					if len != 0 && er != io.EOF {
						page.Length = len
						page.Chunk = string(buffer[0:len])
					}
				}
				return conn.WriteMessage(messageType, page.Bytes())
			})
			if err != nil {
				conn.WriteMessage(messageType, []byte(fmt.Sprintf("{'error':'%v'}", err)))
				continue
			}
		}
	}
}
