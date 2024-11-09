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
	pager *tview.Pages
}

func newApplication() *Application {
	app := &Application{
		term: tview.NewApplication().EnableMouse(true),
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
	err := app.renderTable()
	if err != nil {
		return err
	}

	info := app.helpMessage()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newPrimitive("2FA"), 1, 1, false).
		AddItem(app.table, 0, 100, true).
		AddItem(info, 1, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.term.Stop()
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'e' {
			//fmt.Println("show input page")
			app.pager.ShowPage("infobox")
		}
		return event
	})

	modalFn := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	input := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Name"), 1, 1, false).
		AddItem(tview.NewInputField(), 1, 1, true).
		AddItem(tview.NewBox(), 1, 1, true).
		AddItem(tview.NewTextView().SetText("Code"), 1, 1, false).
		AddItem(tview.NewInputField(), 1, 1, false)
	input.SetTitle("Add")
	input.SetBorder(true)

	mod := modalFn(input, 40, 8)

	page := tview.NewPages()
	app.pager = page
	page.AddPage("home", flex, true, true)
	page.AddPage("input", mod, true, false)

	app.term.SetRoot(page, true)
	app.term.SetFocus(page)
	return nil
}

func (app *Application) renderTable() error {
	table := tview.NewTable()
	app.table = table
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

	for i, obj := range objs {
		// index
		var col int
		cell := tview.NewTableCell(strconv.Itoa(i + 1)).SetExpansion(100)
		cell.SetAlign(tview.AlignCenter)
		table.SetCell(i+1, col, cell)
		// name
		col++
		nameCell := tview.NewTableCell(obj.Name).SetExpansion(500)
		nameCell.SetAlign(tview.AlignCenter)
		table.SetCell(i+1, col, nameCell)
		// code
		col++
		codeCell := tview.NewTableCell("000000").SetExpansion(500)
		codeCell.SetAlign(tview.AlignCenter)
		table.SetCell(i+1, col, codeCell)

		// create
		col++
		createCell := tview.NewTableCell("2006-01-02 15:04:05").SetExpansion(500)
		createCell.SetAlign(tview.AlignCenter)
		table.SetCell(i+1, col, createCell)

		// action
		col++
		actionCell := tview.NewTableCell("delete, edit").SetExpansion(500)
		actionCell.SetAlign(tview.AlignCenter)
		table.SetCell(i+1, col, actionCell)
	}

	//app.term.SetRoot(app.table, true)
	//app.term.SetFocus(app.table)
	//app.term.Draw()
	return nil
}

func (app *Application) buildTableHeader() {
	var col int
	cell := tview.NewTableCell("#").SetExpansion(100)
	cell.SetAlign(tview.AlignCenter)
	cell.NotSelectable = true
	app.table.SetCell(0, col, cell)
	// name
	col++
	nameCell := tview.NewTableCell("Name").SetExpansion(500)
	nameCell.SetAlign(tview.AlignCenter)
	nameCell.NotSelectable = true
	app.table.SetCell(0, col, nameCell)
	// code
	col++
	codeCell := tview.NewTableCell("Code").SetExpansion(500)
	codeCell.SetAlign(tview.AlignCenter)
	codeCell.NotSelectable = true
	app.table.SetCell(0, col, codeCell)

	// create
	col++
	createCell := tview.NewTableCell("Create Time").SetExpansion(500)
	createCell.SetAlign(tview.AlignCenter)
	createCell.NotSelectable = true
	app.table.SetCell(0, col, createCell)

	// action
	col++
	actionCell := tview.NewTableCell("Action").SetExpansion(500)
	actionCell.SetAlign(tview.AlignCenter)
	actionCell.NotSelectable = true
	app.table.SetCell(0, col, actionCell)
}

func (app *Application) helpMessage() tview.Primitive {
	tv := tview.NewTextView().SetText(`A: Add; E: Edit; D: Delete; +: Move Up; -: Move Down; Q: Quit`)
	tv.SetBackgroundColor(tcell.ColorDefault)
	return tv
}
