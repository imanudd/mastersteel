package chatapp

import (
	"context"
	"finalproject/config"

	"github.com/go-pg/pg/v10"
)

type Service struct {
	cfg  *config.Config
	db   *pg.DB
	repo IRepository
}

func NewService(cfg *config.Config, db *pg.DB, repo IRepository) Service {
	return Service{cfg, db, repo}
}

func (s Service) SendMessage(ctx context.Context) (*Users, error) {

	return s.repo.SendMessage(ctx)
}
