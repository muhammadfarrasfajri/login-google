package services

import (
	"context"
	"errors"
	"time"

	firebase "firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/repository"
)

var (
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUserNotRegistered = errors.New("user not registered, please register first")
)

type AuthService struct {
	UserRepo     *repository.UserRepository
	FirebaseAuth *firebase.Client
	JWTSecret    string
}

// --------------------------- REGISTER -----------------------------------

func (s *AuthService) Register(idToken string, customName string) (*models.User, error) {
	ctx := context.Background()

	// 1. Verifikasi Firebase ID Token
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	googleUID := token.UID
	email := token.Claims["email"].(string)
	picture := token.Claims["picture"].(string)

	// 2. Cek apakah sudah ada
	existing, _ := s.UserRepo.FindByGoogleUID(googleUID)
	if existing != nil {
		return nil, errors.New("user already registered, please login")
	}

	name := customName
	if name == "" {
		if n, ok := token.Claims["name"].(string); ok {
			name = n
		}
	}

	// 3. Simpan user
	newUser := models.User{
		GoogleUID: googleUID,
		Name:      name,
		Email:     email,
		Picture:   picture,
	}

	err = s.UserRepo.Create(newUser)
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
	user, err := s.UserRepo.FindByGoogleUID(googleUID)
	if err != nil || user == nil {
		return nil, ErrUserNotRegistered
	}

	// 3. Simpan aktivitas login
	err = s.UserRepo.SaveLoginHistory(user.ID, deviceInfo, ip)
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

// -------------------------- JWT -----------------------------------------

func (s *AuthService) GenerateJWT(userID string, email string, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.JWTSecret))
}
