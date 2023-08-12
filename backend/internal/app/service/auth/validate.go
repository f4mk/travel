package auth

import "github.com/f4mk/api/pkg/web"

func (u LoginUserDTO) Validate() error {
	if err := web.Check(u); err != nil {
		return err
	}
	return nil
}
