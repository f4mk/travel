package user

import "github.com/f4mk/api/pkg/web"

func (nu NewUserDTO) Validate() error {
	return web.Check(nu)
}

func (uu UpdateUserDTO) Validate() error {
	return web.Check(uu)
}
