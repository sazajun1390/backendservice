package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	userv1 "github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1"
	queries "github.com/sazajun1390/backendservice/user/pkg/gen/sqlcmodel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) VerifyUser(
	ctx context.Context,
	req *connect.Request[userv1.VerifyUserRequest],
) (*connect.Response[userv1.VerifyUserResponse], error) {

	db := s.db.DB
	queries := queries.New(db)
	userProfile, err := queries.GetProvisionUser(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if userProfile == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	if req.Msg.VerifyMessage != userProfile[0].UserMultiID.String {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("verify message is incorrect"))
	}

	return &connect.Response[userv1.VerifyUserResponse]{
		Msg: &userv1.VerifyUserResponse{
			User: &userv1.User{
				UserId:          userProfile[0].ResourceID.String,
				UserEmail:       userProfile[0].Email.String,
				UserName:        &userProfile[0].UserMultiID.String,
				UserTel:         &userProfile[0].Tel.String,
				CreatedAt:       timestamppb.New(userProfile[0].CreatedAt.Time),
				UpdatedAt:       timestamppb.New(userProfile[0].UpdatedAt.Time),
				DeletedAt:       timestamppb.New(userProfile[0].DeletedAt.Time),
				PurgedExpiresAt: timestamppb.New(userProfile[0].PurgedExpiresAt.Time),
			},
			UserToken: &userv1.UserToken{
				Token: userProfile[0].UserMultiID.String,
			},
		},
	}, nil
}
