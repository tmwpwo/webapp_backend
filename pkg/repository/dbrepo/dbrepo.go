package dbrepo

import (
	"database/sql"
	"go_server/pkg/config"
	"go_server/pkg/repository"
)

type postgresDBrepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBrepo{
		App: a,
		DB:  conn,
	}
}
