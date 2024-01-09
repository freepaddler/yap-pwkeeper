package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pageModal = "modal"
)

// modalErr shows error in modal window
func (a *App) modalErr(text string) {
	f := a.ui.GetFocus()
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText("Error: " + text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.HidePage(pageModal)
			a.ui.SetFocus(f)
		})
	a.pages.AddPage(pageModal, m, true, true)
}

// modalOk shows success modal window and switches to categories
func (a *App) modalOk(text string) {
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorDarkCyan).
		SetText("OK: " + text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage(pageModal)
			a.ui.SetFocus(a.categories)
		})
	a.pages.AddPage(pageModal, m, true, true)
}

// modalUnauthorized shows Unauthorized modal and returns to application start
func (a *App) modalUnauthorized() {
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText("Session authorisation failed!\nYou need to Login again.\nAll data is cleared.").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage(pageModal)
			a.pages.HidePage(pageRoot)
			a.welcomePage()
		})
	a.pages.AddPage(pageModal, m, true, true)
}

// modalExit is exit alarm window
func (a *App) modalExit() {
	f := a.ui.GetFocus()
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorDarkCyan).
		SetText("Are you sure you want to quit?").
		AddButtons([]string{"No", "Yes, quit application"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 1 {
				a.Stop()
			}
			a.pages.RemovePage(pageModal)
			a.ui.SetFocus(f)
		})
	a.pages.AddPage(pageModal, m, true, true)
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
