package tracker

import (
	"fmt"
	"os"
	"time"

	ai "github.com/devsapp/goutils/aigc/client"
)

// Client is a tracker client used for collecting traces and logs
type Client struct {
	sourceName string
	isEnabled  bool
}

// NewTrackerClient initializes a tracker client
func NewTrackerClient(sourceName string) *Client {
	return &Client{sourceName: sourceName}
}

// SetEnabled controls whether the client collects data or not
func (c *Client) SetEnabled(isEnabled bool) *Client {
	if c == nil {
		return &Client{isEnabled: isEnabled}
	}

	c.isEnabled = isEnabled
	return c
}

// SetSource changes the name of the tracking source
func (c *Client) SetSource(sourceName string) *Client {
	if c == nil {
		return &Client{sourceName: sourceName}
	}

	c.sourceName = sourceName
	return c
}

// SendLogs sends a series of logs
func (c *Client) SendLogs(logs []Log) error {
	if c == nil {
		return fmt.Errorf("tracker.Client is nil")
	}
	defer func() { recover() }()

	if !c.isEnabled {
		return nil
	}

	for i := 0; i < len(logs); i++ {
		if logs[i].AccountID == "" {
			logs[i].AccountID = os.Getenv("FC_ACCOUNT_ID")
		}
		if logs[i].Ts == 0 {
			logs[i].Ts = time.Now().Unix()
		}
		if logs[i].Source == "" {
			logs[i].Source = c.sourceName
		}
		if logs[i].Level == "" {
			logs[i].Level = "info"
		}
	}

	return ai.PostJSON(fmt.Sprintf("/collect/log"), logs)
}

// SendTrackers sends a series of trackers
func (c *Client) SendTrackers(trackers []Tracker) error {
	if c == nil {
		return fmt.Errorf("tracker.Client is nil")
	}
	defer func() { recover() }()

	if !c.isEnabled {
		return nil
	}

	for i := 0; i < len(trackers); i++ {
		if trackers[i].AccountID == "" {
			trackers[i].AccountID = os.Getenv("FC_ACCOUNT_ID")
		}
		if trackers[i].Ts == 0 {
			trackers[i].Ts = time.Now().Unix()
		}
		if trackers[i].Source == "" {
			trackers[i].Source = c.sourceName
		}
	}

	return ai.PostJSON("/collect/tracker", trackers)
}

// SendTracker will send one track
func (c *Client) SendTracker(key string, payload interface{}) {
	if c == nil {
		return
	}

	c.SendTrackers([]Tracker{{
		Key:     string(key),
		Payload: payload,
	}})
}
