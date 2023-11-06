package log

import (
	"regexp"
	"strings"
	"sync"

	"github.com/devsapp/goutils/aigc/tracker"
	"github.com/sirupsen/logrus"
)

// LogrusFormatter implements logrus formatter to collect logs
type LogrusFormatter struct {
	client *tracker.Client

	lock      *sync.Mutex
	formatter logrus.Formatter
	cb        *tracker.Background
}

func NewLogrusFormatter(client *tracker.Client) *LogrusFormatter {
	return &LogrusFormatter{client: client}
}

// Format ...
func (f *LogrusFormatter) Format(log *logrus.Entry) ([]byte, error) {
	if f.cb == nil {
		f.cb = tracker.NewBackground(f.client)
	}
	if f.lock == nil {
		f.lock = new(sync.Mutex)
	}
	if f.formatter == nil {
		f.formatter = new(logrus.TextFormatter)
	}

	readableText := escapeUnexpectedChar(log.Message)

	if len(readableText) > 0 {
		rid := GetRid()

		f.cb.Push(&tracker.Log{
			Msg:       readableText,
			RequestID: rid,
			Level:     log.Level.String(),
		})
	}

	// log.Message = readableText
	// return (f.formatter).Format(log)
	if len(log.Message) > 0 && log.Message[len(log.Message)-1] != '\n' && log.Message[len(log.Message)-1] != '\r' {
		return append([]byte(log.Message), '\n'), nil
	}
	return []byte(log.Message), nil
}

// SendAllLogs will send all cache logs to server
func (f *LogrusFormatter) SendAllLogs() {
	f.cb.SendAll()
}

func escapeUnexpectedChar(s string) string {
	for _, r := range []string{
		"\\033\\[[0-9;]*m",
		"\\033\\[[\\d\\,]+[m]",
		"\\033\\[[\\d\\,]*[a-zA-Z]",
		"\\033\\[\\?25[a-zA-Z]",
		"\\033\\[\\d+;\\d+m",
		"\\033",
	} {
		s = regexp.MustCompile(r).ReplaceAllString(s, "")
	}

	return strings.Map(func(r rune) rune {
		switch r {
		case '\b', '\r', '\a', '\033':
			return -1
		default:
			return r
		}
	}, s)
}
