package user

import (
	"github.com/fluxxu/util"
)

type UserProvider struct {
	*util.SqlProvider
}

func NewUserProvider() *UserProvider {
	p := new(UserProvider)
	p.SqlProvider = util.NewSqlProvider(dbx, "user", User{})
	return p
}
