package userstatus

import (
	"context"
	"time"

	"github.com/sazajun1390/backendservice/tokentestserv/pkg/gen/user"
	"github.com/uptrace/bun"
)

func CreateProvisionalUser(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
	now time.Time,
) (*user.UserProvision, error) {
	userProvision := &user.UserProvision{
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
) (*user.UserActives, error) {
	userActive := &user.UserActives{
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
) (*user.UserDeletes, error) {
	userDelete := &user.UserDeletes{
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
