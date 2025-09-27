package gen

//go:generate mockgen -source=./buf/user/v1/userv1connect/api.connect.go -destination=./moc/user/mock_gen.go -package=user UserService
//go:generate sqlc generate -f ../../deployments/sqlc.yml
//go:generate go run ../../cmd/patch-models/main.go ./sqlcmodel
