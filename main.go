package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/UniversityRadioYork/stream-recorder/recorder"
	"github.com/UniversityRadioYork/stream-recorder/web"
	"gopkg.in/yaml.v2"
)

type configStream struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

type config struct {
	WebPort int             `yaml:"port"`
	BaseURL string          `yaml:"baseURL"`
	Streams []*configStream `yaml:"streams"`
}

func main() {
	log.Println("Stream Recorder")

	var streams []*recorder.Stream
	var recordings []recorder.Recording = make([]recorder.Recording, 0)
	recordingsChannel := make(chan recorder.Recording)

	configYamlFile, err := os.ReadFile("config.yml")

	if err != nil {
		panic(err) // TODO
	}

	var config config
	err = yaml.Unmarshal(configYamlFile, &config)

	if err != nil {
		panic(err) // TODO
	}

	for _, stream := range config.Streams {
		streams = append(streams, &recorder.Stream{
			Name:     stream.Name,
			Endpoint: stream.Endpoint,
			BaseURL:  config.BaseURL,
		})
	}

	go func() {
		for {
			for _, stream := range streams {
				if !stream.Live {
					res, err := http.Get(fmt.Sprintf("%s/%s", stream.BaseURL, stream.Endpoint))

					if err != nil {
						panic(err) // TODO
					}

					if res.StatusCode == 200 {
						go recorder.RecordStream(stream, recordingsChannel)
					} else {
						fmt.Printf("%s Not Live\n", stream.Name)
					}
				}
			}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()

	go func() {
		for msg := range recordingsChannel {
			recordings = append(recordings, msg)
		}
	}()

	web.StartWeb(config.WebPort, &recordings)

}
