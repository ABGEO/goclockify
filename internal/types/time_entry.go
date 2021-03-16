package types

// TimeEntry represents the time entry entity from the API
type TimeEntry struct {
	ID           string
	Description  string
	Tags         []Tag
	Billable     bool
	Project      Project
	TimeInterval TimeInterval
}
