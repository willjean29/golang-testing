package repository

import (
	"database/sql"
	"webapp/pkg/data"
)

type Repository interface {
	Connection() *sql.DB
	AllUsers() ([]*data.User, error)
	GetUser(id int) (*data.User, error)
	GetUserByEmail(email string) (*data.User, error)
	UpdateUser(user *data.User) error
	DeleteUser(id int) error
	InsertUser(user data.User) (int, error)
	ResetPassword(id int, password string) error
	InsertUserImage(i data.UserImage) (int, error)
}
