package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/UniversityRadioYork/stream-recorder/recorder"
)

func StartWeb(port int, recordings *[]recorder.Recording) {

	webFS := http.FileServer(http.Dir("frontend/build"))
	http.Handle("/", webFS)

	recordingsFS := http.FileServer(http.Dir("recordings"))
	http.Handle("/recordings/", http.StripPrefix("/recordings/", recordingsFS))

	http.HandleFunc("/recordings-json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

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
