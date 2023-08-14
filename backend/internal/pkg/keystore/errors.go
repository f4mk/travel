package keystore

import "errors"

var (
	ErrWalkDir            = errors.New("error walking dir")
	ErrOpenKeyFile        = errors.New("error opening key file")
	ErrReadPrivateKeyFile = errors.New("error reading auth private key file")
	ErrParsePrivateKey    = errors.New("error parsing auth private key")
	ErrPrivateKeyLookup   = errors.New("error private key lookup")
)
