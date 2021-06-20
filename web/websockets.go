package web

import (
	"fmt"
	"log"
	"net/http"

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
	defer ws.Close()

	h.clients = append(h.clients, ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Failed to read WebSocket message: %s", err)
		}

		if string(message) == "QUERY" {
			for _, strm := range streams {
				h.PushUpdate(*strm, strm.Live)
			}
		}

	}
}

func (h *websocketH) PushUpdate(strm d.Stream, state bool) {
	for _, client := range h.clients {
		client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s:%v", strm.Name, strm.Live)))
	}
}
