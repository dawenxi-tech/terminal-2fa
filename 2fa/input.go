package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func buildModal(pri tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(pri, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

type InputDialog struct {
	title     string
	nameField *tview.InputField
	codeField *tview.InputField

	view  *tview.Flex
	modal tview.Primitive
}

func newInputDialog(title string) *InputDialog {
	dialog := &InputDialog{
		title:     title,
		nameField: tview.NewInputField(),
		codeField: tview.NewInputField(),
	}
	dialog.layout()
	return dialog
}

func (dialog *InputDialog) getModal() tview.Primitive {
	if dialog.modal != nil {
		return dialog.modal
	}
	dialog.modal = buildModal(dialog.view, 40, 8)
	return dialog.modal
}

func (dialog *InputDialog) setTitle(title string) {
	dialog.view.SetTitle(title)
}

func (dialog *InputDialog) clear() {
	dialog.update("", "")
}

func (dialog *InputDialog) values() (name string, code string) {
	name = dialog.nameField.GetText()
	code = dialog.codeField.GetText()
	return
}

func (dialog *InputDialog) update(name string, code string) {
	dialog.nameField.SetText(name)
	dialog.codeField.SetText(code)
}

func (dialog *InputDialog) onDone(fn func(name string, code string)) {
	dialog.nameField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter && fn != nil {
			fn(dialog.nameField.GetText(), dialog.codeField.GetText())
		}
	})
	dialog.codeField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter && fn != nil {
			fn(dialog.nameField.GetText(), dialog.codeField.GetText())
		}
	})
}

func (dialog *InputDialog) layout() tview.Primitive {
	input := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Name"), 1, 1, false).
		AddItem(dialog.nameField, 1, 1, true).
		AddItem(tview.NewBox(), 1, 1, true).
		AddItem(tview.NewTextView().SetText("Code"), 1, 1, false).
		AddItem(dialog.codeField, 1, 1, false)
	input.SetTitle(dialog.title)
	input.SetBorder(true)
	dialog.view = input
	return input
}
