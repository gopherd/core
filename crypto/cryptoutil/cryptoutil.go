package cryptoutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var (
	ErrNonRSAPublicKey          = errors.New("non-rsa public key")
	ErrInvalidPublicKeyPemBlock = errors.New("invalid public key pem block")
)

func GenerateRSAPemFile(priKey *rsa.PrivateKey, priFilename, pubFilename string) error {
	priFile, err := os.OpenFile(priFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer priFile.Close()
	pubFile, err := os.OpenFile(pubFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer pubFile.Close()

	if err := GenerateRSAPrivatePem(priKey, priFile); err != nil {
		return err
	}
	if err := GenerateRSAPublicPem(&priKey.PublicKey, pubFile); err != nil {
		return err
	}
	return nil
}

func GenerateRSAPrivatePem(priKey *rsa.PrivateKey, w io.Writer) error {
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priKey),
	}
	return pem.Encode(w, block)
}

func GenerateRSAPublicPem(pubKey *rsa.PublicKey, w io.Writer) error {
	block := &pem.Block{
		Type: "RSA PUBLIC KEY",
	}
	var err error
	block.Bytes, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return err
	}
	return pem.Encode(w, block)
}

func LoadRSAPrivateKeyFile(priFilename string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(priFilename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadRSAPublicKeyFile(priFilename string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(priFilename)
	if err != nil {
		return nil, err
	}
	return LoadRSAPublicKeyFromBytes(data)
}

func LoadRSAPublicKeyFromBytes(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidPublicKeyPemBlock
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNonRSAPublicKey
	}
	return rsaPubKey, nil
}
