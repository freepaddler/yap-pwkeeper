package client

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

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

	// layout
	flex.AddItem(categories, 0, 1, true)
	flex.AddItem(items, 0, 1, true)
	flex.AddItem(form, 0, 3, true)

	categories.AddItem(fmt.Sprintf("Cards (%d)", len(a.store.GetCardsList())), "2", 0, func() {
		a.ui.SetFocus(items)
	})
	categories.AddItem(fmt.Sprintf("Credentials (%d)", len(a.store.GetCredentialsList())), "", 0, func() {
		a.ui.SetFocus(items)
	})
	categories.AddItem(fmt.Sprintf("Notes (%d)", len(a.store.GetNotesList())), "", 0, func() {
		a.ui.SetFocus(items)
	})

	categories.SetFocusFunc(func() {
		//if categories.GetCurrentItem() == 0 {
		//	a.cardsList(items, form)
		//}
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
	// TODO: deleted in store
	form.Clear(true).SetTitle("Edit Note")
	form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	form.AddTextArea("Text", doc.Text, 30, 5, 4096, func(text string) {
		doc.Text = text
	})
	drawMetadata(form, &doc.Metadata)
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
	drawMetadata(form, &doc.Metadata)
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
	drawMetadata(form, &doc.Metadata)
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
		drawMetadata(form, metadata)
	}
}

func drawMetadata(form *tview.Form, metadata *[]models.Meta) {
	if idx := form.GetFormItemIndex("Metadata"); idx > -1 {
		form.RemoveFormItem(idx)
	}
	if idx := form.GetFormItemIndex("Delete Meta"); idx > -1 {
		form.RemoveFormItem(idx)
	}
	size := len(*metadata)
	if size == 0 {
		return
	}
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
	form.AddTextView("Metadata", text, 30, size, true, true)
	form.AddDropDown("[red]Delete Meta", options, -1, func(option string, optionIndex int) {
		if optionIndex >= 0 {
			*metadata = append((*metadata)[:optionIndex], (*metadata)[optionIndex+1:]...)
			drawMetadata(form, metadata)
		}
	})
}
