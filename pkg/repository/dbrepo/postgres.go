package dbrepo

import (
	"context"
	"errors"
	"go_server/pkg/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBrepo) AllUsers() bool {
	return true
}

func (m *postgresDBrepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select id, company, email, password, access_level from users where id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.Company,
		&u.Email,
		&u.Password,
		&u.Access_level,
	)

	if err != nil {
		return u, err
	}
	return u, nil
}

// updates a user in the database
func (m *postgresDBrepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "update users set company = $1, email = $2, access_level = $3 "

	_, err := m.DB.ExecContext(ctx, query,
		u.Company, u.Email, u.Access_level)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBrepo) Authentication(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPass string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)

	err := row.Scan(&id, &hashedPass)

	if err != nil {
		return id, "dupa", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrenct password")

	} else if err != nil {
		return 0, "cipa", err
	}

	return id, hashedPass, nil

}
