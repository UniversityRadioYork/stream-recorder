package recorder

import (
	"fmt"
	"io"
	"time"
)

var timeout = fmt.Errorf("timeout")

type writerWithEnd struct {
	standardResponse io.Reader
	endTime          time.Time
}

func (w *writerWithEnd) Read(p []byte) (n int, err error) {
	if time.Now().Before(w.endTime) {
		return w.standardResponse.Read(p)
	}

	// timed out
	return 0, timeout
}
