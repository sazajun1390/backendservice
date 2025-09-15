package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/exaring/otelpgx"
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
	config.Tracer = otelpgx.NewTracer()
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

	h2srv := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
		Addr:    ":8080",
	}

	h3srv := &http3.Server{
		Handler: mux,
		Addr:    ":51051",
	}

	go func() {
		err = h3srv.ListenAndServeTLS("/Users/juntendou/Documents/test-kotlin-connect/tokentestserv/.data/masterCert/cert.pem", "/Users/juntendou/Documents/test-kotlin-connect/tokentestserv/.data/masterCert/key.pem")
		if err != nil {
			slog.Error("failed to listen and serve h3", "error", err)
		}
	}()

	go func() {
		err = h2srv.ListenAndServe()
		if err != nil {
			slog.Error("failed to listen and serve h2", "error", err)
		}
	}()
}
