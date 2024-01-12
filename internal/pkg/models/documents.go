package models

// Meta is a random key-value pair, that may be added to any document
type Meta struct {
	Key   string `bson:"key"`   // meta key
	Value string `bson:"value"` // value
}

// Note is a document, that contains simple test
type Note struct {
	Id       string `bson:"_id,omitempty"` // document id
	UserId   string `bson:"user_id"`       // user id
	Serial   int64  `bson:"serial"`        // update serial number
	Name     string `bson:"name"`          // document name
	Text     string `bson:"text"`          // saved text
	Metadata []Meta `bson:"metadata"`      // document metadata
	State    string `bson:"state"`         // document state
}

// Credential is login-password pair
type Credential struct {
	Id       string `bson:"_id,omitempty"` // document id
	UserId   string `bson:"user_id"`       // user id
	Serial   int64  `bson:"serial"`        // update serial number
	Name     string `bson:"name"`          // document name
	Login    string `bson:"login"`         // saved login
	Password string `bson:"password"`      // saved password
	Metadata []Meta `bson:"metadata"`      // document metadata
	State    string `bson:"state"`         // document state
}

// Card carries credit cards data
type Card struct {
	Id         string `bson:"_id,omitempty"` // document id
	UserId     string `bson:"user_id"`       // user id
	Serial     int64  `bson:"serial"`        // update serial number
	Name       string `bson:"name"`          // document name
	Cardholder string `bson:"cardholder"`    // cardholder name
	Number     string `bson:"number"`        // card number
	Expires    string `bson:"expires"`       // card expiration date
	Pin        string `bson:"pin"`           // card pin
	Code       string `bson:"code"`          // card cvc/cvv2 code
	Metadata   []Meta `bson:"metadata"`      // document metadata
	State      string `bson:"state"`         // document state
}

// File document is a named file with metadata
type File struct {
	Id       string `bson:"_id,omitempty"` // documentId
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
