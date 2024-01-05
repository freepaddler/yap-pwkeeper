package models

const (
	ItemTypeNote  ItemType = "Note"
	ItemTypeLogin ItemType = "Login"
	ItemTypeFile  ItemType = "File"
	ItemTypeCard  ItemType = "Card"
)

type ItemType string

type AuthToken string

type MetaItem struct {
	key   string
	value string
}

type MetaQ []MetaItem

type NoteItem struct {
	Id   string
	Name string
	Text string
	MetaQ
}

type LoginItem struct {
	Id       string
	Name     string
	Login    string
	Password string
	MetaQ
}

func (li *LoginItem) Item() ItemQ {
	return ItemQ{
		Id:   li.Id,
		Name: li.Name,
		Type: ItemTypeLogin,
	}
}

type ItemQ struct {
	Id   string
	Name string
	Type ItemType
}

type Wallet struct {
	Login []LoginItem
	Note  []NoteItem
}

func (w *Wallet) GetCategories() []ItemType {
	return []ItemType{ItemTypeNote, ItemTypeLogin}
}

func (w *Wallet) GetCategoryItems(itemType ItemType) []ItemQ {
	switch itemType {
	default:
		return nil
	}
}

func (w *Wallet) GetItem(itemType ItemType, id string) {

}
