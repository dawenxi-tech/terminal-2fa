package main

import (
	"flag"
	"fmt"
	"log/slog"
)

func main() {
	showUI := flag.Bool("ui", false, "show settings ui")
	flag.Parse()
	if *showUI {
		showTermUI()
		return
	}
	display2FA()
}

func showTermUI() {
	app := newApplication()
	app.showTable()
	err := app.run()
	if err != nil {
		slog.Error("error to run application")
	}
}

func display2FA() {
	fmt.Println("2fa")
}
