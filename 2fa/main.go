package main

import (
	"flag"
	"fmt"
	"github.com/xlzd/gotp"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	config := flag.Bool("config", false, "show settings ui")
	flag.Parse()
	err := defaultStorage.init()
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to init storage")
		os.Exit(1)
	}
	if *config {
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
	now := time.Now().Unix()
	for i, obj := range objs {
		code := gotp.NewDefaultTOTP(obj.Seed).At(now)
		tw.AppendRow(table.Row{strconv.Itoa(i + 1), obj.Name + "\t", code + "\t", "\t10s"})
	}
	tw.SetCaption("Use -config to manager codes.")
	fmt.Println(tw.Render())
}
