package recorder

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func RecordStream(stream *Stream, recordingsChannel chan<- Recording) {
	stream.Live = true
	fmt.Printf("Recording %s\n", stream.Name)

	resp, err := http.Get(fmt.Sprintf("%s/%s", stream.BaseURL, stream.Endpoint))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	startTime := time.Now()
	filename := fmt.Sprintf("recordings/%s%v.mp3", stream.Endpoint, startTime.Unix())
	recording, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer recording.Close()

	_, err = io.Copy(recording, resp.Body)
	if err != nil {
		panic(err)
	}

	stream.Live = false
	fmt.Printf("Stopping Recording %s\n", stream.Name)

	recordingsChannel <- Recording{
		Filename:   filename,
		StreamName: stream.Name,
		StartTime:  startTime,
	}
}
