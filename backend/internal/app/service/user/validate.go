package user

import "github.com/f4mk/api/pkg/web"

func (nu NewUserDTO) Validate() error {
	if err := web.Check(nu); err != nil {
		return err
	}
	return nil
}

func (uu UpdateUserDTO) Validate() error {
	if err := web.Check(uu); err != nil {
		return err
	}
	return nil
}
