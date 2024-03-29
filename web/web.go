package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/UniversityRadioYork/stream-recorder/builder"
	d "github.com/UniversityRadioYork/stream-recorder/data"
	"github.com/gorilla/websocket"
)

type websocketH struct {
	clients map[*websocket.Conn]bool
}

var WebsocketMaster websocketH = websocketH{clients: make(map[*websocket.Conn]bool)}

func StartWeb(port int, recordings *[]d.Recording, streams []*d.Stream, recordingsChannel chan<- d.RecordingInstruction) {

	// webFS := http.FileServer(http.Dir("frontend/build"))
	// http.Handle("/", webFS)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html") })

	recordingsFS := http.FileServer(http.Dir("recordings"))
	http.Handle("/recordings/", http.StripPrefix("/recordings/", recordingsFS))

	http.HandleFunc("/recordings-json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		jsonData, err := json.Marshal(recordings)
		if err != nil {
			log.Printf("Failed to marshall JSON data for recordings: %s\n", err)
		} else {
			w.Write(jsonData)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { WebsocketMaster.websocketHandler(w, r, streams) })

	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s not allowed", r.Method)
			return
		}

		var bodyData map[string]string

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, err)
			return

		}
		err = json.Unmarshal(body, &bodyData)
		if err != nil {
			fmt.Fprint(w, err)
			return

		}

		st, err := time.Parse(time.RFC3339, bodyData["startTime"])
		if err != nil {
			fmt.Fprint(w, err)
			return

		}

		et, err := time.Parse(time.RFC3339, bodyData["endTime"])
		if err != nil {
			fmt.Fprint(w, err)
			return

		}

		var str d.Stream
		var found bool
		for _, v := range streams {
			if v.Endpoint == bodyData["stream"] {
				str = *v
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(w, "%s not found", bodyData["stream"])
			return
		}

		fmt.Fprint(w, builder.RequestRecording(bodyData["name"], st, et, str, recordingsChannel))
	})

	log.Printf("Listening on port %v\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
