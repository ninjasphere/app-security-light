package service

type SecurityLightConfig struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Sensors []string `json:"sensors"`
	Lights  []string `json:"lights"`
	Timeout int      `json:"timeout"` // In minutes
	Time    struct {
		From string `json:"from,omitempty"` // If missing, always on
		To   string `json:"to,omitempty"`
	}
	Enabled bool `json:"enabled"`
}
