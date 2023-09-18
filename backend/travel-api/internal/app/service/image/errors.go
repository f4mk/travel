package image

import "errors"

var (
	ErrGetImageBusiness  = errors.New("error get image from business layer")
	ErrPostImageBusiness = errors.New("error post image from business layer")
)
