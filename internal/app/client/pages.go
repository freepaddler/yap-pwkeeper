package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) modalErr(text string) {
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText("Error: " + text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("error")
		})
	a.pages.AddPage("error", m, true, true)
}

func (a *App) modalExit() {
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorDarkCyan).
		SetText("Are you sure you want to quit?").
		AddButtons([]string{"No", "Yes, quit application"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 1 {
				a.Stop()
			}
			a.pages.RemovePage("error")
		})
	a.pages.AddPage("error", m, true, true)
}

func center(width, height int, p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
