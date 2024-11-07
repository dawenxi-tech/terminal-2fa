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
}

func newApplication() *Application {
	app := &Application{
		term: tview.NewApplication().EnableMouse(false),
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
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	err := app.renderTable()
	if err != nil {
		return err
	}

	//grid := tview.NewGrid().
	//	SetRows(2).
	//	SetColumns(1).
	//	AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
	//	AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newPrimitive("2FA"), 1, 1, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewButton("Add"), 0, 1, false).
			AddItem(tview.NewButton("Delete"), 0, 1, false),
			1, 1, false).
		AddItem(app.table, 0, 100, false)
	//AddItem(tview.NewBox(), 0, 1, false)

	app.term.SetRoot(flex, true)
	app.term.SetFocus(flex)
	return nil
}

func (app *Application) renderTable() error {
	app.table = tview.NewTable()
	app.table.SetBackgroundColor(tcell.ColorDefault)
	app.table.SetEvaluateAllRows(true)
	app.table.SetTitle("2FA")
	app.table.SetBorders(true)
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
		app.table.SetCell(i+1, col, cell)
		// name
		col++
		nameCell := tview.NewTableCell(obj.Name).SetExpansion(500)
		nameCell.SetAlign(tview.AlignCenter)
		app.table.SetCell(i+1, col, nameCell)
		// code
		col++
		codeCell := tview.NewTableCell("000000").SetExpansion(500)
		codeCell.SetAlign(tview.AlignCenter)
		app.table.SetCell(i+1, col, codeCell)

		// create
		col++
		createCell := tview.NewTableCell("2006-01-02 15:04:05").SetExpansion(500)
		createCell.SetAlign(tview.AlignCenter)
		app.table.SetCell(i+1, col, createCell)

		// action
		col++
		actionCell := tview.NewTableCell("delete, edit").SetExpansion(500)
		actionCell.SetAlign(tview.AlignCenter)
		app.table.SetCell(i+1, col, actionCell)
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
	app.table.SetCell(0, col, cell)
	// name
	col++
	nameCell := tview.NewTableCell("Name").SetExpansion(500)
	nameCell.SetAlign(tview.AlignCenter)
	app.table.SetCell(0, col, nameCell)
	// code
	col++
	codeCell := tview.NewTableCell("Code").SetExpansion(500)
	codeCell.SetAlign(tview.AlignCenter)
	app.table.SetCell(0, col, codeCell)

	// create
	col++
	createCell := tview.NewTableCell("Create Time").SetExpansion(500)
	createCell.SetAlign(tview.AlignCenter)
	app.table.SetCell(0, col, createCell)

	// action
	col++
	actionCell := tview.NewTableCell("Action").SetExpansion(500)
	actionCell.SetAlign(tview.AlignCenter)
	app.table.SetCell(0, col, actionCell)
}
