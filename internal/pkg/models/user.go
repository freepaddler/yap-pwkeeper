package models

type User struct {
	Id           string `bson:"_id,omitempty"`
	Login        string `bson:"login"`
	PasswordHash string `bson:"password"`
	Entity       `bson:"inline"`
}

// UserCredentials are aaa login and password
type UserCredentials struct {
	Login    string
	Password string
}
