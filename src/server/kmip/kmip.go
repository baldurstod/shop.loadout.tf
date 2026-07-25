package kmip

import (
	"crypto/rand"
	"errors"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipclient"
	"shop.loadout.tf/src/server/config"
)

var client *kmipclient.Client
var keyId string

func InitKmip(kmsConfig config.Kms) error {
	var err error
	// Connect to KMIP server
	client, err = kmipclient.Dial(
		kmsConfig.Endpoint,
		kmipclient.WithClientCertFiles(kmsConfig.CertificatePath, kmsConfig.PrivateKeyPath),
	)

	keyId = kmsConfig.KeyId

	if err != nil {
		return err
	}
	return nil
}

func KmipEncryptKek(plainKey []byte) (cipherKey []byte, nonce []byte, encryptionTag []byte, err error) {
	if client == nil {
		return nil, nil, nil, errors.New("kmip client is null. Did you forgot to init kmip ?")
	}

	nonce, err = generateNounce(12 /*iv len for AES_GCM*/)
	encrypted, err := client.Encrypt(keyId).
		WithCryptographicParameters(kmip.AES_GCM).
		WithIvCounterNonce(nonce).
		Data(plainKey).
		Exec()

	if err != nil {
		return nil, nil, nil, err
	}

	return encrypted.Data, nonce, encrypted.AuthenticatedEncryptionTag, nil
}

func generateNounce(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func KmipDecryptKey(cipherKey []byte, nonce []byte, encryptionTag []byte) (plainKey []byte, err error) {
	if client == nil {
		return nil, errors.New("kmip client is null. Did you forgot to init kmip ?")
	}

	decrypted, err := client.Decrypt(keyId).
		WithCryptographicParameters(kmip.AES_GCM).
		WithIvCounterNonce(nonce).
		WithAuthTag(encryptionTag).
		Data(cipherKey).
		Exec()

	if err != nil {
		return nil, err
	}

	return decrypted.Data, nil
}
