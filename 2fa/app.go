package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"strconv"
)

type Application struct {
	term  *tview.Application
	table *tview.Table
	pages *tview.Pages

	inputDialog   *InputDialog
	confirmDialog *tview.Modal
}

func newApplication() *Application {
	app := &Application{
		term:          tview.NewApplication(),
		table:         tview.NewTable(),
		pages:         tview.NewPages(),
		inputDialog:   newInputDialog(""),
		confirmDialog: tview.NewModal(),
	}
	return app
}

func (app *Application) run() error {
	err := app.term.Run()
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to run application")
	}
	return err
}

func (app *Application) layout() error {
	newPrimitive := func(text string) tview.Primitive {
		view := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
		view.SetBackgroundColor(tcell.ColorDefault)
		return view
	}
	err := app.configTable()
	if err != nil {
		return err
	}
	app.configureAddDialog()

	info := app.helpMessage()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newPrimitive("2FA"), 1, 1, false).
		AddItem(app.table, 0, 100, true).
		AddItem(info, 1, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.term.Stop()
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'a' {
			//fmt.Println("show input page")
			app.inputDialog.clear()
			app.inputDialog.setTitle("ADD")
			app.pages.ShowPage("inputDialog")
			app.term.EnableMouse(true)
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'e' {
			//fmt.Println("show input page")
			app.pages.ShowPage("inputDialog")
			app.inputDialog.clear()
			app.inputDialog.setTitle("ADD")
			app.term.EnableMouse(true)
		}
		return event
	})

	//mod := modalFn(input, 40, 8)

	page := app.pages
	page.AddPage("home", flex, true, true)
	page.AddPage("inputDialog", app.inputDialog.getModal(), true, false)
	page.AddPage("deleteDialog", app.confirmDialog, true, false)

	app.term.SetRoot(page, true)
	app.term.SetFocus(page)
	return nil
}

func (app *Application) configureAddDialog() {
	app.inputDialog.nameField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			//name, code := app.inputDialog.values()
			app.term.EnableMouse(false)
		case tcell.KeyEsc:
			app.pages.HidePage("addDialog")
			app.term.EnableMouse(false)
		}
	})
}

func (app *Application) configTable() error {
	table := app.table
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetEvaluateAllRows(true)
	table.SetTitle("2FA")
	table.SetBorders(true)
	table.SetSelectable(true, false)
	table.SetBorderPadding(0, 0, 0, 0)
	table.SetSelectedStyle(tcell.Style{}.Background(tcell.ColorGreen))
	table.SetSelectedFunc(func(row, column int) {
		table.GetCell(row, column).SetBackgroundColor(tcell.ColorDefault)
	})
	objs, err := defaultStorage.readConfig()
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to read config")
		return err
	}
	// build header
	app.buildTableHeader()

	addCell := func(txt string, row int, col int, exp int) {
		cell := tview.NewTableCell(txt).SetExpansion(exp)
		cell.SetAlign(tview.AlignCenter)
		table.SetCell(row, col, cell)
	}

	for i, obj := range objs {
		// index
		var col int
		addCell(strconv.Itoa(i+1), i+1, col, 100)
		// name
		col++
		addCell(obj.Name, i+1, col, 500)
		// code
		col++
		addCell("000000", i+1, col, 500)
	}
	return nil
}

func (app *Application) buildTableHeader() {
	addCell := func(txt string, col int) {
		cell := tview.NewTableCell(txt).SetExpansion(100)
		cell.SetAlign(tview.AlignCenter)
		cell.NotSelectable = true
		app.table.SetCell(0, col, cell)
	}
	var col int
	addCell("#", col)
	// name
	col++
	addCell("Name", col)
	// code
	col++
	addCell("Code", col)
}

func (app *Application) helpMessage() tview.Primitive {
	tv := tview.NewTextView().SetText(`A: Add; E: Edit; D: Delete; +: Move Up; -: Move Down; Q: Quit`)
	tv.SetBackgroundColor(tcell.ColorDefault)
	return tv
}
