package widgets

import (
	w "github.com/gizak/termui/v3/widgets"
)

type User struct {
	ID               string
	Email            string
	Name             string
	ProfilePicture   string
	ActiveWorkspace  string
	DefaultWorkspace string
}

type UserWidget struct {
	*w.Table
	User User
}

func NewUserWidget() *UserWidget {
	self := &UserWidget{
		Table: w.NewTable(),
	}

	self.Title = " Current User "
	self.FillRow = true
	self.ColumnWidths = []int{10, -1}
	self.Rows = [][]string{
		{"Name", ""},
		{"Email", ""},
	}

	return self
}

func (self *UserWidget) SetUser(user User) {
	self.User = user
	self.userToTable()
}

func (self *UserWidget) userToTable() {
	self.Rows[0][1] = self.User.Name
	self.Rows[1][1] = self.User.Email
}
