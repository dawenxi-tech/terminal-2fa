package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputDialog struct {
	title    string
	form     *tview.Form
	onSubmit func(name, code string)
	onCancel func()

	*Dialog
}

func newInputDialog(title string) *InputDialog {
	dialog := &InputDialog{
		title: title,
		form:  tview.NewForm(),
	}
	dialog.layoutForm()
	return dialog
}

func (dialog *InputDialog) getModal() tview.Primitive {
	if dialog.Dialog != nil {
		return dialog.modal
	}
	dialog.Dialog = NewDialog(dialog.form, 50, 10)
	return dialog.modal
}

func (dialog *InputDialog) setTitle(title string) {
	dialog.form.SetTitle(title)
}

func (dialog *InputDialog) onKeyPress(field *tview.InputField, key tcell.Key) {
}

func (dialog *InputDialog) clear() {
	dialog.form.Clear(true)
}

func (dialog *InputDialog) values() (name string, code string) {
	return
}

func (dialog *InputDialog) update(name string, code string) {
}

func (dialog *InputDialog) onDone(fn func(name string, code string)) {
}

func (dialog *InputDialog) layoutForm() {
	form := dialog.form
	form.AddInputField("Name", "", 36, nil, nil).
		AddInputField("Code", "", 36, nil, nil).
		AddButton("Save", func() {

		}).
		AddButton("Cancel", func() {
		})
	form.SetBorder(true)
	form.SetTitle(dialog.title)
	form.SetButtonsAlign(tview.AlignCenter)
}
