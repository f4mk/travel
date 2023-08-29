package utils

import (
	"crypto/rsa"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/f4mk/travel/backend/travel-api/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAllTokens(cfg *config.Config, roles []string) ([]map[string]string, error) {

	tokensSlice := make([]map[string]string, len(roles))
	keys := make(map[string]*rsa.PrivateKey)
	fsys := os.DirFS(cfg.Auth.KeyPath)

	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		defer file.Close()

		// limit PEM file size to 1 megabyte. This should be reasonable for
		// almost any PEM file and prevents shenanigans like linking the file
		// to /dev/random or something like that.
		privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth private key: %w", err)
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("parsing auth private key: %w", err)
		}

		keys[strings.TrimSuffix(dirEntry.Name(), ".pem")] = privateKey

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)

	for i, role := range roles {

		tokensSlice[i] = make(map[string]string)

		claims := struct {
			jwt.RegisteredClaims
			Roles        []string
			TokenVersion int32
		}{
			RegisteredClaims: jwt.RegisteredClaims{
				// TODO: change this to a real subject
				Subject:   "12345",
				Issuer:    "travel service",
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(365 * time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles:        []string{role},
			TokenVersion: 0,
		}

		for key, privateKey := range keys {
			method := jwt.GetSigningMethod("RS256")
			token := jwt.NewWithClaims(method, claims)
			token.Header["kid"] = key
			str, err := token.SignedString(privateKey)
			if err != nil {
				return nil, fmt.Errorf("cannot sign token: %w", err)
			}
			tokensSlice[i][key] = str

		}
	}

	return tokensSlice, nil
}

func GenerateToken(cfg *config.Config, kid string, roles []string) (map[string]string, error) {

	tokens := make(map[string]string)

	file, err := os.Open(filepath.Join(cfg.Auth.KeyPath, kid) + ".pem")
	if err != nil {
		return nil, fmt.Errorf("opening key file: %w", err)
	}
	defer file.Close()

	// limit PEM file size to 1 megabyte. This should be reasonable for
	// almost any PEM file and prevents shenanigans like linking the file
	// to /dev/random or something like that.
	privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("reading auth private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, fmt.Errorf("parsing auth private key: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)

	for _, role := range roles {

		claims := struct {
			jwt.RegisteredClaims
			Roles        []string
			TokenVersion int32
		}{
			RegisteredClaims: jwt.RegisteredClaims{
				// TODO: change this to a real subject
				Subject:   "12345",
				Issuer:    "travel service",
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(365 * time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles:        []string{role},
			TokenVersion: 0,
		}

		method := jwt.GetSigningMethod("RS256")
		token := jwt.NewWithClaims(method, claims)
		token.Header["kid"] = kid
		str, err := token.SignedString(privateKey)
		if err != nil {
			return nil, fmt.Errorf("cannot sign token: %w", err)
		}
		tokens[role] = str
	}

	return tokens, nil
}
