package client

import (
	"errors"
	"fmt"

	"yap-pwkeeper/internal/app/client/memstore"
)

// setStatus updates status message
func (a *App) setStatus(text string, color string) {
	if color == "" {
		color = "white"
	}
	if a.statusBar.GetFormItemCount() > 1 {
		a.statusBar.RemoveFormItem(1)
	}
	a.statusBar.AddTextView(fmt.Sprintf("[%s]%s", color, text), "", 1, 1, false, false)
}

func (a *App) statusOK(text string) {
	a.setStatus(text, "green")
}

func (a *App) statusFail(text string) {
	a.setStatus(text, "red")
}

// synchronize run store update
func (a *App) synchronize() {
	if err := a.store.Update(); err != nil {
		if errors.Is(memstore.ErrAuthFailed, err) {
			a.modalUnauthorized()
		} else {
			a.statusFail("Synchronization failed: " + err.Error())
		}
	} else {
		a.statusOK("Synchronized")
		a.categories.SetItemText(0, fmt.Sprintf("Cards (%d)", len(a.store.GetCardsList())), "[yellow](`C` to add new)")
		a.categories.SetItemText(1, fmt.Sprintf("Logins (%d)", len(a.store.GetCredentialsList())), "[yellow](`L` to add new)")
		a.categories.SetItemText(2, fmt.Sprintf("Notes (%d)", len(a.store.GetNotesList())), "[yellow](`N` to add new)")
	}
}

// modifyRequest wraps any data modification call and
func (a *App) modifyRequest(fn func() error, okMsg, failMsg string) {
	if err := fn(); err != nil {
		if errors.Is(memstore.ErrAuthFailed, err) {
			a.modalUnauthorized()
		} else {
			a.modalErr(failMsg + ": " + err.Error())
		}
	} else {
		a.modalOk(okMsg)
	}
}
