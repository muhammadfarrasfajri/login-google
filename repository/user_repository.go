package repository

import (
	"database/sql"
	"fmt"

	"github.com/muhammadfarrasfajri/login-google/models"
)

// --------------------------- GET ALL USERS -----------------------------------

func (r *UserRepository) GetAll() ([]models.User, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, role, profile_picture FROM users`
	rows, err := r.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.GoogleUID, &u.Name, &u.Email, &u.Google_picture, &u.Role, &u.Profile_picture)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// --------------------------- FIND BY ID --------------------------------------

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	row := r.DB.QueryRow(`
        SELECT id, google_uid, name, email, google_picture, role, profile_picture
        FROM users WHERE id = ?
    `, id)

	user := models.User{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.Google_picture, &user.Role, &user.Profile_picture)

	if err == sql.ErrNoRows {
		fmt.Println("No user found with ID:", err)
		return nil, nil
	}

	return &user, err
}

// --------------------------- UPDATE USER -------------------------------------

func (r *UserRepository) Update(user models.User) error {
	_, err := r.DB.Exec(`
        UPDATE users SET 
            name = ?, 
            email = ?, 
            role = ?,
			profile_picture = ?
        WHERE id = ?
    `, user.Name, user.Email, user.Role, user.Profile_picture, user.ID)

	return err
}

// --------------------------- DELETE USER -------------------------------------

func (r *UserRepository) Delete(id string) error {
	// Hapus login history dulu
	_, err := r.DB.Exec("DELETE FROM login_history WHERE user_id = ?", id)
	if err != nil {
		return err
	}

	// Baru hapus user
	_, err = r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// --------------------------- UPDATE PHOTO URL -------------------------------------
func (r *UserRepository) UpdatePhotoURL(userID int, url string) error {
	// Update kolom avatar_url di tabel users
	_, err := r.DB.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", url, userID)
	fmt.Println("errornya:", err)
	return err

}
