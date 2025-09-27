package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	userv1 "github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1"
)

func (s *UserService) GetUser(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (*connect.Response[userv1.GetUserResponse], error) {

	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("GetUserToken is not implemented"))
}
