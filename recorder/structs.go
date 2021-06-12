package recorder

import (
	"time"
)

type Stream struct {
	Name     string
	Endpoint string
	BaseURL  string
	Live     bool
}

type Recording struct {
	Filename   string    `yaml:"filename"`
	StreamName string    `yaml:"streamName"`
	StartTime  time.Time `yaml:"startTime"`
}
