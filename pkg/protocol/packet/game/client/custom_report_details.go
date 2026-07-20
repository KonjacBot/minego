package client

//codec:gen
type ReportDetails struct {
	Title       string
	Description string
}

type CustomReportDetails struct {
	Details map[string]string
}
