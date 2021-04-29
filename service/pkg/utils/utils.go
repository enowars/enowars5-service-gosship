package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"os"

	gossh "golang.org/x/crypto/ssh"
)

var privateKeyPath = "server.key"

func init() {
	if pkp, ok := os.LookupEnv("PRIVATE_KEY_PATH"); ok {
		privateKeyPath = pkp
	}
}

func generateHostSigner() (gossh.Signer, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(privateKeyPath, privateKey, 0400); err != nil {
		return nil, err
	}

	signer, err := gossh.NewSignerFromSigner(privateKey)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func readHostSigner() (gossh.Signer, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := gossh.NewSignerFromSigner(ed25519.PrivateKey(keyBytes))
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func GetHostSigner() (gossh.Signer, error) {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return generateHostSigner()
	}
	return readHostSigner()
}
