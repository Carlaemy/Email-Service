package model

type Recipient struct {
	Name     string
	FileName string
	FileID   string
	Action   string
}

type File struct {
	FileID int
	Name   string
	Size   int
}

type User struct {
	UserID   int
	Name     string
	LastName string
	Email    string
	Password string
}
