package client

import (
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

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
		a.form.GetFormItemIndex("Metadata")
		a.form.GetFormItem(a.form.GetFormItemIndex("Metadata")).SetDisabled(true)
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
