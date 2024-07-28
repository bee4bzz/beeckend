package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"

	"fmt"
	"os"
)

// GenerateSaveKeys generates and saves keys to disk after
// encoding into PEM format.
func GenerateSaveKeys(fb string) error {
	var (
		err   error
		b     []byte
		block *pem.Block
	)

	fbPub := fb + ".pub"

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	pub := priv.Public()
	if err != nil {
		os.Exit(1)
	}

	b, err = x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}

	err = os.WriteFile(fb, pem.EncodeToMemory(block), 0600)
	if err != nil {
		return err
	}

	// public key
	pubPem, err := PublicKeyToPem(pub)
	if err != nil {
		return err
	}

	fileName := fbPub
	err = os.WriteFile(fileName, pubPem, 0600)

	return err
}

func PublicKeyToPem(pub any) ([]byte, error) {
	// public key
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	return pem.EncodeToMemory(block), nil
}

func LoadOrGenerateKeys(fb string) (any, any, error) {
	var (
		err  error
		pub  any
		priv any
	)

	fbPub := fb + ".pub"

	// private key
	priv, err = LoadPrivateKey(fb)
	if err != nil {
		// We generate the keys
		err = GenerateSaveKeys(fb)
		if err != nil {
			return nil, nil, err
		}
		priv, err = LoadPrivateKey(fb)
		if err != nil {
			return nil, nil, err
		}
	}

	// public key
	pub, err = LoadPublicKey(fbPub)
	if err != nil {
		return nil, nil, err
	}

	return pub, priv, nil
}

func LoadPrivateKey(filename string) (any, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParseRawPrivateKey(bytes)
	return key, err
}

func LoadPublicKey(filePath string) (any, error) {
	// Read PEM file
	pemData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return DecodePEMToPublicKey(pemData)
}

func DecodePEMToPublicKey(pemPubKey []byte) (any, error) {
	// Analyse Pem Block
	block, _ := pem.Decode(pemPubKey)
	if block == nil {
		return nil, fmt.Errorf("PEM decode failed")
	}

	// Verify that the type of the public key is indeed PUBLIC KEY
	switch block.Type {
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		return cert.PublicKey, nil
	case "PUBLIC KEY":
		// Analyse certificat X.509
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return pub, nil
	default:
		return nil, fmt.Errorf("wrong bloc type: %s", block.Type)
	}
}
