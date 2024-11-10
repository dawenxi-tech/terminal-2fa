package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dialog struct {
	modal *tview.Flex

	width, height int
	view          tview.Primitive

	close     func()
	closeRune rune
	escClose  bool
}

func NewDialog(view tview.Primitive, width, height int) *Dialog {
	dig := &Dialog{
		width:     width,
		height:    height,
		view:      view,
		escClose:  true,
		closeRune: 'q',
	}
	dig.layout()
	dig.beginListen()
	return dig
}

func (dig *Dialog) setClose(fm func()) {
	dig.close = fm
}

func (dig *Dialog) beginListen() {
	dig.modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			if dig.close != nil {
				dig.close()
			}
		}
		if event.Key() == tcell.KeyRune && event.Rune() == dig.closeRune {
			if dig.close != nil {
				dig.close()
			}
		}
		return event
	})
}

func (dig *Dialog) layout() {
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(dig.view, dig.height, 1, true).
			AddItem(nil, 0, 1, false), dig.width, 1, true).
		AddItem(nil, 0, 1, false)
	dig.modal = flex
}
