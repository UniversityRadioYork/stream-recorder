package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	d "github.com/UniversityRadioYork/stream-recorder/data"
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

	var streams []*d.Stream
	var recordings []d.Recording = make([]d.Recording, 0)
	recordingsChannel := make(chan d.Recording)

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
		streams = append(streams, &d.Stream{
			Name:     stream.Name,
			Endpoint: stream.Endpoint,
			BaseURL:  config.BaseURL,
		})
	}

	recordingsYamlFile, err := os.ReadFile("recordings.yml")

	if err != nil {
		if !strings.Contains(err.Error(), "no such file") {
			panic(err) // TODO
		}
	}

	err = yaml.Unmarshal(recordingsYamlFile, &recordings)

	if err != nil {
		panic(err) // TODO
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
			yamlData, err := yaml.Marshal(recordings)
			if err != nil {
				panic(err)
			}
			os.WriteFile("recordings.yml", yamlData, 0222)
		}
	}()

	web.StartWeb(config.WebPort, &recordings, streams)

}
