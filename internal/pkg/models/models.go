package models

const (
	ItemTypeNote  ItemType = "Note"
	ItemTypeLogin ItemType = "Login"
	ItemTypeFile  ItemType = "File"
	ItemTypeCard  ItemType = "Card"
)

type ItemType string

type LoginRequest struct {
	Login    string
	Password string
}

type AuthToken string

type MetaItem struct {
	key   string
	value string
}

type Meta []MetaItem

type NoteItem struct {
	Id   string
	Name string
	Text string
	Meta
}

type LoginItem struct {
	Id       string
	Name     string
	Login    string
	Password string
	Meta
}

func (li *LoginItem) Item() Item {
	return Item{
		Id:   li.Id,
		Name: li.Name,
		Type: ItemTypeLogin,
	}
}

type Item struct {
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

func (w *Wallet) GetCategoryItems(itemType ItemType) []Item {
	switch itemType {
	default:
		return nil
	}
}

func (w *Wallet) GetItem(itemType ItemType, id string) {

}
