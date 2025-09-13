package user

import (
	"context"
	"time"

	"github.com/sazajun1390/backendservice/user/pkg/gen/user"
	"github.com/sazajun1390/backendservice/user/pkg/models/userstatus"
	"github.com/uptrace/bun"
)

func RevivalUser(
	ctx context.Context,
	tx bun.Tx,
	email string,
	now time.Time,
) error {
	var userProfile user.UserProfiles
	query := GetDeletedUserQuery(tx, email)
	err := query.Model(&userProfile).Scan(ctx)
	if err != nil {
		return err
	}

	_, err = userstatus.RevivalUserStatus(ctx, tx, userProfile.UserID, now)
	if err != nil {
		return err
	}

	return nil
}
