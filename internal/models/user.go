package models

import (
	"database/sql"
	"errors"
	_ "errors"
	"strings"
	_ "strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int
	Name           string
	Email          string
	Hashedpassword []byte
	Created        time.Time
}

type UserModelInterface interface {
	Insert(name string, email string, password string) error
	Authenticate(email string, password string) (int, error)
	Exists(id int) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

// Get by id
// Insert a user
//Authenticate

func (um *UserModel) Insert(name string, email string, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
			VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = um.DB.Exec(stmt, name, email, string(hashedPass))

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

func (um *UserModel) Authenticate(email string, password string) (int, error) {

	var id int
	var hashedPass []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := um.DB.QueryRow(stmt, email).Scan(&id, &hashedPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPass, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (um *UserModel) Exists(id int) (bool, error) {

	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := um.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
