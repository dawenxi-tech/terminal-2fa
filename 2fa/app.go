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
		term: tview.NewApplication(),
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

func (app *Application) showTable() {
	app.table = tview.NewTable()
	app.table.SetBackgroundColor(tcell.ColorDefault)
	app.table.SetEvaluateAllRows(true)
	app.table.SetTitle("2FA")
	app.table.SetBorders(true)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			cell := tview.NewTableCell(strconv.Itoa(i*10 + j)).SetExpansion(500)
			cell.SetAlign(tview.AlignCenter)
			app.table.SetCell(i, j, cell)
		}
	}
	app.term.SetRoot(app.table, true)
	app.term.SetFocus(app.table)
	//app.term.Draw()
}
