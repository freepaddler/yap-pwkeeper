package client

import (
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/app/client/grpccli"
)

func (a *App) requestUpdate(focus *tview.List) {
	var err error
	a.pages.SwitchToPage(logsPage)
	err = a.store.Update()
	if err != nil {
		if errors.Is(grpccli.ErrAuthFail, err) {
			a.modalUnauthorized("Update failed: " + err.Error())
		} else {
			a.modalErr("Update failed: " + err.Error())
		}
	} else {
		a.modalOk("Updates from server received")
	}
}

func (a *App) modalErr(text string) {
	f := a.ui.GetFocus()
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText("Error: " + text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("error")
			a.ui.SetFocus(f)
		})
	a.pages.AddPage("error", m, true, true)
}

func (a *App) modalUnauthorized(text string) {
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText("Error: " + text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.store.Clear()
			a.welcomePage()
			a.pages.RemovePage("unauth")
			a.pages.RemovePage(mainPage)
		})
	a.pages.AddPage("unauth", m, true, true)
}

func (a *App) modalOk(text string) {
	f := a.ui.GetFocus()
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorDarkGreen).
		SetText("OK: " + text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("ok")
			a.ui.SetFocus(f)
		})
	a.pages.AddPage("ok", m, true, true)
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
			a.pages.RemovePage("exit")
		})
	a.pages.AddPage("exit", m, true, true)
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
