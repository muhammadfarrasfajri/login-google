package repository

import (
	"time"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type AuthRepository interface {
	// Register and Login
	Create(user models.BaseUser) error
	SaveLoginHistory(userID int, deviceInfo, ip string) error
	UpdateLoginStatus(id int, status int) error

	// CRUD
	FindByGoogleUID(uid string) (*models.BaseUser, error)
	FindByID(id int) (*models.BaseUser, error)
	GetAll() ([]models.BaseUser, error)
	Update(user models.BaseUser) error
	Delete(id string) error
	UpdatePhotoURL(userID int, url string) error
	

	// Refresh Token
	RefreshToken(userID int, refreshToken string, exp time.Time) error
	FindRefreshToken(userID int) (*models.RefreshToken, error)
	UpdateRefreshToken(userID int, newRefreshToken string, exp time.Time) error
	DeleteRefreshToken(UserID int) error
}
