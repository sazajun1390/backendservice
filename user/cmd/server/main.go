package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sazajun1390/backendservice/user/internal/handlers/user"
	"github.com/sazajun1390/backendservice/user/pkg/gen/buf/user/v1/userv1connect"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.Info("starting server")
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
	}

	slog.Info("database ping success")

	usersrv := user.NewUserService(db)
	mux := http.NewServeMux()

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		slog.Error("failed to create otel interceptor", "error", err)
	}

	path, handler := userv1connect.NewUserServiceHandler(
		usersrv,
		connect.WithInterceptors(
			otelInterceptor,
		),
	)
	mux.Handle(path, handler)

	err = http3.ListenAndServeTLS(":51051", "/Users/juntendou/Documents/test-kotlin-connect/tokentestserv/.data/masterCert/cert.pem", "/Users/juntendou/Documents/test-kotlin-connect/tokentestserv/.data/masterCert/key.pem", mux)
	if err != nil {
		slog.Error("failed to listen and serve tls", "error", err)
	}

	http.ListenAndServe(
		":51051",
		h2c.NewHandler(mux, &http2.Server{}),
	)

}
