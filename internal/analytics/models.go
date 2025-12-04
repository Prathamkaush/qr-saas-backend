package analytics

type SummaryResponse struct {
	TotalScans int            `json:"total_scans"`
	Countries  map[string]int `json:"countries"`
	Devices    map[string]int `json:"devices"`
	Browsers   map[string]int `json:"browsers"`
}
type TimeSeriesPoint struct {
	Timestamp string `json:"timestamp"`
	Scans     int    `json:"scans"`
}
