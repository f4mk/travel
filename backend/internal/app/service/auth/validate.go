package auth

import "github.com/f4mk/api/pkg/web"

func (u LoginUserDTO) Validate() error {
	return web.Check(u)
}
