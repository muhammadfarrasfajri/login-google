package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FindByGoogleUID(uid string) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, google_uid, name, email, picture, role FROM users WHERE google_uid = ?", uid)

	user := models.User{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.Picture, &user.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepository) Create(user models.User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (google_uid, name, email, picture) VALUES (?, ?, ?, ?)",
		user.GoogleUID, user.Name, user.Email, user.Picture,
	)
	return err
}

func (r *UserRepository) SaveLoginHistory(userID int, deviceInfo, ip string) error {
    _, err := r.DB.Exec(`
        INSERT INTO login_history (user_id, login_at, device_info, ip_address)
        VALUES (?, NOW(), ?, ?)
    `, userID, deviceInfo, ip)

    return err
}

