package list

import "errors"

var (
	ErrGetListsBusiness     = errors.New("error query lists from business layer")
	ErrListValidateListUUID = errors.New("error query list validate list uuid")
	ErrListValidateItemUUID = errors.New("error query list validate item uuid")
)
