package types

import "time"

// TimeInterval represents the time interval entity from the API
type TimeInterval struct {
	Start    time.Time
	End      time.Time
	Duration string
}
