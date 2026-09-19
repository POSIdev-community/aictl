package statistic

type Statistic struct {
	Total        int32  `json:"total"`
	High         int32  `json:"high"`
	Low          int32  `json:"low"`
	Medium       int32  `json:"medium"`
	Potential    int32  `json:"potential"`
	FilesTotal   int32  `json:"filesTotal"`
	FilesScanned int32  `json:"filesScanned"`
	UrlsTotal    int32  `json:"urlsTotal"`
	UrlsScanned  int32  `json:"urlsScanned"`
	ScanDuration string `json:"scanDuration"`
	PolicyState  string `json:"policyState"`
}
