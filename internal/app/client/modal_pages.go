package client

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pageModal    = "modal"
	pageDownload = "pageDownload"
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

// fileInput draws input field to download file
func (a *App) fileInput(documentId, fname string) {
	var filename string
	cancelFunc := func() {
		a.pages.RemovePage(pageDownload)
		a.ui.SetFocus(a.form)
	}
	form := tview.NewForm().SetCancelFunc(cancelFunc)
	form.AddInputField("Download to:", "", 45, nil, func(text string) {
		filename = text
	})
	form.AddButton("Download", func() {
		if filename == "" {
			a.modalErr("Please enter path to save file")
			return
		}
		st, _ := os.Stat(filename)
		if st != nil && st.IsDir() {
			filename = filename + string(os.PathSeparator) + fname
		}
		if _, err := os.Stat(filename); err == nil {
			a.modalErr("File already exists, please enter another path.")
			return
		}
		a.pages.RemovePage(pageDownload)
		a.ui.SetFocus(a.form)
		a.modifyRequest(
			func() error {
				return a.store.GetFile(documentId, filename)
			},
			"File saved to: "+filename,
			"Failed to save File",
		)

	})
	form.AddButton("<< Cancel", func() {
		cancelFunc()
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.SetBorder(true)
	form.SetTitle(" Enter path where to save " + fname + " ")
	a.pages.AddPage(pageDownload, center(62, 7, form), true, true)
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
