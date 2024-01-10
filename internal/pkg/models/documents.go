package models

// Meta is a random key-value pair, that may be added to any document
type Meta struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

// Note is a document, that contains simple test
type Note struct {
	Id       string `bson:"_id,omitempty"`
	UserId   string `bson:"user_id"`
	Serial   int64  `bson:"serial"`
	Name     string `bson:"name"`
	Text     string `bson:"text"`
	Metadata []Meta `bson:"metadata"`
	State    string `bson:"state"`
}

// Credential is login-password pair
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

// Card carries credit cards data
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

// File document consists of a fiie
type File struct {
	Id       string `bson:"_id,omitempty"`
	UserId   string `bson:"user_id"`
	Serial   int64  `bson:"serial"`
	Name     string `bson:"name"`
	Filename string `bson:"filename"`
	Size     int64  `bson:"size"`
	Sha265   string `bson:"sha265"`
	Data     []byte `bson:"data"`
	Metadata []Meta `bson:"metadata"`
	State    string `bson:"state"`
}
