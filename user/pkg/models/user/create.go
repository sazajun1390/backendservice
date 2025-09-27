package user

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	sqlcmodel "github.com/sazajun1390/backendservice/user/pkg/gen/sqlcmodel"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(
	ctx context.Context,
	idb bun.IDB,
	email string,
	password string,
	now time.Time,
) (
	*sqlcmodel.Users,
	*sqlcmodel.UserProfiles,
	error,
) {
	userMaster := &sqlcmodel.Users{
		CreatedAt: now,
	}
	_, err := idb.NewInsert().Model(userMaster).Exec(ctx)
	if err != nil {
		return nil, nil, err
	}

	userprovision := &sqlcmodel.UserProvision{
		UserID:    userMaster.ID,
		CreatedAt: now,
	}
	_, err = idb.NewInsert().Model(userprovision).Exec(ctx)
	if err != nil {
		return nil, nil, err
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	userProfile := &sqlcmodel.UserProfiles{
		UserID:     userMaster.ID,
		ResourceID: ulid.Make().String(),
		Email:      email,
		Password:   string(passhash),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = idb.NewInsert().Model(userProfile).Exec(ctx)
	if err != nil {
		return nil, nil, err
	}

	return userMaster, userProfile, nil
}
