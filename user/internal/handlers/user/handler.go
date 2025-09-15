package user

import (
	"github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1/userv1connect"
	"github.com/uptrace/bun"
)

type UserService struct {
	db *bun.DB
}

var _ userv1connect.UserServiceHandler = (*UserService)(nil)

func NewUserService(db *bun.DB) *UserService {
	return &UserService{db: db}
}
