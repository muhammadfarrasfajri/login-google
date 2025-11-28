package models

type LoginHistory struct {
	ID        int
	UserID    string
	LoginTime string
	Device    string
	IP        string
}
