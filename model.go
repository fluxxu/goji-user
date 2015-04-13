package user

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/fluxxu/util"
	"github.com/lann/squirrel"
	"time"
)

func hashPassword(plain string) string {
	hasher := sha256.New()
	hasher.Write([]byte(plain))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

type User struct {
	Id          int64         `db:"id" json:"id" provider:"sort filter"`
	Email       string        `db:"email" json:"email" provider:"sort filter"`
	Password    string        `db:"password" json:"-"`
	DisplayName string        `db:"display_name" json:"display_name" provider:"sort filter"`
	CreatedAt   util.NullTime `db:"created_at" json:"created_at" provider:"sort"`
	UpdatedAt   util.NullTime `db:"updated_at" json:"updated_at" provider:"sort"`

	Roles []string `json:"roles"`
}

func (u *User) SetPassword(password string) *User {
	u.Password = hashPassword(password)
	return u
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetDisplayName() string {
	return u.DisplayName
}

func (u *User) validateInsert() (*util.ValidationContext, error) {
	ctx := util.NewValidationContext()
	if !govalidator.IsEmail(u.Email) {
		ctx.AddError("email", "Invalid Email")
	} else {
		var count int
		if err := dbx.Get(&count, "SELECT COUNT(*) FROM user WHERE email = ?", u.Email); err != nil {
			return nil, err
		}
		if count != 0 {
			ctx.AddError("email", "Email already used")
		}
	}

	if u.DisplayName == "" {
		ctx.AddError("display_name", "Display name is required")
	}

	return ctx, nil
}

func (u *User) validateUpdate() (*util.ValidationContext, error) {
	ctx := util.NewValidationContext()

	if u.Id == 0 {
		return nil, errors.New("Invalid Id")
	}

	if !govalidator.IsEmail(u.Email) {
		ctx.AddError("email", "Invalid Email")
	} else {
		var count int
		if err := dbx.Get(&count, "SELECT COUNT(*) FROM user WHERE email = ? AND id != ?", u.Email, u.Id); err != nil {
			return nil, err
		}
		if count != 0 {
			ctx.AddError("email", "Email already used")
		}
	}

	if u.DisplayName == "" {
		ctx.AddError("display_name", "Display name is required")
	}

	return ctx, nil
}

func (u *User) Insert() error {
	v, err := u.validateInsert()
	if err != nil {
		return fmt.Errorf("can not validate user: %s", err)
	}

	if v.HasError() {
		return v.ToError()
	}

	now := time.Now()

	r, err := squirrel.Insert("user").
		Columns("email", "password", "display_name", "created_at").
		Values(u.Email, u.Password, u.DisplayName, now).
		RunWith(dbx.DB).Exec()

	if err != nil {
		return fmt.Errorf("can not insert user: %s", err)
	}

	if u.Id, err = r.LastInsertId(); err != nil {
		return fmt.Errorf("can not get user id: %s", err)
	}

	u.CreatedAt.Valid = true
	u.CreatedAt.Time = now

	return nil
}

func (u *User) Update() error {
	v, err := u.validateUpdate()
	if err != nil {
		return fmt.Errorf("can not validate user: %s", err)
	}
	if v.HasError() {
		return v.ToError()
	}

	now := time.Now()

	r, err := squirrel.Update("user").
		Set("email", u.Email).
		Set("password", u.Password).
		Set("display_name", u.DisplayName).
		Set("updated_at", now).
		Where("id = ?", u.Id).
		Limit(1).
		RunWith(dbx.DB).Exec()
	if err != nil {
		return fmt.Errorf("can not update user: %s", err)
	}

	var n int64
	n, err = r.RowsAffected()
	if err != nil {
		return fmt.Errorf("check update user error: %s", err)
	}
	if n != 1 {
		return fmt.Errorf("user to update not found")
	}

	u.UpdatedAt.Valid = true
	u.UpdatedAt.Time = now

	return nil
}

func (u *User) Delete() error {
	if _, err := dbx.Exec("DELETE FROM user WHERE id = ?", u.Id); err != nil {
		return fmt.Errorf("can not delete user: %s", err)
	}
	return nil
}

func FindUserByEmailPassword(email, password string) (*User, error) {
	p := NewUser()
	err := dbx.Get(p, `
		SELECT id, email, display_name, created_at, updated_at
		FROM user
		WHERE email = ? AND password = ?`, email, hashPassword(password))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func NewUser() *User {
	return &User{}
}
