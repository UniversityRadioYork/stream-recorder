package web

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = 3000

func StartWeb() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	log.Printf("Listening on port %v", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}
