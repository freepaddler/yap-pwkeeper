package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

// mainPage is application root page listing and modifying documents
func (a *App) mainPage() {
	flex := tview.NewFlex()

	// list of categories
	a.categories.ShowSecondaryText(true)
	a.categories.SetBorder(true).SetTitle(" Categories ")

	flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(a.categories, 0, 1, false).
			AddItem(a.itemsList, 0, 1, false).
			AddItem(a.form, 0, 3, true), 0, 1, true).
		AddItem(a.statusBar, 3, 1, false), 0, 1, true)

	a.categories.AddItem("", "", 0, nil)
	a.categories.AddItem("", "", 0, nil)
	a.categories.AddItem("", "", 0, nil)
	a.categories.AddItem("", "", 0, nil)

	a.categories.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		if a.itemsList.GetItemCount() > 0 {
			a.ui.SetFocus(a.itemsList)
		}
	})
	a.categories.SetFocusFunc(func() {
		a.clearForm()
		a.synchronize()
		switch a.categories.GetCurrentItem() {
		case 0:
			a.cardsList()
		case 1:
			a.credentialsList()
		case 2:
			a.notesList()
		case 3:
			a.filesList()
		}
	})

	a.categories.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			a.cardsList()
		case 1:
			a.credentialsList()
		case 2:
			a.notesList()
		case 3:
			a.filesList()
		}
	})

	a.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			if a.itemsList.GetItemCount() > 0 {
				a.ui.SetFocus(a.itemsList)
			}
			return nil
		case tcell.KeyLeft:
			return nil
		case tcell.KeyEsc:
			a.modalExit()
			return nil
		}
		switch event.Rune() {
		case 's':
			a.synchronize()
			return nil
		case 'c':
			a.cardsForm(&models.Card{}, formAdd)
			a.ui.SetFocus(a.form)
		case 'l':
			a.credentialsForm(&models.Credential{}, formAdd)
			a.ui.SetFocus(a.form)
		case 'n':
			a.notesForm(&models.Note{}, formAdd)
			a.ui.SetFocus(a.form)
		case 'f':
			a.filesForm(&models.File{}, formAdd)
			a.ui.SetFocus(a.form)
		}

		return event
	})

	// list of documents
	a.itemsList.ShowSecondaryText(false)
	a.itemsList.SetBorder(true).SetTitle(" Items ")

	a.itemsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			a.ui.SetFocus(a.form)
			return nil
		case tcell.KeyLeft:
			a.ui.SetFocus(a.categories)
			return nil
		case tcell.KeyEsc:
			a.ui.SetFocus(a.categories)
			return nil
		}
		switch event.Rune() {
		case 's':
			a.ui.SetFocus(a.categories)
			return nil
		case 'c':
			a.cardsForm(&models.Card{}, formAdd)
			a.ui.SetFocus(a.form)
		case 'l':
			a.credentialsForm(&models.Credential{}, formAdd)
			a.ui.SetFocus(a.form)
		case 'n':
			a.notesForm(&models.Note{}, formAdd)
			a.ui.SetFocus(a.form)
		case 'f':
			a.filesForm(&models.File{}, formAdd)
			a.ui.SetFocus(a.form)
		}
		return event
	})

	// document edit form
	a.form.SetBorder(true)
	a.clearForm()

	a.pages.AddPage(pageRoot, flex, true, false)
}

// notesList displays list of Notes
func (a *App) notesList() {
	a.itemsList.Clear().SetTitle("Notes")
	for _, v := range a.store.GetNotesList() {
		v := *v
		a.itemsList.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(a.form)
		})
		a.itemsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if a.itemsList.HasFocus() {
				a.notesForm(a.store.GetNote(secondaryText), formModify)
			} else {
				a.clearForm()
			}
		})
	}
	a.itemsList.SetFocusFunc(func() {
		_, id := a.itemsList.GetItemText(a.itemsList.GetCurrentItem())

		a.notesForm(a.store.GetNote(id), formModify)
	})
}

// cardsList displays list of Cards
func (a *App) cardsList() {
	a.itemsList.Clear().SetTitle("Cards")
	for _, v := range a.store.GetCardsList() {
		v := *v
		a.itemsList.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(a.form)
		})
		a.itemsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if a.itemsList.HasFocus() {
				a.cardsForm(a.store.GetCard(secondaryText), formModify)
			} else {
				a.clearForm()
			}
		})
	}
	a.itemsList.SetFocusFunc(func() {
		_, id := a.itemsList.GetItemText(a.itemsList.GetCurrentItem())
		a.cardsForm(a.store.GetCard(id), formModify)
	})
}

// credentialsList displays list of Credentials
func (a *App) credentialsList() {
	a.itemsList.Clear().SetTitle("Credentials")
	for _, v := range a.store.GetCredentialsList() {
		v := *v
		a.itemsList.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(a.form)
		})
		a.itemsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if a.itemsList.HasFocus() {
				a.credentialsForm(a.store.GetCredential(secondaryText), formModify)
			} else {
				a.clearForm()
			}
		})
	}
	a.itemsList.SetFocusFunc(func() {
		_, id := a.itemsList.GetItemText(a.itemsList.GetCurrentItem())
		a.credentialsForm(a.store.GetCredential(id), formModify)
	})
}

// filesList displays list of Files
func (a *App) filesList() {
	a.itemsList.Clear().SetTitle("Files")
	for _, v := range a.store.GetFilesList() {
		v := *v
		a.itemsList.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(a.form)
		})
		a.itemsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if a.itemsList.HasFocus() {
				a.filesForm(a.store.GetFileInfo(secondaryText), formModify)
			} else {
				a.clearForm()
			}
		})
	}
	a.itemsList.SetFocusFunc(func() {
		_, id := a.itemsList.GetItemText(a.itemsList.GetCurrentItem())
		a.filesForm(a.store.GetFileInfo(id), formModify)
	})
}
