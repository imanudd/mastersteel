package chatapp

import (
	"context"
	"finalproject/config"

	"github.com/go-pg/pg/v10"
)

type Repository struct {
	db  *pg.DB
	cfg *config.Config
}

func NewRepository(db *pg.DB, cfg *config.Config) Repository {
	return Repository{db, cfg}
}

func (r Repository) SendMessage(ctx context.Context) (*Users, error) {
	var user Users
	err := r.db.Model(&user).Select()
	if err != nil {
		return nil, err
	}

	return &user, nil
}
