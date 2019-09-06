package packdata

import (
	"fmt"
	"testing"
)

func TestSignerHashSha1_Sign(t *testing.T) {
	key := []byte("token")
	sh1Signer := NewSignerHashSha1(key)
	out, err := sh1Signer.Sign([]byte("hello"))
	if err != nil {
		t.Error(err)
		return
	}
	ft := fmt.Sprintf("%x", out)
	if ft != "8e4f581f0b7be2dc6db4ecfb5309361c42f9714a" {
		t.Fail()
	}
}

func TestNewSignerHashSha256_Sign(t *testing.T) {
	key := []byte("token")
	sh1Signer := NewSignerHashSha256(key)
	out, err := sh1Signer.Sign([]byte("hello"))
	if err != nil {
		t.Error(err)
		return
	}
	ft := fmt.Sprintf("%x", out)
	if ft != "df3178e409a68446314d5be83b911b78dc6fa272a556429fd9d2092575dcf174" {
		t.Fail()
	}
}

func TestNewSignerHashSha512_Sign(t *testing.T) {
	key := []byte("token")
	sh1Signer := NewSignerHashSha512(key)
	out, err := sh1Signer.Sign([]byte("hello"))
	if err != nil {
		t.Error(err)
		return
	}
	ft := fmt.Sprintf("%x", out)
	if ft != "a7b337f71e739bd6fa63320a1b1220c694317b9ae5f7ed45d521b414ebed434621f78815cd6a76f06f55ce9527e41a08b252d7bbb584d3d60401e9de5227648f" {
		t.Fail()
	}
}