package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/UniversityRadioYork/stream-recorder/data"
	"github.com/google/uuid"
)

func buildRecording(id uuid.UUID, startTime time.Time, endTime time.Time, stream data.Stream, recordingsChannel chan<- data.RecordingInstruction) {
	var filesToJoin []string

	filepath.WalkDir("recordings", func(path string, de fs.DirEntry, _ error) error {

		if de.IsDir() {
			return nil
		}

		reducedPath := strings.Split(path, "/")[1]

		splitFilePath := strings.Split(reducedPath, ".")

		if len(splitFilePath) != 3 || splitFilePath[2] != "mp3" || splitFilePath[0] != stream.Endpoint {
			return nil
		}

		fileTime, err := strconv.Atoi(splitFilePath[1])
		if err != nil {
			return err
		}

		if fileTime < int(startTime.Add(time.Duration(-1*data.RecordingLength)*time.Minute).Unix()) || fileTime > int(endTime.Add(time.Duration(data.RecordingLength)*time.Minute).Unix()) {
			return nil
		}

		filesToJoin = append(filesToJoin, path)

		return nil
	})

	sort.Strings(filesToJoin)

	fileStartTimeUnix, _ := strconv.Atoi(strings.Split(filesToJoin[0], ".")[1])
	fileStartTime := time.Unix(int64(fileStartTimeUnix), 0)

	fileEndTimeUnix, _ := strconv.Atoi(strings.Split(filesToJoin[len(filesToJoin)-1], ".")[1])
	fileEndTime := time.Unix(int64(fileEndTimeUnix), 0).Add(time.Duration(data.RecordingLength) * time.Minute)

	// TODO all the errors around here

	outFile, _ := os.Create(fmt.Sprintf("recordings/%v.mp3", id.String()))
	defer outFile.Close()
	for _, clip := range filesToJoin {
		clipFile, _ := os.Open(clip)
		io.Copy(outFile, clipFile)
		clipFile.Close()
	}

	// TODO updated start times
	recordingsChannel <- data.RecordingInstruction{
		Instruction: data.Update,
		Recording: data.Recording{
			ID:        id.String(),
			StartTime: fileStartTime,
			EndTime:   fileEndTime,
			State:     data.Ready,
		},
	}

}

func RequestRecording(name string, startTime time.Time, endTime time.Time, stream data.Stream, recordingsChannel chan<- data.RecordingInstruction) uuid.UUID {
	requestId := uuid.New()

	go buildRecording(requestId, startTime, endTime, stream, recordingsChannel)

	recordingsChannel <- data.RecordingInstruction{
		Instruction: data.Create,
		Recording: data.Recording{
			ID:        requestId.String(),
			Name:      name,
			StartTime: startTime,
			EndTime:   endTime,
			State:     data.Creating,
		},
	}

	return requestId

}
