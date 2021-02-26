package components

import (
	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
)

type TimeEntryFormData struct {
	Title string
}

type TimeEntryForm struct {
	Title          *Input
	ActiveInput    int
	Data           TimeEntryFormData
	SubmitCallback func()

	inputs []*Input
}

func NewTimeEntryForm() *TimeEntryForm {
	f := &TimeEntryForm{
		Title:       NewInput(),
		ActiveInput: 0,
	}

	actionMapping := map[string]func(){
		"<Tab>":   f.FocusNext,
		"<Enter>": f.Submit,
	}

	f.Title.ActionMapping = actionMapping

	// Add all focusable inputs to inputs field.
	f.inputs = append(f.inputs, f.Title)

	return f
}

func (f *TimeEntryForm) Render() {
	// Setup help box.
	help := w.NewParagraph()
	help.Title = " Help "
	help.Text = "Navigation keys:\n" +
		"- <Tab> - Focus next input\n" +
		"- <Escape> - If form is focused left focus, else - close the form\n" +
		"- <Enter> - Submit the form"
	help.SetRect(0, 0, 100, 10)

	// Setup inputs.
	f.Title.SetRect(0, 10, 100, 13)
	f.Title.AddText(f.Data.Title)
	f.Title.Title = " Title "

	// Setup pseudo buttons.
	btnNegative := w.NewParagraph()
	btnNegative.Text = "Cancel <Escape>"
	btnNegative.TextStyle = ui.Style{
		Fg: ui.ColorBlack,
		Bg: ui.ColorYellow,
	}
	btnNegative.SetRect(0, 13, 17, 16)

	btnPositive := w.NewParagraph()
	btnPositive.Text = "Submit <Enter>"
	btnPositive.TextStyle = ui.Style{
		Fg: ui.ColorBlack,
		Bg: ui.ColorGreen,
	}
	btnPositive.SetRect(17, 13, 33, 16)

	// Render elements.
	ui.Render(help, f.Title, btnNegative, btnPositive)
	f.GetActiveInput().Capture()
}

func (f *TimeEntryForm) FocusNext() {
	f.ActiveInput++

	if f.ActiveInput >= len(f.inputs) {
		f.ActiveInput = 0
	}

	f.GetActiveInput().Capture()
}

func (f *TimeEntryForm) Submit() {
	f.SubmitCallback()
}

func (f *TimeEntryForm) GetActiveInput() *Input {
	return f.inputs[f.ActiveInput]
}

func (f *TimeEntryForm) GetFormData() TimeEntryFormData {
	return TimeEntryFormData{
		Title: f.Title.GetText(),
	}
}
