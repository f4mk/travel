package keystore

import (
	"crypto/rsa"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

type KeyStore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

func New() *KeyStore {

	return &KeyStore{
		mu:    sync.RWMutex{},
		store: make(map[string]*rsa.PrivateKey),
	}
}

// NewMap constructs a KeyStore with an initial set of keys.
func NewMap(store map[string]*rsa.PrivateKey) *KeyStore {
	return &KeyStore{
		store: store,
	}
}

// NewFS constructs a KeyStore based on a set of PEM files rooted inside
// of a directory. The name of each PEM file will be used as the key id.
// Example: keystore.NewFS(os.DirFS("/secret/jwt/"))
// Example: /secret/jwt/77a6ddf0-c968-4800-829e-27a26e3b3cbd.pem
func NewFS(fsys fs.FS) (*KeyStore, error) {
	// TODO: need logging here?
	ks := KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}

	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return ErrWalkDir
		}

		if dirEntry.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return ErrOpenKeyFile
		}
		defer file.Close()

		// limit PEM file size to 1 megabyte. This should be reasonable for
		// almost any PEM file and prevents shenanigans like linking the file
		// to /dev/random or something like that.
		privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return ErrReadPrivateKeyFile
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
		if err != nil {
			return ErrParsePrivateKey
		}

		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = privateKey

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return &ks, nil
}

func (ks *KeyStore) Add(keyID string, key *rsa.PrivateKey) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.store[keyID] = key
}

func (ks *KeyStore) Remove(keyID string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	delete(ks.store, keyID)
}

// PrivateKey searches the key store for a given kid and returns the private key.
func (ks *KeyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, ErrPrivateKeyLookup
	}

	return privateKey, nil
}

func (ks *KeyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	privateKey, err := ks.PrivateKey(kid)
	if err != nil {
		return nil, err
	}

	return &privateKey.PublicKey, nil
}

func (ks *KeyStore) CollectKeyIDs() []string {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	var ids []string
	for k := range ks.store {
		ids = append(ids, k)
	}
	return ids
}
