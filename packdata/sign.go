package packdata

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)


type SignerHashSha1 struct {
	key []byte
}

func NewSignerHashSha1(key []byte) *SignerHashSha1 {
	return &SignerHashSha1{key: key}
}

func (s *SignerHashSha1) Sign(data []byte) ([]byte, error) {
	mac := hmac.New(sha1.New, s.key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

type SignerHashSha256 struct {
	key []byte
}

func NewSignerHashSha256(key []byte) *SignerHashSha256 {
	return &SignerHashSha256{key: key}
}

func (s *SignerHashSha256) Sign(data []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, s.key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

type SignerHashSha512 struct {
	key []byte
}

func NewSignerHashSha512(key []byte) *SignerHashSha512 {
	return &SignerHashSha512{key: key}
}

func (s *SignerHashSha512) Sign(data []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, s.key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}
