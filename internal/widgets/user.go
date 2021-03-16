// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package widgets

import (
	"github.com/abgeo/goclockify/internal/types"
	w "github.com/gizak/termui/v3/widgets"
)

// UserWidget is a component that displays the user data
type UserWidget struct {
	*w.Table
	User types.User
}

// NewUserWidget creates new UserWidget
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

// SetUser sets the value of UserWidget.User
func (u *UserWidget) SetUser(user types.User) {
	u.User = user
	u.userToTable()
}

func (u *UserWidget) userToTable() {
	u.Rows[0][1] = u.User.Name
	u.Rows[1][1] = u.User.Email
}
