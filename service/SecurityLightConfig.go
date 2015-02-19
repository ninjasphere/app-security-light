package service

import "github.com/ninjasphere/go-ninja/suit"

type SecurityLightConfig struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Sensors []string `json:"sensors"`
	Lights  []string `json:"lights"`
	Timeout int      `json:"timeout"` // In minutes
	Time    suit.TimeRange
	Enabled bool `json:"enabled"`
}
