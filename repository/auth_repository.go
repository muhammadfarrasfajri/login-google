package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FindByGoogleUID(uid string) (*models.User, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, role FROM users WHERE google_uid = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, uid)
	user := models.User{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.Google_picture, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

func (r *UserRepository) Create(user models.User) error {
	sqlQuery := `INSERT INTO users (google_uid, name, email, google_picture, role) VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(sqlQuery, user.GoogleUID, user.Name, user.Email, user.Google_picture, user.Role)
	return err
}

func (r *UserRepository) SaveLoginHistory(userID int, deviceInfo, ip string) error {
	_, err := r.DB.Exec(`
        INSERT INTO login_history (user_id, login_at, device_info, ip_address)
        VALUES (?, NOW(), ?, ?)
    `, userID, deviceInfo, ip)

	return err
}
