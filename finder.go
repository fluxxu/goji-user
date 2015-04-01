package user

import (
	"github.com/fluxxu/goji-auth"
)

// implement auth.UserFinderInterface
type Finder struct{}

func (f *Finder) FindUserByEmailAndPassword(email, password string) (auth.UserInterface, error) {
	p, err := FindUserByEmailPassword(email, password)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, auth.ErrUserNotFound
	}
	return p, nil
}
