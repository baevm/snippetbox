package mocks

import (
	"snippetbox/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Get(id int) (*models.User, error) {
	if id == 1 {
		u := &models.User{
			Id:      1,
			Name:    "User",
			Email:   "user@test.com",
			Created: time.Now(),
		}

		return u, nil
	}

	return nil, models.ErrNoRecord
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "password" {
			return models.ErrInvalidCredentials
		}
		return nil
	}

	return models.ErrNoRecord
}

func (m *UserModel) Create(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "user@test.com" && password == "password" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
