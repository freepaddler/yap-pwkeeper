package client

import (
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/app/client/grpccli"
	"yap-pwkeeper/internal/pkg/models"
)

func (a *App) browser() {
	_ = a.store.Update()
	flex := tview.NewFlex()

	// list of documents
	items := tview.NewList().ShowSecondaryText(false)
	items.SetBorder(true).SetTitle("Items")

	// document edit form
	form := tview.NewForm()
	form.SetBorder(true)

	// list of categories
	categories := tview.NewList().ShowSecondaryText(false)
	categories.SetBorder(true).SetTitle("Categories")

	help := tview.NewForm()
	help.SetHorizontal(true)
	help.AddTextView("Quit: `Esc`    Add New Document: `A`    Get Server Updates `U`", "", 1, 1, true, false)

	flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(categories, 0, 1, false).
			AddItem(items, 0, 1, false).
			AddItem(form, 0, 3, true), 0, 1, true).
		AddItem(help, 3, 1, false), 0, 1, true)

	categories.AddItem("Cards", "", 0, func() {
		a.ui.SetFocus(items)
	})
	categories.AddItem("Credentials", "", 0, func() {
		a.ui.SetFocus(items)
	})
	categories.AddItem("Notes", "", 0, func() {
		a.ui.SetFocus(items)
	})

	categories.SetFocusFunc(func() {
		form.Clear(true)
		switch categories.GetCurrentItem() {
		case 0:
			a.cardsList(items, form)
		case 1:
			a.credentialsList(items, form)
		case 2:
			a.notesList(items, form)
		}
	})

	categories.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			a.cardsList(items, form)
		case 1:
			a.credentialsList(items, form)
		case 2:
			a.notesList(items, form)
		}
	})

	categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			a.ui.SetFocus(items)
			return nil
		case tcell.KeyLeft:
			return nil
		case tcell.KeyEsc:
			a.modalExit()
			return nil
		}
		switch event.Rune() {
		case 'a':
			a.modalErr("test")
			return nil
		case 'A':
			a.modalUnauthorized("relogin")
			return nil
		case 'u':
			a.requestUpdate(categories)
		case 'U':
			a.requestUpdate(categories)
		}
		return event
	})

	items.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			a.ui.SetFocus(form)
			return nil
		case tcell.KeyLeft:
			a.ui.SetFocus(categories)
			return nil
		case tcell.KeyEsc:
			a.ui.SetFocus(categories)
			return nil
		}
		return event
	})

	// logs page
	if a.debug {
		a.logsPage()
		flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlG {
				a.pages.SwitchToPage(logsPage)
				return nil
			}
			return event
		})
	}
	a.pages.AddPage(mainPage, flex, true, true)
	a.pages.SwitchToPage(mainPage)
	a.ui.SetFocus(categories)
}

func (a *App) notesList(list *tview.List, form *tview.Form) {
	list.Clear().SetTitle("Notes")
	for _, v := range a.store.GetNotesList() {
		v := *v
		list.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(form)
		})
		list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if list.HasFocus() {
				a.notesForm(secondaryText, form, list)
			} else {
				form.Clear(true)
			}
		})
	}
	list.SetFocusFunc(func() {
		_, id := list.GetItemText(list.GetCurrentItem())
		a.notesForm(id, form, list)
	})
}

func (a *App) notesForm(id string, form *tview.Form, list *tview.List) {
	doc := *a.store.GetNote(id)
	if doc.State == models.StateDeleted {
		form.Clear(true)
		form.SetCancelFunc(func() {
			a.ui.SetFocus(list)
		})
		return
	}

	form.Clear(true).SetTitle("Edit Note")
	form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddTextArea("Text", doc.Text, 30, 5, 4096, func(text string) {
		doc.Text = text
	})
	//form.AddTextView("Metadata", "", 30, 1, true, true)
	//form.AddDropDown("Delete Meta", []string{}, -1, nil)
	a.drawMetadata(form, &doc.Metadata)
	form.AddButton("Back", func() {
		a.ui.SetFocus(list)
	})
	form.AddButton("Add  Meta", func() {
		a.addMeta(form, &doc.Metadata)
	})
	form.AddButton("Save", func() {
		if err := a.store.UpdateNote(doc); err != nil {
			if errors.Is(grpccli.ErrAuthFail, err) {
				a.modalUnauthorized(err.Error())
			} else {
				a.modalErr(err.Error())
			}
		} else {
			if err := a.store.Update(); err != nil {
				if errors.Is(grpccli.ErrAuthFail, err) {
					a.modalUnauthorized(err.Error())
				} else {
					a.modalErr("Changes saved, but failed to get server updates" + err.Error())
				}
			} else {
				a.modalOk("Changes saved.")
			}

		}
	})
	form.AddButton("[red]Delete", func() {
		if err := a.store.DeleteNote(doc); err != nil {
			if errors.Is(grpccli.ErrAuthFail, err) {
				a.modalUnauthorized(err.Error())
			} else {
				a.modalErr(err.Error())
			}
		} else {
			a.modalOk("Changes saved.")
			go func() { _ = a.store.Update() }()
		}
	})
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetCancelFunc(func() {
		a.ui.SetFocus(list)
	})
}

func (a *App) cardsList(list *tview.List, form *tview.Form) {
	list.Clear().SetTitle("Cards")
	for _, v := range a.store.GetCardsList() {
		v := *v
		list.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(form)
		})
		list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if list.HasFocus() {
				a.cardsForm(secondaryText, form, list)
			} else {
				form.Clear(true)
			}
		})
	}
	list.SetFocusFunc(func() {
		_, id := list.GetItemText(list.GetCurrentItem())
		a.cardsForm(id, form, list)
	})
}

func (a *App) cardsForm(id string, form *tview.Form, list *tview.List) {
	doc := *a.store.GetCard(id)
	// TODO: deleted in store
	form.Clear(true).SetTitle("Edit Card")
	form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddInputField("Cardholder Name", doc.Cardholder, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddInputField("Card Number", doc.Number, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddInputField("Expires (mm/yy)", doc.Expires, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddInputField("CVC or CVV", doc.Code, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddInputField("PIN", doc.Pin, 30, nil, func(text string) {
		doc.Name = text
	})
	a.drawMetadata(form, &doc.Metadata)
	form.AddButton("Back", func() {
		a.ui.SetFocus(list)
	})
	form.AddButton("Add  Meta", func() {
		a.addMeta(form, &doc.Metadata)
	})
	form.AddButton("Save", nil)
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetCancelFunc(func() {
		a.ui.SetFocus(list)
	})
}

func (a *App) credentialsList(list *tview.List, form *tview.Form) {
	list.Clear().SetTitle("Credentials")
	for _, v := range a.store.GetCredentialsList() {
		v := *v
		list.AddItem(v.Name, v.Id, 0, func() {
			a.ui.SetFocus(form)
		})
		list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if list.HasFocus() {
				a.credentialsForm(secondaryText, form, list)
			} else {
				form.Clear(true)
			}
		})
	}
	list.SetFocusFunc(func() {
		_, id := list.GetItemText(list.GetCurrentItem())
		a.credentialsForm(id, form, list)
	})
}

func (a *App) credentialsForm(id string, form *tview.Form, list *tview.List) {
	doc := *a.store.GetCredential(id)
	// TODO: deleted in store
	form.Clear(true).SetTitle("Edit Credential")
	form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddInputField("Login", doc.Login, 30, nil, func(text string) {
		doc.Login = text
	})
	form.AddInputField("Password", doc.Password, 30, nil, func(text string) {
		doc.Password = text
	})
	a.drawMetadata(form, &doc.Metadata)
	form.AddButton("Back", func() {
		a.ui.SetFocus(list)
	})
	form.AddButton("Add  Meta", func() {
		a.addMeta(form, &doc.Metadata)
	})
	form.AddButton("Save", nil)
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetCancelFunc(func() {
		a.ui.SetFocus(list)
	})
}

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
		a.saveMeta(form, &newMeta, metadata)
	})
}

func (a *App) saveMeta(form *tview.Form, meta *models.Meta, metadata *[]models.Meta) {
	if meta.Key == "" {
		a.modalErr("Meta key should not be empty")
	} else {
		button := form.GetButton(1)
		*metadata = append(*metadata, *meta)
		button.SetLabel("Add  Meta")
		button.SetSelectedFunc(func() {
			a.addMeta(form, metadata)
		})
		a.drawMetadata(form, metadata)
	}
}

func (a *App) drawMetadata(form *tview.Form, metadata *[]models.Meta) {
	if idx := form.GetFormItemIndex("Metadata"); idx > -1 {
		form.RemoveFormItem(idx)
	}
	if idx := form.GetFormItemIndex("Delete Meta"); idx > -1 {
		form.RemoveFormItem(idx)
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
	switch {
	case size > 5:
		size = 5
	case size == 0:
		size = 1
	}

	form.AddTextView("Metadata", text, 30, size, true, true)
	form.AddDropDown("Delete Meta", options, -1, func(option string, optionIndex int) {
		//a.ui.SetFocus(form)
		if optionIndex >= 0 {
			*metadata = append((*metadata)[:optionIndex], (*metadata)[optionIndex+1:]...)
			a.drawMetadata(form, metadata)
			//a.ui.Draw()
		}
	})

	p, ok := a.ui.GetFocus().(*tview.DropDown)
	if ok && p.GetLabel() == "Delete Meta" {
		a.ui.SetFocus(form)
	}
}
