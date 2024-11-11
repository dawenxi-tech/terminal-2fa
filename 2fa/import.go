package main

import (
	"github.com/dim13/otpauth/migration"
	"github.com/rivo/tview"
)

type ImportDialog struct {
	form     *tview.Form
	onSubmit func(uri string)
	onCancel func()
	title    string
	*Dialog

	value string
}

func newImportDialog() *ImportDialog {
	dialog := &ImportDialog{
		form: tview.NewForm(),
	}
	dialog.layoutForm()
	return dialog
}

func (dig *ImportDialog) getModal() tview.Primitive {
	if dig.Dialog != nil {
		return dig.Dialog.modal
	}
	dig.Dialog = NewDialog(dig.form, 50, 8)
	dig.Dialog.setClose(func() {
		if dig.onCancel != nil {
			dig.onCancel()
		}
	})
	return dig.Dialog.modal
}

func (dig *ImportDialog) clear() {
	dig.form.GetFormItemByLabel("URL").(*tview.InputField).SetText("")
}

func (dig *ImportDialog) layoutForm() {
	form := dig.form
	form.SetTitle("IMPORT")
	form.AddInputField("URL", "", 40, nil, func(text string) {
		dig.value = text
	}).AddButton("Cancel", func() {
		if dig.onCancel != nil {
			dig.onCancel()
		}
	}).AddButton("Submit", func() {
		_, err := migration.UnmarshalURL(dig.value)
		if err != nil {
			return
		}
		if dig.onSubmit != nil {
			dig.onSubmit(dig.value)
		}
	})
	form.SetBorder(true)
	form.SetButtonsAlign(tview.AlignCenter)
}
