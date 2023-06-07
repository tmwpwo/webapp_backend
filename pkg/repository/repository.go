package repository

import "go_server/pkg/models"

type DatabaseRepo interface {
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authentication(email, testPassword string) (int, string, error)
}
