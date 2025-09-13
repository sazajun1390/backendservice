package migrations

import (
	"context"

	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time `bun:"created_at,notnull"`

	UserProfiles *UserProfile `bun:"rel:has-one,join:id=user_id"`
	UserActives  *UserActive  `bun:"rel:has-one,join:id=user_id"`
}

type UserProfile struct {
	bun.BaseModel  `bun:"table:user_profiles,alias:profiles"`
	UserID         int64          `bun:"user_id,notnull,unique"`
	UserMultiID    string         `bun:"user_multi_id,notnull,unique"`
	ResourceID     string         `bun:"resource_id,notnull,unique"`
	Email          string         `bun:"email,unique,notnull"`
	Password       string         `bun:"password,notnull"`
	PostCode       string         `bun:"post_code,notnull"`
	Address        int64          `bun:"address,notnull"`
	AddressKana    string         `bun:"address_kana,notnull"`
	Nonce          string         `bun:"nonce,notnull"`
	Tel            sql.NullString `bun:"tel"`
	CreatedAt      time.Time      `bun:"created_at,notnull"`
	UpdatedAt      time.Time      `bun:"updated_at,notnull"`
	DeletedAt      bun.NullTime   `bun:"deleted_at,soft_delete"`
	PurgeExpiredAt bun.NullTime   `bun:"purged_expires_at"`

	// User is the user that this profile belongs to.
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

// Provisionテーブルはユーザーのプロビジョニング状態を管理するテーブルです。
type UserProvision struct {
	bun.BaseModel `bun:"table:user_provision"`
	UserID        int64     `bun:"user_id,notnull,unique"`
	CreatedAt     time.Time `bun:"created_at,notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

// Activeテーブルはユーザーのアクティブ状態を管理するテーブルです。
type UserActive struct {
	bun.BaseModel `bun:"table:user_actives"`
	UserID        int64     `bun:"user_id,notnull,unique"`
	CreatedAt     time.Time `bun:"created_at,notnull"`

	// User is the user that this active belongs to.
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

// Deleteテーブルは論理削除用のテーブルです。
type UserDelete struct {
	bun.BaseModel   `bun:"table:user_deletes"`
	UserID          int64        `bun:"user_id,notnull,unique"`
	CreatedAt       time.Time    `bun:"created_at,notnull"`
	PargedExpiresAt bun.NullTime `bun:"purged_expires_at"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		var err error
		_, err = db.NewCreateTable().Model((*User)(nil)).WithForeignKeys().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*UserProfile)(nil)).WithForeignKeys().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*UserActive)(nil)).WithForeignKeys().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*UserProvision)(nil)).WithForeignKeys().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*UserDelete)(nil)).WithForeignKeys().Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().Model((*UserActive)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model((*UserProfile)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model((*UserProvision)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model((*UserDelete)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
