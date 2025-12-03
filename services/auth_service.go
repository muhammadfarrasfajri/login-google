package services

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/repository"
)

var (
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUserNotRegistered = errors.New("user not registered, please register first")
)

type AuthService struct {
	Repo     repository.AuthRepository
	FirebaseAuth *firebase.Client
	JWTSecret    string
}

// --------------------------- REGISTER -----------------------------------

func (s *AuthService) Register(idToken string, customName string) (*models.BaseUser, error) {
	ctx := context.Background()

	// 1. Verifikasi Firebase ID Token
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	googleUID := token.UID

	email, _ := token.Claims["email"].(string)
	googlePicture, _ := token.Claims["picture"].(string)

	// 2. Cek apakah user sudah ada
	existing, _ := s.Repo.FindByGoogleUID(googleUID)
	if existing != nil {
		return nil, errors.New("user already registered, please login")
	}

	// 3. Tentukan nama
	name := customName
	if name == "" {
		if n, ok := token.Claims["name"].(string); ok {
			name = n
		}
	}
	
	// 5. Simpan user
	newUser := models.BaseUser{
		GoogleUID: googleUID,
		Name:      name,
		Email:     email,
		GooglePicture:   googlePicture,
	}

	err = s.Repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// -------------------------- LOGIN ----------------------------------------

func (s *AuthService) Login(idToken string, deviceInfo string, ip string) (map[string]interface{}, error) {
	ctx := context.Background()

	// 1. Verifikasi Firebase Token
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	googleUID := token.UID

	// 2. Cek user di DB
	user, err := s.Repo.FindByGoogleUID(googleUID)
	if err != nil || user == nil {
		return nil, ErrUserNotRegistered
	}

	// 3. Simpan aktivitas login
	err = s.Repo.SaveLoginHistory(user.ID, deviceInfo, ip)
	if err != nil {
		return nil, err
	}

	jwtToken, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "login success",
		"token":   jwtToken,
		"user":    user,
	}, nil
}