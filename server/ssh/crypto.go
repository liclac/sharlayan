package ssh

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func LoadOrGenerateHostKey(path string) (ssh.Signer, error) {
	key, err := LoadHostKey(path)
	if os.IsNotExist(err) {
		if err := GenerateHostKey(path); err != nil {
			return nil, err
		}
		return LoadHostKey(path)
	}
	return key, err
}

func LoadHostKey(path string) (ssh.Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(data)
}

func GenerateHostKey(path string) error {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	privDER, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privDER,
	}), 0600)
}
