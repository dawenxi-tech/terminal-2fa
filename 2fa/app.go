package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xlzd/gotp"
	"slices"
	"strconv"
	"time"
)

type Application struct {
	term  *tview.Application
	table *tview.Table
	pages *tview.Pages

	manual *tview.TextView

	inputDialog   *InputDialog
	importDialog  *ImportDialog
	confirmDialog *tview.Modal
	errorDialog   *tview.Modal
}

func newApplication() *Application {
	app := &Application{
		term:          tview.NewApplication(),
		table:         tview.NewTable(),
		pages:         tview.NewPages(),
		inputDialog:   newInputDialog(""),
		confirmDialog: tview.NewModal(),
		errorDialog:   tview.NewModal(),
		importDialog:  newImportDialog(),
	}
	app.inputDialog.onCancel = app.onCancelInput
	app.inputDialog.onSubmit = app.onSaveInput
	app.importDialog.onCancel = app.onCancelImport
	app.importDialog.onSubmit = app.onSubmitImport
	app.confirmDialog.SetTitle("Confirm").SetBorder(true)
	app.errorDialog.SetTitle("Error").SetBorder(true)
	return app
}

func (app *Application) run() error {
	err := app.term.Run()
	if err != nil {
		fmt.Println("error to run", err)
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
			app.exit()
		}
		return event
	})
	page := app.pages
	page.AddPage("home", flex, true, true)
	page.AddPage("inputDialog", app.inputDialog.getModal(), true, false)
	page.AddPage("deleteDialog", app.confirmDialog, true, false)
	page.AddPage("importDialog", app.importDialog.getModal(), true, false)
	page.AddPage("errorDialog", app.errorDialog, true, false)

	app.term.SetRoot(page, true)
	app.term.SetFocus(page)
	return nil
}

func (app *Application) exit() {
	app.term.Stop()
}

func (app *Application) handleKeyPressed(r rune) {
	switch r {
	case 'q':
		app.exit()
	case 'a':
		app.inputDialog.clear()
		app.inputDialog.setTitle("ADD")
		app.pages.ShowPage("inputDialog")
		app.term.EnableMouse(true)
		app.inputDialog.focus()
	case 'e':
		row, _ := app.table.GetSelection()
		row = row - 1
		records, _ := defaultStorage.LoadRecords()
		if row < 0 || row >= len(records) {
			return
		}
		id, name, code := records[row].ID, records[row].Name, records[row].Secret.Val()
		app.inputDialog.setTitle("EDIT")
		app.pages.ShowPage("inputDialog")
		app.term.EnableMouse(true)
		app.inputDialog.update(id, name, code)
		app.inputDialog.focus()
	case '+':
		// move up
		row, _ := app.table.GetSelection()
		row = row - 1
		_ = defaultStorage.Move(row, "", -1)
		app.reloadTable()
		app.table.Select(row, 0)
	case '-':
		// move down
		row, _ := app.table.GetSelection()
		row = row - 1
		_ = defaultStorage.Move(row, "", 1)
		app.reloadTable()
		app.table.Select(row+2, 0)
	case 'd':
		app.showDeleteDialog()
	case 'i':
		// import
		app.importDialog.clear()
		app.pages.ShowPage("importDialog")
		app.term.EnableMouse(true)
	}
}

func (app *Application) showErrorMessage(msg string) {
	app.errorDialog.SetText(msg)
	app.errorDialog.
		ClearButtons().
		AddButtons([]string{"ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.pages.HidePage("errorDialog")
		app.term.EnableMouse(false)
	})
	app.pages.ShowPage("errorDialog")
	app.term.EnableMouse(true)
}

func (app *Application) reloadTable() {
	app.table.Clear()
	_ = app.configTable()
}

func (app *Application) showDeleteDialog() {
	row, _ := app.table.GetSelection()
	row = row - 1
	records, _ := defaultStorage.LoadRecords()
	if row < 0 || row >= len(records) {
		return
	}
	record := records[row]
	msg := fmt.Sprintf("Are you sure you want to delete \"%s\"?", record.Name)
	app.confirmDialog.SetText(msg)
	app.confirmDialog.
		ClearButtons().
		AddButtons([]string{"Cancel", "Delete"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Delete":
			_ = defaultStorage.DeleteRecord(row, "")
			app.reloadTable()
			app.table.Select(row, 0)
		case "Cancel":
		}
		app.pages.HidePage("deleteDialog")
		app.term.EnableMouse(false)
	})
	app.pages.ShowPage("deleteDialog")
	app.term.EnableMouse(true)
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
	objs, err := defaultStorage.LoadRecords()
	if err != nil {
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
		code := gotp.NewDefaultTOTP(obj.Secret.Val()).At(time.Now().Unix())
		addCell(code, i+1, col, 500)
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
	app.manual = tv
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
	if !isValidTOTPCode(code) {
		app.showErrorMessage("invalid 2fa secret")
		return
	}
	app.dismissInputDialog()
	if id != "" {
		records, _ := defaultStorage.LoadRecords()
		idx := slices.IndexFunc(records, func(entry Entry) bool {
			return entry.ID == id
		})
		_ = defaultStorage.Update(idx, name, code)
	} else {
		_ = defaultStorage.SaveEntry(Entry{
			ID:       newId(),
			Name:     name,
			Secret:   newSecret(code),
			CreateAt: time.Now(),
		})
	}
	app.reloadTable()
}

func (app *Application) onCancelImport() {
	app.pages.HidePage("importDialog")
	app.term.EnableMouse(false)
}

func (app *Application) onSubmitImport(uri string) {
	err := defaultStorage.Import(uri)
	if err != nil {
		app.showErrorMessage(err.Error())
		return
	}
	app.reloadTable()
	app.pages.HidePage("importDialog")
	app.term.EnableMouse(false)
}

func newPrimitive(text string) tview.Primitive {
	view := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
	view.SetBackgroundColor(tcell.ColorDefault)
	return view
}
