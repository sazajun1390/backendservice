package user

import (
	"github.com/uptrace/bun"
)

type UserService struct {
	db *bun.DB
}

func NewUserService(db *bun.DB) *UserService {
	return &UserService{db: db}
}
