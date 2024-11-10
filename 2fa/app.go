package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"strconv"
	"time"
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
		inputDialog:   newInputDialog("Update"),
		confirmDialog: tview.NewModal(),
	}
	app.inputDialog.onCancel = app.onCancelInput
	app.inputDialog.onSubmit = app.onSaveInput
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
	err := app.configTable()
	if err != nil {
		return err
	}
	info := app.helpMessage()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newPrimitive("2FA"), 1, 1, false).
		AddItem(app.table, 0, 100, true).
		AddItem(info, 1, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			app.handleKeyPressed(event.Rune())
		}
		if event.Key() == tcell.KeyESC {
			app.term.Stop()
		}
		return event
	})
	page := app.pages
	page.AddPage("home", flex, true, true)
	page.AddPage("inputDialog", app.inputDialog.getModal(), true, false)
	page.AddPage("deleteDialog", app.confirmDialog, true, false)

	app.term.SetRoot(page, true)
	app.term.SetFocus(page)
	return nil
}

func (app *Application) handleKeyPressed(r rune) {
	switch r {
	case 'q':
		app.term.Stop()
	case 'a':
		app.inputDialog.clear()
		app.inputDialog.setTitle("ADD")
		app.pages.ShowPage("inputDialog")
		app.term.EnableMouse(true)
		app.inputDialog.focus()
	case 'e':
		row, _ := app.table.GetSelection()
		row = row - 1
		records, _ := defaultStorage.readConfig()
		if row < 0 || row >= len(records) {
			return
		}
		id, name, code := records[row].ID, records[row].Name, records[row].Seed
		app.inputDialog.setTitle("EDIT")
		app.pages.ShowPage("inputDialog")
		app.term.EnableMouse(true)
		app.inputDialog.update(id, name, code)
		app.inputDialog.focus()
	}
}

func (app *Application) reloadTable() {
	app.table.Clear()
	_ = app.configTable()
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

func (app *Application) onCancelInput() {
	app.dismissInputDialog()
}

func (app *Application) dismissInputDialog() {
	app.pages.ShowPage("home")
	app.pages.HidePage("inputDialog")
	app.term.EnableMouse(false)
}

func (app *Application) onSaveInput(id string, name, code string) {
	app.dismissInputDialog()
	records, _ := defaultStorage.readConfig()
	if id != "" {
		for i, record := range records {
			if record.ID == id {
				records[i].Name = name
				records[i].Seed = code
				break
			}
		}
	} else {
		records = append(records, Entry{
			ID:       newId(),
			Name:     name,
			Seed:     code,
			Order:    len(records),
			CreateAt: time.Now(),
		})
	}
	_ = defaultStorage.saveConfig(records)
	app.reloadTable()
}

func newPrimitive(text string) tview.Primitive {
	view := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
	view.SetBackgroundColor(tcell.ColorDefault)
	return view
}
