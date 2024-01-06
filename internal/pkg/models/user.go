package models

// UserCredentials are user login and password to register and login
type UserCredentials struct {
	Login    string
	Password string
}

// User is user storage implementation
type User struct {
	Id           string `bson:"_id,omitempty"`
	Login        string `bson:"login"`
	PasswordHash string `bson:"password"`
	State        string `bson:"state"`
}
