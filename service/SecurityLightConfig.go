package service

type SecurityLightConfig struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Sensors   []string `json:"sensors"`
	Lights    []string `json:"lights"`
	Timeout   int      `json:"timeout"`             // In minutes
	TimeStart string   `json:"timeStart,omitempty"` // If missing, always on
	TimeEnd   string   `json:"timeEnd,omitempty"`
	Enabled   bool     `json:"enabled"`
}
