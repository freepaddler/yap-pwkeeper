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

// notesList display list of Notes
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

// cardsList display list of Cards
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

// credentialsList display list of Credentials
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

// filesList display list of Credentials
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

// addMeta adds controls to enter new Meta
func (a *App) addMeta(form *tview.Form, metadata *[]models.Meta) {
	newMeta := models.Meta{}
	button := form.GetButton(1)
	form.AddInputField("Add meta key", "", 30, nil, func(text string) {
		newMeta.Key = text
	})
	form.AddInputField("Add meta value", "", 30, nil, func(text string) {
		newMeta.Value = text
	})
	button.SetLabel("Save Meta")
	button.SetSelectedFunc(func() {
		a.saveMeta(&newMeta, metadata)
	})
}

// saveMeta applies new Meta to document
func (a *App) saveMeta(meta *models.Meta, metadata *[]models.Meta) {
	if meta.Key == "" {
		a.modalErr("Meta key should not be empty")
	} else {
		button := a.form.GetButton(1)
		*metadata = append(*metadata, *meta)
		button.SetLabel("Add  Meta")
		button.SetSelectedFunc(func() {
			a.addMeta(a.form, metadata)
		})
		a.drawMetadata(metadata)
	}
}

// drawMetadata displays metadata content and metadata delete dropdown
func (a *App) drawMetadata(metadata *[]models.Meta) {
	if idx := a.form.GetFormItemIndex("Add meta key"); idx > -1 {
		a.form.RemoveFormItem(idx)
	}
	if idx := a.form.GetFormItemIndex("Add meta value"); idx > -1 {
		a.form.RemoveFormItem(idx)
	}
	if idx := a.form.GetFormItemIndex("Metadata"); idx > -1 {
		a.form.RemoveFormItem(idx)
	}
	if idx := a.form.GetFormItemIndex("Delete Meta"); idx > -1 {
		a.form.RemoveFormItem(idx)
	}
	size := len(*metadata)
	text := ""
	options := make([]string, 0, size)
	for i, v := range *metadata {
		text += "[green]" + v.Key + "[white]: " + v.Value
		options = append(options, v.Key+": "+v.Value)
		if i < size-1 {
			text += "\n"
		}
	}
	// metadata view size
	if size > 5 {
		size = 5
	}
	if size == 0 {
		a.form.AddTextView("Metadata", text, 30, 1, true, true)
		//a.form.GetFormItemIndex("Metadata")
		//a.form.GetFormItem(a.form.GetFormItemIndex("Metadata")).SetDisabled(true)
	} else {
		a.form.AddTextView("Metadata", text, 30, size, true, true)
		a.form.AddDropDown("Delete Meta", options, -1, func(option string, optionIndex int) {
			//a.ui.SetFocus(form)
			if optionIndex >= 0 {
				*metadata = append((*metadata)[:optionIndex], (*metadata)[optionIndex+1:]...)
				a.drawMetadata(metadata)
			}
		})
	}

	p, ok := a.ui.GetFocus().(*tview.DropDown)
	if ok && p.GetLabel() == "Delete Meta" {
		a.ui.SetFocus(a.form)
	}
}
