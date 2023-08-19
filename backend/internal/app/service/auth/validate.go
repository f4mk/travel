package auth

import "github.com/f4mk/api/pkg/web"

func (u LoginUser) Validate() error {
	return web.Check(u)
}

func (cp ChangePassword) Validate() error {
	return web.Check(cp)
}

func (cp ResetPassword) Validate() error {
	return web.Check(cp)
}

func (cp SubmitResetPassword) Validate() error {
	return web.Check(cp)
}
