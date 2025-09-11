package userstatus

import (
	"context"
	"database/sql"
	"time"

	"github.com/sazajun1390/backendservice/tokentestserv/pkg/gen/user"
	"github.com/uptrace/bun"
)

func RevivalUserStatus(
	ctx context.Context,
	tx bun.Tx,
	userID int64,
	now time.Time,
) (*user.UserActives, error) {
	err := deleteUserDelete(ctx, tx, userID)
	// 一件も出ない場合にはsql.ErrNoRowsが返るはず
	if err != nil {
		return nil, err
	}
	userActive, err := createUserActive(ctx, tx, userID, now)
	if err != nil {
		return nil, err
	}
	return userActive, nil
}

func PargeUser(
	ctx context.Context,
	tx bun.Tx,
	userID int64,
	now time.Time,
) error {
	exists, err := tx.NewSelect().Model((*user.UserDeletes)(nil)).Where("user_id = ?", userID).Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// 削除ユーザーの削除処理
	err = deleteUserDelete(ctx, tx, userID)
	if err != nil {
		return err
	}
	return nil
}

func ActivateUserStatus(
	ctx context.Context,
	tx bun.Tx,
	userID int64,
	now time.Time,
) (*user.UserActives, error) {
	err := deleteUserProvision(ctx, tx, userID)
	// 一件も出ない場合にはsql.ErrNoRowsが返るはず
	if err != nil {
		return nil, err
	}
	userActive, err := createUserActive(ctx, tx, userID, now)
	if err != nil {
		return nil, err
	}
	return userActive, nil
}

func WithdrawalUserStatus(
	ctx context.Context,
	tx bun.Tx,
	userID int64,
	now time.Time,
	purgedExpiresAt time.Time,
) (*user.UserDeletes, error) {
	err := deleteUserActive(ctx, tx, userID)
	// 一件も出ない場合にはsql.ErrNoRowsが返るはず
	if err != nil {
		return nil, err
	}
	userDelete, err := createUserDelete(ctx, tx, userID, now, purgedExpiresAt)
	if err != nil {
		return nil, err
	}
	return userDelete, nil
}
