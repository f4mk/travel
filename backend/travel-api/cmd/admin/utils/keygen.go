package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/f4mk/travel/backend/travel-api/config"
)

func GenerateKey(cfg *config.Config) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Errorf("cannot generate private key: %w", err))
	}

	filename := uuid.New().String()
	privateFile, err := os.Create(filepath.Join(cfg.Auth.KeyPath, filename+".pem"))
	if err != nil {
		panic(fmt.Errorf("creating private file: %w", err))
	}
	defer privateFile.Close()

	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		panic(fmt.Errorf("encoding to private file: %w", err))
	}

	// Create a file for the public key information in PEM form.
	publicFile, err := os.Create(filepath.Join(cfg.Auth.KeyPath, filename+".pub"))
	if err != nil {
		panic(fmt.Errorf("creating public file: %w", err))
	}
	defer publicFile.Close()

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(fmt.Errorf("marshaling public key: %w", err))
	}

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the public key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		panic(fmt.Errorf("encoding to public file: %w", err))
	}

	fmt.Println("private and public key files generated")
}
