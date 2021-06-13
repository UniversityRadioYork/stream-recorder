package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/UniversityRadioYork/stream-recorder/recorder"
)

func StartWeb(port int, recordings *[]recorder.Recording) {

	fs := http.FileServer(http.Dir("frontend/build"))
	http.Handle("/", fs)

	http.HandleFunc("/recordings", func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := json.Marshal(recordings)
		if err != nil {
			panic(err) // TODO
		}
		w.Write(jsonData)
	})

	log.Printf("Listening on port %v", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
