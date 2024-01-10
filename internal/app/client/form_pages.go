package client

import (
	"fmt"

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
	a.form.AddInputField("Name", doc.Name, 50, nil, func(text string) {
		doc.Name = text
	})
	a.form.AddTextArea("Text", doc.Text, 50, 5, 4096, func(text string) {
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
	a.form.AddInputField("Name", doc.Name, 50, nil, func(text string) {
		doc.Name = text
	})
	a.form.AddInputField("Cardholder Name", doc.Cardholder, 50, nil, func(text string) {
		doc.Cardholder = text
	})
	a.form.AddInputField("Card Number", doc.Number, 50, nil, func(text string) {
		doc.Number = text
	})
	a.form.AddInputField("Expires (mm/yy)", doc.Expires, 50, nil, func(text string) {
		doc.Expires = text
	})
	a.form.AddInputField("CVC or CVV", doc.Code, 50, nil, func(text string) {
		doc.Code = text
	})
	a.form.AddInputField("PIN", doc.Pin, 50, nil, func(text string) {
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
	a.form.AddInputField("Name", doc.Name, 50, nil, func(text string) {
		doc.Name = text
	})
	a.form.AddInputField("Login", doc.Login, 50, nil, func(text string) {
		doc.Login = text
	})
	a.form.AddInputField("Password", doc.Password, 50, nil, func(text string) {
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

// filesForm draws forms for operations with Credentials
func (a *App) filesForm(cred *models.File, formType int) {
	a.form.Clear(true)
	doc := models.File{}
	var newFile string
	switch formType {
	case formAdd:
		a.form.SetTitle(" Add File ")
	case formModify:
		a.form.SetTitle(" Edit File ")
		if cred == nil {
			a.form.SetTitle(" [red]INVALID DOCUMENT ")
			return
		}
		doc = *cred
	default:
		a.form.SetTitle(" [red]INVALID FORM ")
	}
	a.form.AddInputField("Name", doc.Name, 50, nil, func(text string) {
		doc.Name = text
	})
	switch formType {
	case formAdd:
		a.form.AddInputField("Path to file", "", 50, nil, func(text string) {
			newFile = text
		})
	case formModify:
		a.form.AddInputField("Path to NEW file", "", 50, nil, func(text string) {
			newFile = text
		})
		a.form.AddTextView("File Info", fmt.Sprintf("[green]File Name[white]: %s\n[green]File Size: [white]%d bytes", doc.Filename, doc.Size), 50, 2, true, false)
		a.form.GetFormItemByLabel("File Info").SetDisabled(true)
	}
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
			if newFile == "" {
				a.modalErr("Path to file should not be empty")
				return
			}
			if err := checkFile(newFile); err != nil {
				a.modalErr(err.Error())
				return
			}
			a.modifyRequest(
				func() error {
					return a.store.AddFile(doc, newFile)
				},
				"New File saved",
				"Failed to save File",
			)
		})
	case formModify:
		a.form.AddButton("Download", func() {
			a.fileInput(doc.Id, doc.Filename)
		})
		a.form.AddButton("Save", func() {
			if doc.Name == "" {
				a.modalErr("Document name should not be empty")
				return
			}
			if newFile != "" {
				if err := checkFile(newFile); err != nil {
					a.modalErr(err.Error())
					return
				}
				a.modifyRequest(
					func() error {
						return a.store.UpdateFile(doc, newFile)
					},
					"File saved",
					"Failed to save File",
				)
			} else {
				a.modifyRequest(
					func() error {
						return a.store.UpdateFileInfo(doc)
					},
					"File saved",
					"Failed to save File",
				)
			}
		})
		a.form.AddButton("[red]Delete", func() {
			a.modifyRequest(
				func() error {
					return a.store.DeleteFile(doc)
				},
				"File deleted",
				"Failed to delete File",
			)
		})
	}

	a.form.SetButtonsAlign(tview.AlignCenter)
	a.form.SetCancelFunc(func() {
		a.ui.SetFocus(a.itemsList)
	})
}
