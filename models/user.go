package models

type BaseUser struct {
	ID             int
	GoogleUID      string
	Name           string
	Email          string
	GooglePicture  string
	ProfilePicture string
	Role           string
	IsLoggedIn     int
}
