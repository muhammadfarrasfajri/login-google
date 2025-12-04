package models

type User struct {
	ID              int
	GoogleUID       string
	Name            string
	Email           string
	Google_picture  string
	Role            string
	Profile_picture string
}
