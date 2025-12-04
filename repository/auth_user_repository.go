package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FindByGoogleUID(uid string) (*models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture FROM users WHERE google_uid = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, uid)
	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.GooglePicture)
	if err != nil {
		if err == sql.ErrNoRows {
		return nil, nil
	}
	return nil, err
}
	return &user, nil
}

func (r *UserRepository) Create(user models.BaseUser) error {
	sqlQuery := `INSERT INTO users (google_uid, name, email, google_picture) VALUES (?, ?, ?, ?)`
	_, err := r.DB.Exec(sqlQuery, user.GoogleUID, user.Name, user.Email, user.GooglePicture)
	return err
}

func (r *UserRepository) UpdateLoginStatus(id int, status int) error {
    query := `UPDATE users SET is_logged_in = ? WHERE id = ?`
    _, err := r.DB.Exec(query, status, id)
    return err
}

func (r *UserRepository) SaveLoginHistory(userID int, deviceInfo, ip string) error {
	sqlQuery := `INSERT INTO login_history_user (user_id, login_at, device_info, ip_address) VALUES (?, NOW(), ?, ?)`
	_, err := r.DB.Exec(sqlQuery, userID, deviceInfo, ip)
	return err
}
