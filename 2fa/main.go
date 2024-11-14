package main

import (
	"flag"
	"fmt"
	"github.com/xlzd/gotp"
	"log/slog"
	"os"
	"strconv"
	"strings"
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
	for {
		str := render2FA(objs)
		_, _ = os.Stdout.WriteString(str)
		time.Sleep(time.Second)
		lineCount := len(strings.Split(str, "\n"))
		clearStr := strings.Repeat("\033[K\033[A", lineCount-1)
		_, _ = os.Stdout.WriteString(clearStr)
	}
}

func render2FA(objs []Entry) string {
	tw := table.NewWriter()
	tw.SetTitle("2FA")
	tw.SetIndexColumn(1)
	tw.AppendHeader(table.Row{"#", "Name\t\t", "Code\t\t", "Remain\t"})
	now := time.Now().Unix()
	for i, obj := range objs {
		code := gotp.NewDefaultTOTP(obj.Seed).At(now)
		tw.AppendRow(table.Row{strconv.Itoa(i + 1), " " + obj.Name + "\t", " " + code + "\t", fmt.Sprintf(" %ds", 30-time.Now().Second()%30)})
	}
	tw.SetCaption("Use -config to manager codes.")
	return "\n" + tw.Render()
}
