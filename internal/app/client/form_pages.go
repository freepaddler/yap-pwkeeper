package client

import (
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

const (
	formAdd = iota
	formModify
)

func (a *App) clearForm() {
	a.form.Clear(true).SetCancelFunc(func() {
		if a.itemsList.GetItemCount() > 0 {
			a.ui.SetFocus(a.itemsList)
			return
		}
		a.ui.SetFocus(a.categories)
	}).SetTitle("Documents View")
}

// notesForm draws forms for operations with Notes
func (a *App) notesForm(note *models.Note, formType int) {
	a.form.Clear(true)
	doc := models.Note{}
	switch formType {
	case formAdd:
		a.form.SetTitle(" Add Note ")
	case formModify:
		a.form.SetTitle(" Edit Note ")
		if note == nil {
			a.form.SetTitle(" [red]INVALID DOCUMENT ")
			return
		}
		doc = *note
	default:
		a.form.SetTitle(" [red]INVALID FORM ")
	}
	a.form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	a.form.AddTextArea("Text", doc.Text, 30, 5, 4096, func(text string) {
		doc.Text = text
	})
	a.drawMetadata(&doc.Metadata)
	a.form.AddButton("<< Back", func() {
		a.ui.SetFocus(a.itemsList)
	})
	a.form.AddButton("Add  Meta", func() {
		a.addMeta(a.form, &doc.Metadata)
	})

	// buttons
	switch formType {
	case formAdd:
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.AddNote(doc)
				},
				"New Note saved",
				"Failed to save Note",
			)
		})
	case formModify:
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.UpdateNote(doc)
				},
				"Note saved",
				"Failed to save Note",
			)
		})
		a.form.AddButton("[red]Delete", func() {
			a.modifyRequest(
				func() error {
					return a.store.DeleteNote(doc)
				},
				"Note deleted",
				"Failed to delete Note",
			)
		})
	}

	a.form.SetButtonsAlign(tview.AlignCenter)
	a.form.SetCancelFunc(func() {
		a.ui.SetFocus(a.itemsList)
	})
}

// cardsForm draws forms for operations with Cards
func (a *App) cardsForm(card *models.Card, formType int) {
	a.form.Clear(true)
	doc := models.Card{}
	switch formType {
	case formAdd:
		a.form.SetTitle(" Add Card ")
	case formModify:
		a.form.SetTitle(" Edit Card ")
		if card == nil {
			a.form.SetTitle(" [red]INVALID DOCUMENT ")
			return
		}
		doc = *card
	default:
		a.form.SetTitle(" [red]INVALID FORM ")
	}
	a.form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	a.form.AddInputField("Cardholder Name", doc.Cardholder, 30, nil, func(text string) {
		doc.Cardholder = text
	})
	a.form.AddInputField("Card Number", doc.Number, 30, nil, func(text string) {
		doc.Number = text
	})
	a.form.AddInputField("Expires (mm/yy)", doc.Expires, 30, nil, func(text string) {
		doc.Expires = text
	})
	a.form.AddInputField("CVC or CVV", doc.Code, 30, nil, func(text string) {
		doc.Code = text
	})
	a.form.AddInputField("PIN", doc.Pin, 30, nil, func(text string) {
		doc.Pin = text
	})
	a.drawMetadata(&doc.Metadata)
	a.form.AddButton("<< Back", func() {
		a.ui.SetFocus(a.itemsList)
	})
	a.form.AddButton("Add  Meta", func() {
		a.addMeta(a.form, &doc.Metadata)
	})

	// buttons
	switch formType {
	case formAdd:
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.AddCard(doc)
				},
				"New Card saved",
				"Failed to save Card",
			)
		})
	case formModify:
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.UpdateCard(doc)
				},
				"Card saved",
				"Failed to save Card",
			)
		})
		a.form.AddButton("[red]Delete", func() {
			a.modifyRequest(
				func() error {
					return a.store.DeleteCard(doc)
				},
				"Card deleted",
				"Failed to delete Card",
			)
		})
	}

	a.form.SetButtonsAlign(tview.AlignCenter)
	a.form.SetCancelFunc(func() {
		a.ui.SetFocus(a.itemsList)
	})
}

// credentialsForm draws forms for operations with Credentials
func (a *App) credentialsForm(cred *models.Credential, formType int) {
	a.form.Clear(true)
	doc := models.Credential{}
	switch formType {
	case formAdd:
		a.form.SetTitle(" Add Credential ")
	case formModify:
		a.form.SetTitle(" Edit Credential ")
		if cred == nil {
			a.form.SetTitle(" [red]INVALID DOCUMENT ")
			return
		}
		doc = *cred
	default:
		a.form.SetTitle(" [red]INVALID FORM ")
	}
	a.form.AddInputField("Name", doc.Name, 30, nil, func(text string) {
		doc.Name = text
	})
	a.form.AddInputField("Login", doc.Login, 30, nil, func(text string) {
		doc.Login = text
	})
	a.form.AddInputField("Password", doc.Password, 30, nil, func(text string) {
		doc.Password = text
	})
	a.drawMetadata(&doc.Metadata)
	a.form.AddButton("Back", func() {
		a.ui.SetFocus(a.itemsList)
	})
	a.form.AddButton("Add  Meta", func() {
		a.addMeta(a.form, &doc.Metadata)
	})

	// buttons
	switch formType {
	case formAdd:
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.AddCredential(doc)
				},
				"New Credential saved",
				"Failed to save Credential",
			)
		})
	case formModify:
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.UpdateCredential(doc)
				},
				"Credential saved",
				"Failed to save Credential",
			)
		})
		a.form.AddButton("[red]Delete", func() {
			a.modifyRequest(
				func() error {
					return a.store.DeleteCredential(doc)
				},
				"Credential deleted",
				"Failed to delete Credential",
			)
		})
	}

	a.form.SetButtonsAlign(tview.AlignCenter)
	a.form.SetCancelFunc(func() {
		a.ui.SetFocus(a.itemsList)
	})
}
