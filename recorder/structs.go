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
	Filename  string
	Stream    *Stream
	StartTime time.Time
}
