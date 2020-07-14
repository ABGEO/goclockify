package widgets

import (
	ui "github.com/gizak/termui/v3"
	"strconv"
)

type TimeInterval struct {
	Start    string
	End      string
	Duration string
}

type TimeEntry struct {
	ID           string
	Description  string
	TagIds       []string
	Billable     bool
	ProjectId    string
	TimeInterval TimeInterval
}

type TimeEntriesWidget struct {
	*Table
	TimeEntries []TimeEntry
}

func NewTimeEntriesWidget() *TimeEntriesWidget {
	self := &TimeEntriesWidget{
		Table: NewTable(),
	}

	self.Title = " Time Entries "
	self.ShowCursor = true
	self.ShowLocation = true
	self.ColGap = 3
	self.PadLeft = 2
	self.UniqueCol = 1
	self.Header = []string{"Description", "Duration", "Billable"}
	self.ColResizer = func() {
		self.ColWidths = []int{ui.MaxInt(self.Inner.Dx()-26, 10), 10, 5}
	}

	return self
}

func (self *TimeEntriesWidget) SetTimeEntries(timeEntries []TimeEntry) {
	self.TimeEntries = timeEntries
	self.entriesToRows()
}

func (self *TimeEntriesWidget) entriesToRows() {
	var timeEntries *[]TimeEntry
	timeEntries = &self.TimeEntries
	strings := make([][]string, len(*timeEntries))
	for i := range *timeEntries {
		strings[i] = make([]string, 3)
		strings[i][0] = (*timeEntries)[i].Description
		strings[i][1] = (*timeEntries)[i].TimeInterval.Duration
		strings[i][2] = strconv.FormatBool((*timeEntries)[i].Billable)
	}
	self.Rows = strings
}
