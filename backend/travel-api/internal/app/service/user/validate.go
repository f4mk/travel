package user

import "github.com/f4mk/travel/backend/travel-api/pkg/web"

func (nu NewUser) Validate() error {
	return web.Check(nu)
}

func (uu UpdateUser) Validate() error {
	return web.Check(uu)
}
