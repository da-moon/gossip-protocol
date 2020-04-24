package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"

	"github.com/palantir/stacktrace"
)

const defaultSize = 4096

// Key ...
type Key struct {
	private *rsa.PrivateKey
}

// Default ...
func Default() (*Key, error) {
	return New(defaultSize)
}

// New ...
func New(size int) (*Key, error) {
	private, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		err = stacktrace.Propagate(err, "could not generate a new RSA key with size '%d'", size)
		return nil, err
	}
	result := &Key{
		private: private,
	}
	return result, nil
}
func (i *Key) EncodedPublicKey() ([]byte, error) {
	pubKey := i.private.PublicKey
	result, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		err = stacktrace.Propagate(err, "could not marshall public key to bytes")
		return nil, err
	}
	if result == nil {
		err = stacktrace.NewError("marshalled public key was a nil slice")
		return nil, nil
	}
	if len(result) == 0 {
		err = stacktrace.NewError("marshalled public key was an empty slice")
		return nil, err
	}
	return result, nil
}
func (i *Key) PublicKeyBase64() (string, error) {
	marshalled, err := i.EncodedPublicKey()
	if err != nil {
		err = stacktrace.Propagate(err, "could not encode public key as base64")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(marshalled), nil
}

// Sha256 ...
func (i *Key) Sha256() ([]byte, error) {
	pubKey := i.private.PublicKey
	derEncoded, err := i.EncodedPublicKey()
	if err != nil {
		err = stacktrace.Propagate(err, "could not get sha256 hash of public key due to encoding issue")
		return nil, err
	}
	result := sha256.Sum256(derEncoded)
	return result[:], nil
}

// Sha256String ...
func (i *Key) Sha256String() (string, error) {
	hash, err := i.Sha256()
	if err != nil {
		err = stacktrace.Propagate(err, "could not calculate sha256 hash string of public key")
		return "", err
	}
	return hex.EncodeToString(hash), nil
}
