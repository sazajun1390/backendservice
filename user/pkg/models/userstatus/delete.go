package userstatus

import (
	"context"

	"github.com/sazajun1390/backendservice/user/pkg/gen/sqlcmodel"
	"github.com/uptrace/bun"
)

func deleteUserActive(
	ctx context.Context,
	idb bun.IDB,
	userID int64,
) error {
	_, err := idb.NewDelete().Model((*sqlcmodel.UserActives)(nil)).Where("user_id = ?", userID).Exec(ctx)
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
	_, err := idb.NewDelete().Model((*sqlcmodel.UserProvision)(nil)).Where("user_id = ?", userID).Exec(ctx)
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
	_, err := idb.NewDelete().Model((*sqlcmodel.UserDeletes)(nil)).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
