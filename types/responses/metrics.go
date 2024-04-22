package responses

// Response from the Metrics endpoint
type MetricsResponse struct {
	ServerFPS        int     `json:"serverfps"`
	CurrentPlayerNum int     `json:"currentplayernum"`
	ServerFrameTime  float64 `json:"serverframetime"`
	MaxPlayerNum     int     `json:"maxplayernum"`
	ServerUptime     int     `json:"uptime"`
}
