package data

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
	Filename   string    `yaml:"filename" json:"filename"`
	StreamName string    `yaml:"streamName" json:"streamName"`
	StartTime  time.Time `yaml:"startTime" json:"startTime"`
}
