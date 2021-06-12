package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// "github.com/UniversityRadioYork/stream-recorder/web"
	"github.com/UniversityRadioYork/stream-recorder/recorder"
)

func main() {
	log.Println("Stream Recorder")

	// web.StartWeb()

	var recordings []recorder.Recording
	recordingsChannel := make(chan recorder.Recording)

	streams := []recorder.Stream{
		{
			Name:     "OB-Line",
			BaseURL:  "https://audio.ury.org.uk",
			Endpoint: "OB-Line",
		},
		{
			Name:     "OB-Line2",
			BaseURL:  "https://audio.ury.org.uk",
			Endpoint: "OB-Line2",
		},
	}

	go func() {
		for {
			for idx, stream := range streams {
				if !stream.Live {
					res, err := http.Get(fmt.Sprintf("%s/%s", stream.BaseURL, stream.Endpoint))

					if err != nil {
						panic(err) // TODO
					}

					if res.StatusCode == 200 {
						go recorder.RecordStream(&streams[idx], recordingsChannel)
					} else {
						fmt.Printf("%s Not Live\n", stream.Name)
					}
				}
			}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()

	for msg := range recordingsChannel {
		recordings = append(recordings, msg)
	}

}
