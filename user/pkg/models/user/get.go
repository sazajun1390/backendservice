package user

import (
	"context"

	"github.com/sazajun1390/backendservice/user/pkg/gen/sqlcmodel"
	"github.com/uptrace/bun"
)

func GetAliveUser(
	ctx context.Context,
	db *bun.DB,
	email string,
) (*sqlcmodel.UserProfiles, error) {
	var userProfile sqlcmodel.UserProfiles

	query := GetUserQuery(db, email)
	err := query.Model(&userProfile).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func GetUserQuery(
	db bun.IDB,
	email string,
) *bun.SelectQuery {
	return db.NewSelect().ColumnExpr("user_profiles.*").Join("LEFT JOIN user_actives ON users.id = user_actives.user_id").Join("LEFT JOIN user_provision ON users.id = user_provision.user_id").Where("email = ?", email)
}

func GetDeletedUserQuery(
	db bun.IDB,
	email string,
) *bun.SelectQuery {
	return db.NewSelect().ColumnExpr("user_profiles.*").Join("JOIN deleted_users ON users.id = user_profiles.user_id").Where("email = ?", email)
}
