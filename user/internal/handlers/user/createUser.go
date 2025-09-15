package user

import (
	"context"
	"log/slog"
	"math/rand"
	"os"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/oklog/ulid/v2"
	userv1 "github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1"
	"github.com/sazajun1390/backendservice/user/pkg/models/user"

	"github.com/cockroachdb/errors"
	"github.com/dimonomid/clock"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/uptrace/bun"
)

func (s *UserService) CreateUser(ctx context.Context, req *connect.Request[userv1.CreateUserRequest]) (*connect.Response[userv1.CreateUserResponse], error) {

	profileUserQuery, err := user.GetAliveUser(ctx, s.db, req.Msg.GetUserEmail())
	if err != nil {
		slog.WarnContext(ctx, "failed to get user", slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if profileUserQuery != nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("user already exists"))
	}

	// 時刻設定
	now := clock.New().Now()
	// ULID生成
	entropy := rand.New(rand.NewSource(now.UnixNano()))
	ulid, err := ulid.New(uint64(now.Unix()), entropy)
	if err != nil {
		slog.WarnContext(ctx, "failed to generate ULID", slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// sendgridに
	from := mail.NewEmail("Example User", os.Getenv("SENDGRID_FROM_EMAIL"))
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", req.Msg.GetUserEmail())
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	// client.Request, _ = sendgrid.SetDataResidency(client.Request, "eu")
	// uncomment the above line if you are sending mail using a regional EU subuser
	response, err := client.Send(message)
	if err != nil {
		slog.WarnContext(ctx, "failed to send email", slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	slog.InfoContext(ctx, "email sent successfully",
		slog.Int("status_code", response.StatusCode),
		slog.String("body", response.Body),
		slog.Any("headers", response.Headers),
	)

	// ormの都合もあり、トランザクションで設定
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		userMaster, err := user.CreateUser(ctx, tx, req.Msg.GetUserEmail(), req.Msg.GetPassword(), now)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
		slog.InfoContext(ctx, "user created", slog.Any("user", userMaster))
		return nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &userv1.CreateUserResponse{
		User: &userv1.User{
			UserId:    "users/" + ulid.String(),
			UserEmail: req.Msg.GetUserEmail(),
			CreatedAt: timestamppb.New(now),
		},
	}

	return connect.NewResponse(res), nil
}
