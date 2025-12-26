package report

//go:generate enumer -type ReportType -text -transform=lower

type ReportType uint8

const (
	AutoCheck ReportType = iota
	Custom
	Gitlab
	Json
	Markdown
	Nist
	Oud4
	Owasp
	Owaspm
	Pcidss
	PlainReport
	Sans
	Sarif
	Xml
)
