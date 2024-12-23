package dbrepo

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/repository"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *pgxpool.Pool
}

func NewPostgresRepo(conn *pgxpool.Pool, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  conn,
	}
}
