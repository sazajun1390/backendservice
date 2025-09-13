package user

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sazajun1390/backendservice/user/pkg/gen/user"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(
	ctx context.Context,
	idb bun.IDB,
	email string,
	password string,
	now time.Time,
) (*user.Users, error) {
	userMaster := &user.Users{
		CreatedAt: now,
	}
	_, err := idb.NewInsert().Model(userMaster).Exec(ctx)
	if err != nil {
		return nil, err
	}

	userprovision := &user.UserProvision{
		UserID:    userMaster.ID,
		CreatedAt: now,
	}
	_, err = idb.NewInsert().Model(userprovision).Exec(ctx)
	if err != nil {
		return nil, err
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userProfile := &user.UserProfiles{
		UserID:     userMaster.ID,
		ResourceID: ulid.Make().String(),
		Email:      email,
		Password:   string(passhash),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = idb.NewInsert().Model(userProfile).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return userMaster, nil
}
