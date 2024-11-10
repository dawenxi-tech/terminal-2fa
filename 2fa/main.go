package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	showUI := flag.Bool("ui", false, "show settings ui")
	flag.Parse()
	err := defaultStorage.init()
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to init storage")
		os.Exit(1)
	}
	if *showUI {
		showTermUI()
		return
	}
	display2FA()
}

func showTermUI() {
	app := newApplication()
	_ = app.layout()
	err := app.run()
	if err != nil {
		slog.Error("error to run application")
	}
}

func display2FA() {
	objs, err := defaultStorage.readConfig()
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to read configure")
		os.Exit(1)
	}
	tw := table.NewWriter()
	tw.SetTitle("2FA")
	tw.SetIndexColumn(1)
	tw.AppendHeader(table.Row{"#", "Name\t\t", "Code\t\t", "Time\t\t"})
	for i, obj := range objs {
		tw.AppendRow(table.Row{strconv.Itoa(i + 1), obj.Name + "\t", "000000" + "\t", "\t10s"})
	}
	tw.SetCaption("Use -ui to manager codes.")
	fmt.Println(tw.Render())
}
