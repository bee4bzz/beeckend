package crypto

import (
	"crypto/x509"
	"encoding/pem"
)

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
