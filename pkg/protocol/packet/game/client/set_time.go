package client

//codec:gen
type SetTime struct {
	WorldAge            int64
	TimeOfDay           int64
	TimeOfDayIncreasing bool
}
