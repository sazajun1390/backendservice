package pkg

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"github.com/quic-go/quic-go/http3"

	"connectrpc.com/connect"
	v1 "github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1"
	userv1connect "github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1/userv1connect"
	"github.com/uptrace/bun"
)

type Proxy struct {
	client   userv1connect.UserServiceClient
	shutdown func(context.Context) error
	db       *bun.DB
}

var _ userv1connect.UserServiceHandler = (*Proxy)(nil)

func (p *Proxy) CreateUser(ctx context.Context, req *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreateUserResponse], error) {
	return p.client.CreateUser(ctx, req)
}

func (p *Proxy) GetUserToken(ctx context.Context, req *connect.Request[v1.GetUserTokenRequest]) (*connect.Response[v1.GetUserTokenResponse], error) {
	header := req.Header().Clone()
	token := header.Get("Authorization")
	if token == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token is required"))
	}

	return p.client.GetUserToken(ctx, req)
}

func NewProxy(db *bun.DB) *Proxy {
	h3Tansport := &http3.Transport{
		TLSClientConfig: &tls.Config{
			// we need this because our certificate is self signed
			InsecureSkipVerify: true,
		},
	}
	shutdown := func(ctx context.Context) error {
		return h3Tansport.Close()
	}
	h3Client := http.Client{
		Transport: h3Tansport,
	}

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	userServiceClient := userv1connect.NewUserServiceClient(&h3Client, userServiceURL)

	return &Proxy{client: userServiceClient, shutdown: shutdown, db: db}
}
