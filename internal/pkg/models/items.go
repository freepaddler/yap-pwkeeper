package models

type Meta struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type Document struct {
	UserId string `bson:"user_id"`
	Entity `bson:"inline"`
}

type Note struct {
	Id       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	Text     string `bson:"text"`
	Metadata []Meta `bson:"metadata"`
}

type NoteDocument struct {
	Note     `bson:"inline"`
	Document `bson:"inline"`
}
