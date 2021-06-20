package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	d "github.com/UniversityRadioYork/stream-recorder/data"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO
	},
}

func (h *websocketH) websocketHandler(w http.ResponseWriter, r *http.Request, streams []*d.Stream) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to generate upgrader: %s", err)
		return
	}

	h.clients[ws] = true

	defer func() {
		delete(h.clients, ws)
		ws.Close()
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "1001") {
				log.Printf("Failed to read WebSocket message: %s", err)
			} else {
				// Client Disconnected
				return
			}
		}

		if string(message) == "QUERY" {
			for _, strm := range streams {
				h.PushUpdate(*strm, strm.Live)
			}
		}

	}
}

func (h *websocketH) PushUpdate(strm d.Stream, state bool) {
	for client := range h.clients {
		client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s:%v", strm.Name, strm.Live)))
	}
}
