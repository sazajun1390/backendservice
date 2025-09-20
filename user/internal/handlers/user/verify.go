package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	userv1 "github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1"
	queries "github.com/sazajun1390/backendservice/user/pkg/gen/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) VerifyUser(ctx context.Context, req *connect.Request[userv1.VerifyUserRequest]) (*connect.Response[userv1.VerifyUserResponse], error) {

	db := s.db
	queries := queries.New(db)
	userProfile, err := queries.GetProvisionUser(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if userProfile == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	return &connect.Response[userv1.VerifyUserResponse]{
		Msg: &userv1.VerifyUserResponse{
			User: &userv1.User{
				UserId:          userProfile[0].ResourceID,
				UserEmail:       userProfile[0].Email,
				UserName:        &userProfile[0].UserMultiID,
				UserTel:         &userProfile[0].UserTel,
				CreatedAt:       timestamppb.New(userProfile[0].CreatedAt),
				UpdatedAt:       timestamppb.New(userProfile[0].UpdatedAt),
				DeletedAt:       timestamppb.New(userProfile[0].DeletedAt.Time),
				PurgedExpiresAt: timestamppb.New(userProfile[0].PurgedExpiresAt.Time),
			},
			UserToken: &userv1.UserToken{
				Token: userProfile[0].UserMultiID,
			},
		},
	}, nil
}
