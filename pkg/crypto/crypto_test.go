package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fb    = "test"
	fbPub = "test.pub"
)

func TestLoadKeys(t *testing.T) {
	t.Run("generate and save keys succeed", func(t *testing.T) {
		// We generate the keys
		err := GenerateSaveKeys(fb)
		defer removeFiles(t, fb, fbPub)
		assert.NoError(t, err)

		// We load the private key
		privateKey, err := LoadPrivateKey(fb)
		assert.NoError(t, err)
		publicKey, err := LoadPublicKey(fbPub)
		assert.NoError(t, err)

		assertKeys(t, privateKey.(*rsa.PrivateKey), publicKey.(*rsa.PublicKey))
	})
}

func TestLoadOrGenerate(t *testing.T) {
	t.Run("generate, save and load keys succeed", func(t *testing.T) {
		// We generate the keys as the keys are not presents
		publicKey, privateKey, err := LoadOrGenerateKeys(fb)
		defer removeFiles(t, fb, fbPub)

		assert.NoError(t, err)

		assertKeys(t, privateKey.(*rsa.PrivateKey), publicKey.(*rsa.PublicKey))
	})

	// We first generate the keys
	err := GenerateSaveKeys(fb)
	defer removeFiles(t, fb, fbPub)
	assert.NoError(t, err)
	t.Run("keys already exists, they should be loaded", func(t *testing.T) {
		// The keys should be loaded
		publicKey, privateKey, err := LoadOrGenerateKeys(fb)
		assert.NoError(t, err)

		assertKeys(t, privateKey.(*rsa.PrivateKey), publicKey.(*rsa.PublicKey))
	})
}

func TestDecodeCertPEMToPublicKey(t *testing.T) {
	t.Run("decode public key from a certificate PEM", func(t *testing.T) {
		// Create a certificate
		caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
		assert.NoError(t, err)
		caBytes, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{SerialNumber: big.NewInt(1)}, &x509.Certificate{SerialNumber: big.NewInt(2)}, &caPrivKey.PublicKey, caPrivKey)
		assert.NoError(t, err)
		certPEM := bytes.NewBufferString("")
		err = pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
		assert.NoError(t, err)

		// Sign a message
		signature, err := SignMessageSHA256(caPrivKey, "test")
		assert.NoError(t, err)

		// Get the public key from the certificate
		pb, err := DecodePEMToPublicKey(certPEM.Bytes())
		assert.NoError(t, err)

		// Assert
		assert.True(t, VerifySignatureSHA256(pb.(*rsa.PublicKey), "test", signature))
	})
}

func removeFiles(t *testing.T, fb, fbPub string) {
	t.Helper()
	err := os.Remove(fb)

	if err != nil {
		t.Fatal(err)
	}
	os.Remove(fbPub)

	if err != nil {
		t.Fatal(err)
	}
}

func assertKeys(t *testing.T, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	t.Helper()
	// The message to sign
	message := "test"

	// We sign the message
	signedString, err := SignMessageSHA512(privateKey, message)

	assert.NoError(t, err)

	// We verify the signature
	res := VerifySignatureSHA512(publicKey, message, signedString)

	assert.True(t, res)
}

// This is a helper function, see package jwt to verify signature.
func VerifySignatureSHA512(publicKey *rsa.PublicKey, message string, signature []byte) bool {
	// message hash.
	hashed := sha512.Sum512([]byte(message))

	// signature verification.
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashed[:], signature)
	return err == nil
}

// This is a helper function, see package jwt to verify signature.
func VerifySignatureSHA256(publicKey *rsa.PublicKey, message string, signature []byte) bool {
	// message hash
	hashed := sha256.Sum256([]byte(message))

	// signature verification.
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	return err == nil
}

// This is a helper function, see package jwt to sign a message.
func SignMessageSHA512(privateKey *rsa.PrivateKey, message string) ([]byte, error) {
	// message hash
	hashed := sha512.Sum512([]byte(message))

	// sign the message
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashed[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// This is a helper function, see package jwt to sign a message.
func SignMessageSHA256(privateKey *rsa.PrivateKey, message string) ([]byte, error) {
	// message hash
	hashed := sha256.Sum256([]byte(message))

	// sign the message
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}
