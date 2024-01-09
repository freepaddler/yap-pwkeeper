package models

type Meta struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type Note struct {
	Id       string `bson:"_id,omitempty"`
	UserId   string `bson:"user_id"`
	Serial   int64  `bson:"serial"`
	Name     string `bson:"name"`
	Text     string `bson:"text"`
	Metadata []Meta `bson:"metadata"`
	State    string `bson:"state"`
}

type Credential struct {
	Id       string `bson:"_id,omitempty"`
	UserId   string `bson:"user_id"`
	Serial   int64  `bson:"serial"`
	Name     string `bson:"name"`
	Login    string `bson:"login"`
	Password string `bson:"password"`
	Metadata []Meta `bson:"metadata"`
	State    string `bson:"state"`
}

type Card struct {
	Id         string `bson:"_id,omitempty"`
	UserId     string `bson:"user_id"`
	Serial     int64  `bson:"serial"`
	Name       string `bson:"name"`
	Cardholder string `bson:"cardholder"`
	Number     string `bson:"number"`
	Expires    string `bson:"expires"`
	Pin        string `bson:"pin"`
	Code       string `bson:"code"`
	Metadata   []Meta `bson:"metadata"`
	State      string `bson:"state"`
}
