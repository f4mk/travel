package list

import "errors"

var (
	ErrGetListsBusiness     = errors.New("error query lists from business layer")
	ErrCreateListBusiness   = errors.New("error create list from business layer")
	ErrCreateItemBusiness   = errors.New("error create item from business layer")
	ErrUpdateListBusiness   = errors.New("error update list from business layer")
	ErrUpdateItemBusiness   = errors.New("error update item from business layer")
	ErrDeleteListBusiness   = errors.New("error delete list from business layer")
	ErrDeleteItemBusiness   = errors.New("error delete item from business layer")
	ErrListValidateListUUID = errors.New("error list validate list uuid")
	ErrListValidateItemUUID = errors.New("error list validate item uuid")
	ErrItemValidateListUUID = errors.New("error item validate list uuid")
	ErrItemValidateItemUUID = errors.New("error item validate item uuid")
	ErrListCreateValidate   = errors.New("error create list parsing user input")
	ErrListUpdateValidate   = errors.New("error update list parsing user input")
	ErrListDeleteValidate   = errors.New("error delete list parsing user input")
	ErrItemCreateValidate   = errors.New("error create item parsing user input")
	ErrItemUpdateValidate   = errors.New("error update item parsing user input")
	ErrItemDeleteValidate   = errors.New("error delete item parsing user input")
)
