package repository

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionNotFound   = errors.New("session not found")
)

type UserRepository struct {
	DB *sql.DB
}

type User struct {
	UID          string    `json:"uid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Picture      string    `json:"picture"`
	IsRegistered bool      `json:"is_registered"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginHistory struct {
	ID         int       `json:"id"`
	UserUID    string    `json:"user_uid"`
	DeviceInfo string    `json:"device_info"`
	IPAddress  string    `json:"ip_address"`
	LoginTime  time.Time `json:"login_time"`
}

type Session struct {
	ID           int       `json:"id"`
	UserUID      string    `json:"user_uid"`
	SessionToken string    `json:"session_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetUserByUID retrieves a user by their Firebase UID
func (r *UserRepository) GetUserByUID(uid string) (*User, error) {
	query := `
		SELECT uid, name, email, picture, is_registered, created_at, updated_at 
		FROM users 
		WHERE uid = ?
	`
	
	var u User
	err := r.DB.QueryRow(query, uid).Scan(
		&u.UID,
		&u.Name,
		&u.Email,
		&u.Picture,
		&u.IsRegistered,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	
	return &u, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT uid, name, email, picture, is_registered, created_at, updated_at 
		FROM users 
		WHERE email = ?
	`
	
	var u User
	err := r.DB.QueryRow(query, email).Scan(
		&u.UID,
		&u.Name,
		&u.Email,
		&u.Picture,
		&u.IsRegistered,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	
	return &u, nil
}

// UserExists checks if a user exists by UID
func (r *UserRepository) UserExists(uid string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE uid = ?)"
	err := r.DB.QueryRow(query, uid).Scan(&exists)
	return exists, err
}

// InsertUser creates a new user record
func (r *UserRepository) InsertUser(uid, name, email, picture string, isRegistered bool) error {
	// Check if user already exists
	exists, err := r.UserExists(uid)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExists
	}
	
	query := `
		INSERT INTO users (uid, name, email, picture, is_registered, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	
	_, err = r.DB.Exec(query, uid, name, email, picture, isRegistered)
	return err
}

// UpdateIsRegistered updates the registration status of a user
func (r *UserRepository) UpdateIsRegistered(uid string, isRegistered bool) error {
	query := `
		UPDATE users 
		SET is_registered = ?, updated_at = NOW() 
		WHERE uid = ?
	`
	
	result, err := r.DB.Exec(query, isRegistered, uid)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// UpdateUserProfile updates user profile information
func (r *UserRepository) UpdateUserProfile(uid, name, picture string) error {
	query := `
		UPDATE users 
		SET name = ?, picture = ?, updated_at = NOW() 
		WHERE uid = ?
	`
	
	result, err := r.DB.Exec(query, name, picture, uid)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// DeleteUser removes a user from the database
func (r *UserRepository) DeleteUser(uid string) error {
	// Start transaction to delete user and related data
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete sessions
	_, err = tx.Exec("DELETE FROM sessions WHERE user_uid = ?", uid)
	if err != nil {
		return err
	}
	
	// Delete login history
	_, err = tx.Exec("DELETE FROM login_history WHERE user_uid = ?", uid)
	if err != nil {
		return err
	}
	
	// Delete user
	result, err := tx.Exec("DELETE FROM users WHERE uid = ?", uid)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return tx.Commit()
}

// InsertLoginHistory records a login attempt
func (r *UserRepository) InsertLoginHistory(uid, deviceInfo, ipAddress string) error {
	query := `
		INSERT INTO login_history (user_uid, device_info, ip_address, login_time) 
		VALUES (?, ?, ?, NOW())
	`
	
	_, err := r.DB.Exec(query, uid, deviceInfo, ipAddress)
	return err
}

// GetLoginHistory retrieves login history for a user
func (r *UserRepository) GetLoginHistory(uid string, limit int) ([]LoginHistory, error) {
	query := `
		SELECT id, user_uid, device_info, ip_address, login_time 
		FROM login_history 
		WHERE user_uid = ? 
		ORDER BY login_time DESC 
		LIMIT ?
	`
	
	rows, err := r.DB.Query(query, uid, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var history []LoginHistory
	for rows.Next() {
		var h LoginHistory
		err := rows.Scan(&h.ID, &h.UserUID, &h.DeviceInfo, &h.IPAddress, &h.LoginTime)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	
	return history, rows.Err()
}

// InsertSession creates a new session
func (r *UserRepository) InsertSession(uid, sessionToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (user_uid, session_token, created_at, expires_at) 
		VALUES (?, ?, NOW(), ?)
	`
	
	_, err := r.DB.Exec(query, uid, sessionToken, expiresAt)
	return err
}

// GetSessionByToken retrieves a session by token
func (r *UserRepository) GetSessionByToken(sessionToken string) (*Session, error) {
	query := `
		SELECT id, user_uid, session_token, created_at, expires_at 
		FROM sessions 
		WHERE session_token = ? AND expires_at > NOW()
	`
	
	var s Session
	err := r.DB.QueryRow(query, sessionToken).Scan(
		&s.ID,
		&s.UserUID,
		&s.SessionToken,
		&s.CreatedAt,
		&s.ExpiresAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	
	return &s, nil
}

// ValidateSession checks if a session is valid
func (r *UserRepository) ValidateSession(sessionToken string) (bool, string, error) {
	session, err := r.GetSessionByToken(sessionToken)
	if err != nil {
		if err == ErrSessionNotFound {
			return false, "", nil
		}
		return false, "", err
	}
	
	return true, session.UserUID, nil
}

// DeleteSession removes a session (logout)
func (r *UserRepository) DeleteSession(sessionToken string) error {
	query := "DELETE FROM sessions WHERE session_token = ?"
	
	result, err := r.DB.Exec(query, sessionToken)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrSessionNotFound
	}
	
	return nil
}

// DeleteAllUserSessions removes all sessions for a user (logout all devices)
func (r *UserRepository) DeleteAllUserSessions(uid string) error {
	query := "DELETE FROM sessions WHERE user_uid = ?"
	_, err := r.DB.Exec(query, uid)
	return err
}

// CleanupExpiredSessions removes expired sessions from database
func (r *UserRepository) CleanupExpiredSessions() (int64, error) {
	query := "DELETE FROM sessions WHERE expires_at < NOW()"
	
	result, err := r.DB.Exec(query)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}

// GetActiveSessions retrieves all active sessions for a user
func (r *UserRepository) GetActiveSessions(uid string) ([]Session, error) {
	query := `
		SELECT id, user_uid, session_token, created_at, expires_at 
		FROM sessions 
		WHERE user_uid = ? AND expires_at > NOW() 
		ORDER BY created_at DESC
	`
	
	rows, err := r.DB.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sessions []Session
	for rows.Next() {
		var s Session
		err := rows.Scan(&s.ID, &s.UserUID, &s.SessionToken, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	
	return sessions, rows.Err()
}

// UpdateSessionExpiry extends a session's expiration time
func (r *UserRepository) UpdateSessionExpiry(sessionToken string, newExpiresAt time.Time) error {
	query := "UPDATE sessions SET expires_at = ? WHERE session_token = ?"
	
	result, err := r.DB.Exec(query, newExpiresAt, sessionToken)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrSessionNotFound
	}
	
	return nil
}