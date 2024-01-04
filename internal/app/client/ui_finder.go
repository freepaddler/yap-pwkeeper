package client

import "github.com/rivo/tview"

func (a *App) finder() {
	categories := tview.NewList().ShowSecondaryText(false)
	categories.SetBorder(true).SetTitle("Categories")

	items := tview.NewList().ShowSecondaryText(false)
	items.SetBorder(true).SetTitle("Items")

	//flex := tview.NewFlex().
	//	AddItem(categories, 0, 1, true).
	//	AddItem(items, 0, 1, false)

}
