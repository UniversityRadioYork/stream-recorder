package recorder

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	d "github.com/UniversityRadioYork/stream-recorder/data"
	"github.com/UniversityRadioYork/stream-recorder/web"
)

func RecordStream(stream *d.Stream, recordingsChannel chan<- d.Recording) {
	stream.Live = true
	go web.WebsocketMaster.PushUpdate(*stream, true)

	defer func() {
		stream.Live = false
		go web.WebsocketMaster.PushUpdate(*stream, false)
	}()

	log.Printf("Recording %s\n", stream.Name)

	strm := fmt.Sprintf("%s/%s", stream.BaseURL, stream.Endpoint)
	resp, err := http.Get(strm)
	if err != nil {
		log.Printf("Failed to GET stream %s (despite just getting it anyway): %s", strm, err)
		return
	}
	defer resp.Body.Close()

	startTime := time.Now()
	filename := fmt.Sprintf("recordings/%s%v.mp3", stream.Endpoint, startTime.Unix())
	recording, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed making recording file %s: %s", filename, err)
		return
	}
	defer recording.Close()

	_, err = io.Copy(recording, resp.Body)
	if err != nil {
		log.Printf("Failed copying response body to file: %s", err)
		return
	}

	log.Printf("Stopping Recording %s\n", stream.Name)

	recordingsChannel <- d.Recording{
		Filename:   filename,
		StreamName: stream.Name,
		StartTime:  startTime,
	}
}
