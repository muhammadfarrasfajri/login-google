package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)


type AdminRepository struct {
	DB *sql.DB
}

func (r *AdminRepository) FindByGoogleUID(uid string) (*models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, is_logged_in FROM admins WHERE google_uid = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, uid)
	admin := models.BaseUser{}
	err := row.Scan(&admin.ID, &admin.GoogleUID, &admin.Name, &admin.Email, &admin.GooglePicture, &admin.IsLoggedIn)
	if err != nil {
		if err == sql.ErrNoRows {
		return nil, err
	}
	return nil, err
}
	return &admin, err
}

func (r *AdminRepository) Create(admin models.BaseUser) error {
	sqlQuery := `INSERT INTO admins (google_uid, name, email, google_picture) VALUES (?, ?, ?, ?)`
	_, err := r.DB.Exec(sqlQuery, admin.GoogleUID, admin.Name, admin.Email, admin.GooglePicture)
	return err
}

func (r *AdminRepository) UpdateLoginStatus(id int, status int) error {
    query := `UPDATE admins SET is_logged_in = ? WHERE id = ?`
    _, err := r.DB.Exec(query, status, id)
    return err
}

func (r *AdminRepository) SaveLoginHistory(adminID int, deviceInfo, ip string) error {
	sqlQuery := `INSERT INTO login_history_admin (admin_id, login_at, device_info, ip_address) VALUES (?, NOW(), ?, ?)`
	_, err := r.DB.Exec(sqlQuery, adminID, deviceInfo, ip)
	return err
}

func (r *AdminRepository) FindByID(id int) (*models.BaseUser, error) {
	row := r.DB.QueryRow(`
        SELECT id, google_uid, name, email, google_picture, role, profile_picture
        FROM users WHERE id = ?
    `, id)

	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.GooglePicture, &user.Role, &user.ProfilePicture)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return &user, err
}
