package userstatus

import (
	"context"

	"github.com/sazajun1390/backendservice/tokentestserv/pkg/gen/user"
	"github.com/uptrace/bun"
)

func deleteUserActive(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
) error {
	_, err := idb.NewDelete().Model((*user.UserActives)(nil)).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func deleteUserProvision(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
) error {
	_, err := idb.NewDelete().Model((*user.UserProvision)(nil)).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func deleteUserDelete(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
) error {
	_, err := idb.NewDelete().Model((*user.UserDeletes)(nil)).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
