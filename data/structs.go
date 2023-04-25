package data

import (
	"time"
)

const RecordingLength int = 5

type Stream struct {
	Name     string
	Endpoint string
	BaseURL  string
	Live     bool
}

type StreamState string

const (
	Creating StreamState = "Creating"
	Ready    StreamState = "Ready"
)

type Recording struct {
	ID        string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	State     StreamState
}

type Instruction int

const Create Instruction = 0
const Update Instruction = 1

type RecordingInstruction struct {
	Instruction Instruction
	Recording   Recording
}
