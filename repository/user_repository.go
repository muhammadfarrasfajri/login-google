package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)

// --------------------------- GET ALL USERS -----------------------------------

func (r *UserRepository) GetAll() ([]models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, picture, role FROM users`
	rows, err := r.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.BaseUser{}

	for rows.Next() {
		u := models.BaseUser{}
		err := rows.Scan(&u.ID, &u.GoogleUID, &u.Name, &u.Email, &u.GooglePicture,&u.ProfilePicture, &u.Role)
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

func (r *UserRepository) FindByID(id string) (*models.BaseUser, error) {
	row := r.DB.QueryRow(`
        SELECT id, google_uid, name, email, picture, role
        FROM users WHERE id = ?
    `, id)

	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.GooglePicture,&user.ProfilePicture, &user.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

// --------------------------- UPDATE USER -------------------------------------

func (r *UserRepository) Update(user models.BaseUser) error {
	_, err := r.DB.Exec(`
        UPDATE users SET 
            name = ?, 
            email = ?, 
            role = ?
        WHERE id = ?
    `, user.Name, user.Email, user.Role, user.ID)

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
