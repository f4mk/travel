package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

type Auth struct {
	activeKIDs []string
	method     jwt.SigningMethod
	keyFunc    func(t *jwt.Token) (any, error)
	parser     *jwt.Parser
	keyLookup  KeyLookup
}

func New(activeKIDs []string, keyLookup KeyLookup) (*Auth, error) {

	for _, activeKID := range activeKIDs {
		_, err := keyLookup.PrivateKey(activeKID)
		if err != nil {
			return nil, fmt.Errorf("cannot find active key: %w", err)
		}
	}

	method := jwt.SigningMethodRS256

	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]

		if !ok {
			return nil, fmt.Errorf("missing key ID in token header")
		}

		kidStr, ok := kid.(string)

		if !ok {
			return nil, fmt.Errorf("key ID must be of type string")
		}

		return keyLookup.PublicKey(kidStr)
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{method.Name}))

	a := Auth{
		method:     jwt.GetSigningMethod(method.Name),
		keyLookup:  keyLookup,
		keyFunc:    keyFunc,
		parser:     parser,
		activeKIDs: activeKIDs,
	}

	return &a, nil
}

func (a *Auth) GenerateToken(claims Claims) (string, error) {

	token := jwt.NewWithClaims(a.method, claims)
	// TODO: select the most fresh KID
	currentKID := a.activeKIDs[0]
	token.Header["kid"] = currentKID
	privateKey, err := a.keyLookup.PrivateKey(currentKID)
	if err != nil {
		return "", fmt.Errorf("no private key for key ID: %v", currentKID)
	}

	str, err := token.SignedString(privateKey)

	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}

func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, a.keyFunc)

	if err != nil {
		return Claims{}, fmt.Errorf("cannot parse token: %w", err)
	}

	if !token.Valid {
		return Claims{}, errors.New("token is invalid")
	}

	return claims, nil
}
