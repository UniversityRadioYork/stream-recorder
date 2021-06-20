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
		log.Panicf("Error reading YAML config file: %s", err)
	}

	var config config
	err = yaml.Unmarshal(configYamlFile, &config)

	if err != nil {
		log.Panicf("Error decoding YAML config: %s", err)
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
			log.Panicf("Failed reading recordings file: %s", err)
		}
	}

	err = yaml.Unmarshal(recordingsYamlFile, &recordings)

	if err != nil {
		log.Panicf("Error decoding YAML recordings file: %s", err)
	}

	go func() {
		for {
			for _, stream := range streams {
				if !stream.Live {
					strm := fmt.Sprintf("%s/%s", stream.BaseURL, stream.Endpoint)
					res, err := http.Get(strm)

					if err != nil {
						log.Printf("Failed to GET stream %s: %s\n", strm, err)
					} else {
						if res.StatusCode == 200 {
							go recorder.RecordStream(stream, recordingsChannel)
						}
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
				log.Printf("Failed to marshal YAML for all recordings: %s\n", err)
			} else {
				err = os.WriteFile("recordings.yml", yamlData, 0222)
				if err != nil {
					log.Printf("Failed to write YAML recordings file: %s\n", err)
				}
			}
		}
	}()

	web.StartWeb(config.WebPort, &recordings, streams)

}
