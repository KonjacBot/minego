package client

//codec:gen
type ReportDetails struct {
	Title       string
	Description string
}

//codec:gen
type CustomReportDetails struct {
	Details []ReportDetails
}
