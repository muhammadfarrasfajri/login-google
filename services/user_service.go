package services

import (
	"errors"
	"strconv"

	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	UserRepo *repository.UserRepository
}

// ------------------------- GET ALL USERS -----------------------------

func (s *UserService) GetAll() ([]models.BaseUser, error) {
	users, err := s.UserRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ------------------------- GET USER BY ID ----------------------------

func (s *UserService) GetByID(id string) (*models.BaseUser, error) {
	id_user, _:= strconv.Atoi(id)
	user, err := s.UserRepo.FindByID(id_user)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// --------------------------- UPDATE USER -----------------------------

func (s *UserService) Update(id, name, email, role, ProfilePicture string) (*models.BaseUser, error) {
	id_user, _:= strconv.Atoi(id)
	// cek apakah user ada
	existing, err := s.UserRepo.FindByID(id_user)
	if err != nil || existing == nil {
		return nil, ErrUserNotFound
	}

	// update field
	existing.Name = name
	existing.Email = email
	existing.Role = role
	existing.ProfilePicture = ProfilePicture

	err = s.UserRepo.Update(*existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// --------------------------- DELETE USER -----------------------------

func (s *UserService) Delete(id string) error {
	// cek user dulu
	id_user, _:= strconv.Atoi(id)
	user, err := s.UserRepo.FindByID(id_user)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	return s.UserRepo.Delete(id)
}
