package image

import "errors"

var (
	ErrGetImageBusiness     = errors.New("error get image from business layer")
	ErrPostImageBusiness    = errors.New("error post image from business layer")
	ErrPostImageDecode      = errors.New("error post image parsing user input")
	ErrPostImageDecodeLen   = errors.New("error post image parsing user input: no images found")
	ErrPostImageRead        = errors.New("error post image parsing user input: cannot open image")
	ErrPostImageReadContent = errors.New("error post image parsing user input: cannot read image content")
)
