package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Command interface {
	exec()
}

type ConfigCommand struct {
	args    []string
	command Command
}

type addCommand struct {
	name   string
	secret string
}

func newAddCommand(args []string) addCommand {
	ac := addCommand{}
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	cmd.StringVar(&ac.name, "name", "", "code name")
	cmd.StringVar(&ac.secret, "secret", "", "secret")
	err := cmd.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ac
}

func (add addCommand) exec() {
	if add.name == "" || add.secret == "" {
		fmt.Println("usage: 2fa config add -name=foo -secret=bar")
		return
	}
	if !isValidTOTPCode(add.secret) {
		fmt.Println("invalid 2fa secret")
		return
	}
	err := defaultStorage.SaveEntry(Entry{
		ID:       newId(),
		Name:     add.name,
		Secret:   add.secret,
		CreateAt: time.Now(),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("success add 2fa")
}

type deleteCommand struct {
	id   int
	name string
}

func newDeleteCommand(args []string) deleteCommand {
	dc := deleteCommand{}
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	cmd.StringVar(&dc.name, "name", "", "code name")
	cmd.IntVar(&dc.id, "id", -1, "id")
	err := cmd.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return dc
}

func (d deleteCommand) exec() {
	err := defaultStorage.DeleteRecord(d.id-1, d.name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type editCommand struct {
	id     int
	name   string
	secret string
}

func newEditCommand(args []string) editCommand {
	ec := editCommand{}
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	cmd.StringVar(&ec.name, "name", "", "code name")
	cmd.StringVar(&ec.secret, "secret", "", "secret")
	cmd.IntVar(&ec.id, "id", 0, "id")
	err := cmd.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ec
}

func (e editCommand) exec() {
	err := defaultStorage.Update(e.id-1, e.name, e.secret)
	if err != nil {
		fmt.Println("error saving config:", err)
		os.Exit(1)
	}
}

type importCommand struct {
	url string
}

func newImportCommand(args []string) importCommand {
	ic := importCommand{}
	cmd := flag.NewFlagSet("import", flag.ExitOnError)
	cmd.StringVar(&ic.url, "url", "", "url")
	err := cmd.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ic
}

func (i importCommand) exec() {
	err := defaultStorage.Import(i.url)
	if err != nil {
		fmt.Println("error saving config:", err)
		os.Exit(1)
	}
}

type moveCommand struct {
	id     int
	offset int
}

func newMoveCommand(args []string) moveCommand {
	mc := moveCommand{}
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	cmd.IntVar(&mc.id, "id", 0, "id")
	cmd.IntVar(&mc.offset, "offset", 0, "offset")
	err := cmd.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return mc
}

func (m moveCommand) exec() {
	err := defaultStorage.Move(m.id-1, m.offset)
	if err != nil {
		fmt.Println("error saving config:", err)
		os.Exit(1)
	}
}

type listCommand struct{}

func newListCommand(_ []string) listCommand {
	lc := listCommand{}
	return lc
}

func (l listCommand) exec() {
	records, err := defaultStorage.readConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, record := range records {
		fmt.Printf("%2d  %s\n", i+1, record.Name)
	}
}
func doConfigCommand(args []string) {
	cmd := newConfigCommand(args)
	cmd.exec()
}

func newConfigCommand(args []string) *ConfigCommand {
	cmd := &ConfigCommand{args: args}
	//fmt.Println("args:", args)
	var subCommand string
	if len(args) >= 2 {
		subCommand = args[1]
	}
	switch subCommand {
	case "list":
		cmd.command = newListCommand(args[1:])
	case "add":
		cmd.command = newAddCommand(args[1:])
	case "delete":
		cmd.command = newDeleteCommand(args[1:])
	case "edit":
		cmd.command = newEditCommand(args[1:])
	case "import":
		cmd.command = newImportCommand(args[1:])
	case "move":
		cmd.command = newMoveCommand(args[1:])
	default:
		fmt.Println(strings.TrimSpace(usageConfig))
		os.Exit(0)
	}
	return cmd
}

func (c ConfigCommand) exec() {
	c.command.exec()
}
