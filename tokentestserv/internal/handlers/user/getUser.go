package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	userv1 "github.com/sazajun1390/backendservice/tokentestserv/pkg/gen/buf/user/v1"
)

func (s *UserService) GetUserToken(ctx context.Context, req *connect.Request[userv1.GetUserTokenRequest]) (*connect.Response[userv1.GetUserTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("GetUserToken is not implemented"))
}
