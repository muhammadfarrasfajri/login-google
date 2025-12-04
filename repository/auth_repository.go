package repository

import "github.com/muhammadfarrasfajri/login-google/models"

type AuthRepository interface {
	FindByGoogleUID(uid string) (*models.BaseUser, error)
	Create(user models.BaseUser) error
	SaveLoginHistory(userID int, deviceInfo, ip string) error
	UpdateLoginStatus(id int, status int) error
}
