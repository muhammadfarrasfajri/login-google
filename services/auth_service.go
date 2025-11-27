package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	firebase "firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/muhammadfarrasfajri/login-google/repository"
)

var (
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserNotRegistered  = errors.New("user not registered, please complete registration")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrSessionCreation    = errors.New("failed to create session")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	UserRepo     *repository.UserRepository
	FirebaseAuth *firebase.Client
}

// LoginResponse contains the login result data
type LoginResponse struct {
	User         *UserData `json:"user"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserData represents user information
type UserData struct {
	UID     string `json:"uid"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo *repository.UserRepository, firebaseAuth *firebase.Client) *AuthService {
	return &AuthService{
		UserRepo:     userRepo,
		FirebaseAuth: firebaseAuth,
	}
}

// Register handles user registration with Firebase authentication
func (s *AuthService) Register(idToken string, deviceInfo string, ipAddress string) (*LoginResponse, error) {
	// Verify Firebase ID token
	decoded, err := s.verifyFirebaseToken(idToken)
	if err != nil {
		return nil, err
	}

	uid := decoded.UID

	// Check if user already exists
	existingUser, err := s.UserRepo.GetUserByUID(uid)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		if existingUser.IsRegistered {
			return nil, ErrUserAlreadyExists
		}
		// User exists but not fully registered, update registration status
		return s.completeRegistration(uid, deviceInfo, ipAddress)
	}

	// Get user data from Firebase
	userRecord, err := s.FirebaseAuth.GetUser(context.Background(), uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from Firebase: %w", err)
	}

	// Create new user in database
	err = s.UserRepo.InsertUser(
		uid,
		userRecord.DisplayName,
		userRecord.Email,
		userRecord.PhotoURL,
		true, // is_registered = true on first registration
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Record login history
	if err := s.UserRepo.InsertLoginHistory(uid, deviceInfo, ipAddress); err != nil {
		// Log error but don't fail the registration
		fmt.Printf("Warning: failed to record login history: %v\n", err)
	}

	// Create session
	sessionToken, expiresAt, err := s.createSession(uid)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: &UserData{
			UID:     uid,
			Name:    userRecord.DisplayName,
			Email:   userRecord.Email,
			Picture: userRecord.PhotoURL,
		},
		SessionToken: sessionToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Login handles user login with Firebase authentication
func (s *AuthService) Login(idToken string, deviceInfo string, ipAddress string) (*LoginResponse, error) {
	// Verify Firebase ID token
	decoded, err := s.verifyFirebaseToken(idToken)
	if err != nil {
		return nil, err
	}

	uid := decoded.UID

	// Get user from database
	user, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotRegistered
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has completed registration
	if !user.IsRegistered {
		return nil, ErrUserNotRegistered
	}

	// Record login history
	if err := s.UserRepo.InsertLoginHistory(uid, deviceInfo, ipAddress); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: failed to record login history: %v\n", err)
	}

	// Create new session
	sessionToken, expiresAt, err := s.createSession(uid)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: &UserData{
			UID:     user.UID,
			Name:    user.Name,
			Email:   user.Email,
			Picture: user.Picture,
		},
		SessionToken: sessionToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout invalidates a user's session
func (s *AuthService) Logout(sessionToken string) error {
	err := s.UserRepo.DeleteSession(sessionToken)
	if err != nil {
		if err == repository.ErrSessionNotFound {
			return nil // Already logged out
		}
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

// LogoutAllDevices invalidates all sessions for a user
func (s *AuthService) LogoutAllDevices(uid string) error {
	err := s.UserRepo.DeleteAllUserSessions(uid)
	if err != nil {
		return fmt.Errorf("failed to logout from all devices: %w", err)
	}
	return nil
}

// ValidateSession checks if a session is valid and returns the user UID
func (s *AuthService) ValidateSession(sessionToken string) (string, error) {
	valid, uid, err := s.UserRepo.ValidateSession(sessionToken)
	if err != nil {
		return "", fmt.Errorf("failed to validate session: %w", err)
	}

	if !valid {
		return "", errors.New("invalid or expired session")
	}

	return uid, nil
}

// RefreshSession extends a session's expiration time
func (s *AuthService) RefreshSession(sessionToken string) (time.Time, error) {
	// Validate current session
	uid, err := s.ValidateSession(sessionToken)
	if err != nil {
		return time.Time{}, err
	}

	// Check if user still exists and is active
	user, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsRegistered {
		return time.Time{}, ErrUserNotRegistered
	}

	// Extend session expiration
	newExpiresAt := time.Now().Add(24 * time.Hour)
	err = s.UserRepo.UpdateSessionExpiry(sessionToken, newExpiresAt)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to refresh session: %w", err)
	}

	return newExpiresAt, nil
}

// GetUserProfile retrieves user profile information
func (s *AuthService) GetUserProfile(uid string) (*UserData, error) {
	user, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &UserData{
		UID:     user.UID,
		Name:    user.Name,
		Email:   user.Email,
		Picture: user.Picture,
	}, nil
}

// UpdateUserProfile updates user profile information
func (s *AuthService) UpdateUserProfile(uid string, name string, picture string) error {
	// Validate user exists
	_, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to check user: %w", err)
	}

	// Update profile
	err = s.UserRepo.UpdateUserProfile(uid, name, picture)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// GetLoginHistory retrieves user's login history
func (s *AuthService) GetLoginHistory(uid string, limit int) ([]repository.LoginHistory, error) {
	// Validate user exists
	_, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	history, err := s.UserRepo.GetLoginHistory(uid, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get login history: %w", err)
	}

	return history, nil
}

// GetActiveSessions retrieves all active sessions for a user
func (s *AuthService) GetActiveSessions(uid string) ([]repository.Session, error) {
	// Validate user exists
	_, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	sessions, err := s.UserRepo.GetActiveSessions(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	return sessions, nil
}

// DeleteAccount permanently deletes a user account
func (s *AuthService) DeleteAccount(uid string, sessionToken string) error {
	// Validate session belongs to the user
	sessionUID, err := s.ValidateSession(sessionToken)
	if err != nil {
		return err
	}

	if sessionUID != uid {
		return errors.New("unauthorized: session does not match user")
	}

	// Delete user from database (cascade deletes sessions and login history)
	err = s.UserRepo.DeleteUser(uid)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	// Optionally: Delete from Firebase Authentication
	// err = s.FirebaseAuth.DeleteUser(context.Background(), uid)
	// if err != nil {
	//     fmt.Printf("Warning: failed to delete Firebase user: %v\n", err)
	// }

	return nil
}

// CleanupExpiredSessions removes expired sessions (should be run periodically)
func (s *AuthService) CleanupExpiredSessions() (int64, error) {
	count, err := s.UserRepo.CleanupExpiredSessions()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return count, nil
}

// Private helper methods

// verifyFirebaseToken verifies the Firebase ID token
func (s *AuthService) verifyFirebaseToken(idToken string) (*firebase.Token, error) {
	decoded, err := s.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	return decoded, nil
}

// createSession creates a new session for a user
func (s *AuthService) createSession(uid string) (string, time.Time, error) {
	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	err := s.UserRepo.InsertSession(uid, sessionToken, expiresAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: %v", ErrSessionCreation, err)
	}

	return sessionToken, expiresAt, nil
}

// completeRegistration completes registration for existing but unregistered user
func (s *AuthService) completeRegistration(uid string, deviceInfo string, ipAddress string) (*LoginResponse, error) {
	// Update registration status
	err := s.UserRepo.UpdateIsRegistered(uid, true)
	if err != nil {
		return nil, fmt.Errorf("failed to complete registration: %w", err)
	}

	// Get updated user data
	user, err := s.UserRepo.GetUserByUID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Record login history
	if err := s.UserRepo.InsertLoginHistory(uid, deviceInfo, ipAddress); err != nil {
		fmt.Printf("Warning: failed to record login history: %v\n", err)
	}

	// Create session
	sessionToken, expiresAt, err := s.createSession(uid)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: &UserData{
			UID:     user.UID,
			Name:    user.Name,
			Email:   user.Email,
			Picture: user.Picture,
		},
		SessionToken: sessionToken,
		ExpiresAt:    expiresAt,
	}, nil
}