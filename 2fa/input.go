package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputDialog struct {
	title    string
	form     *tview.Form
	onSubmit func(id string, name, code string)
	onCancel func()

	model struct {
		id   string
		name string
		code string
	}

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
	dialog.Dialog.close = func() {
		if dialog.onCancel != nil {
			dialog.onCancel()
		}
	}
	return dialog.modal
}

func (dialog *InputDialog) setTitle(title string) {
	dialog.form.SetTitle(title)
}

func (dialog *InputDialog) onKeyPress(field *tview.InputField, key tcell.Key) {
}

func (dialog *InputDialog) clear() {
	dialog.update("", "", "")
}

func (dialog *InputDialog) values() (name string, code string) {
	name = dialog.model.name
	code = dialog.model.code
	return
}

func (dialog *InputDialog) update(id string, name string, code string) {
	dialog.model.name = id
	dialog.model.name = name
	dialog.model.code = code
	dialog.form.GetFormItemByLabel("Name").(*tview.InputField).SetText(name)
	dialog.form.GetFormItemByLabel("Code").(*tview.InputField).SetText(code)
}

func (dialog *InputDialog) layoutForm() {
	form := dialog.form
	form.AddInputField("Name", "", 36, nil, func(text string) {
		dialog.model.name = text
	}).AddInputField("Code", "", 36, nil, func(text string) {
		dialog.model.code = text
	}).AddButton("Save", func() {
		name, code := dialog.values()
		if name == "" || code == "" {
			return
		}
		if dialog.onSubmit != nil {
			dialog.onSubmit(dialog.model.id, name, code)
		}
	}).AddButton("Cancel", func() {
		if dialog.onCancel != nil {
			dialog.onCancel()
		}
	})
	form.SetBorder(true)
	form.SetTitle(dialog.title)
	form.SetButtonsAlign(tview.AlignCenter)
}
