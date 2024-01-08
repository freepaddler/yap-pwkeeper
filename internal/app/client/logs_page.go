package client

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) logsPage() {
	//f := a.ui.GetFocus()
	logview := tview.NewTextView().SetScrollable(true).ScrollToEnd()
	log.SetOutput(logview)
	log.Println("test message")
	logview.SetBorder(true).SetTitle("Logs")
	logview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlG {
			a.pages.SwitchToPage(mainPage)
		}
		return event

	})
	logview.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			a.pages.SwitchToPage(mainPage)
		}
	})
	a.pages.AddPage(logsPage, logview, true, true)
}
