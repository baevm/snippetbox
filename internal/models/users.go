package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserRepo interface {
	Create(title, content string, expires int) error
	Exists(id int) (bool, error)
	Authenticate(email, password string) (int, error)
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Create(name, email, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	// ? used as placeholder to avoid SQL injections
	query := `
	INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())
	`
	_, err = u.DB.Exec(query, name, email, string(hashedPass))

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	user := &User{}

	query := `
	
	`
	err := u.DB.
		QueryRow(query, id).
		Scan(&user.Id, &user.Name, &user.Email, &user.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrNoRecord
		} else {
			return false, err
		}
	}

	return true, nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := `
	SELECT id, hashed_password from users where email = ?
	`
	err := u.DB.
		QueryRow(query, email).
		Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
