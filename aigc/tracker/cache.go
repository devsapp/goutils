package tracker

import (
	"fmt"
)

// LogCache type Log or Tracker
type LogCache interface {
	Size() int
}

// Log struct from AIGC server
type Log struct {
	AccountID string `json:"accountID"`
	Level     string `json:"level"`
	Ts        int64  `json:"ts"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestID"`
	Source    string `json:"source"`
}

// Size returns payload size of log
func (l *Log) Size() int { return len(l.Msg) }

// Tracker struct from AIGC server
type Tracker struct {
	Key       string      `json:"key"`
	AccountID string      `json:"accountID"`
	Ts        int64       `json:"ts"`
	Payload   interface{} `json:"payload"`
	Source    string      `json:"source"`
}

// Size returns payload size of tracker
func (l *Tracker) Size() int { return len(fmt.Sprint(l.Payload)) }
