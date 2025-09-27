package userstatus

import (
	"context"
	"time"

	"github.com/sazajun1390/backendservice/user/pkg/gen/sqlcmodel"
	"github.com/uptrace/bun"
)

func CreateProvisionalUser(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
	now time.Time,
) (*sqlcmodel.UserProvision, error) {
	userProvision := &sqlcmodel.UserProvision{
		UserID:    userID,
		CreatedAt: now,
	}
	_, err := idb.NewInsert().Model(userProvision).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return userProvision, nil
}

func createUserActive(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
	now time.Time,
) (*sqlcmodel.UserActives, error) {
	userActive := &sqlcmodel.UserActives{
		UserID:    userID,
		CreatedAt: now,
	}
	_, err := idb.NewInsert().Model(userActive).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return userActive, nil
}

func createUserDelete(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
	now time.Time,
	expiresAt time.Time,
) (*sqlcmodel.UserDeletes, error) {
	userDelete := &sqlcmodel.UserDeletes{
		UserID:    userID,
		CreatedAt: now,
		PurgedExpiresAt: bun.NullTime{
			Time: expiresAt,
		},
	}
	_, err := idb.NewInsert().Model(userDelete).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return userDelete, nil
}
