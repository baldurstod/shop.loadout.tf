package assets

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipclient"
	"github.com/ovh/okms-sdk-go"
	"github.com/ovh/okms-sdk-go/types"
	"shop.loadout.tf/src/server/config"
)

var client *kmipclient.Client
var keyId string

func init() {
	config := initConfig()

	var err error
	// Connect to KMIP server
	client, err = kmipclient.Dial(
		config.Kms.Endpoint,
		kmipclient.WithClientCertFiles(config.Kms.CertificatePath, config.Kms.PrivateKeyPath),
	)

	if err != nil {
		panic(err)
	}

	keyId = config.Kms.KeyId
}

func initConfig() *config.Config {
	// load config
	config := config.Config{}
	if content, err := os.ReadFile("config.json"); err == nil {
		if err = json.Unmarshal(content, &config); err != nil {
			panic(err)
		}
	}

	return &config
}

func TestKmsCreateKey(t *testing.T) {
	config := initConfig()

	cert, err := tls.LoadX509KeyPair(config.Kms.CertificatePath, config.Kms.PrivateKeyPath)
	if err != nil {
		t.Error(err)
		return
	}
	httpClient := http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}},
	}

	okmsIdString := "26bb1589-ad39-4008-adad-58ef8fa8a92f"
	okmsId, err := uuid.Parse(okmsIdString)
	if err != nil {
		t.Error(err)
		return
	}

	okmsClient, err := okms.NewRestAPIClientWithHttp("https://eu-west-rbx.okms.ovh.net", &httpClient)
	if err != nil {
		t.Error(err)
		return
	}

	// Then start using the kmsClient
	// Create a new AES 256 key
	respAes, err := okmsClient.GenerateSymmetricKey(context.Background(), okmsId, types.N256, "AES key example", "SOFTWARE", "", []types.CryptographicUsages{types.Encrypt, types.Decrypt, types.WrapKey, types.UnwrapKey})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("AES KEY:", respAes.Id)
}

func TestKmipCreateKey(t *testing.T) {
	// Create an AES key
	resp, err := client.Create().
		AES(256, kmip.CryptographicUsageEncrypt|kmip.CryptographicUsageDecrypt).
		//WithName("my-encryption-key").
		Exec()

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Created AES key: %s\n", resp.UniqueIdentifier)
	// Activate the key
	_, err = client.Activate(resp.UniqueIdentifier).Exec()

	if err != nil {
		t.Error(err)
		return
	}
}

func TestKmipEncrypt(t *testing.T) {
	defer client.Close()

	nonce, _ := generateNounce(12 /*iv len for AES_GCM*/)
	plaintext := []byte("sensitive data")
	encrypted, err := client.Encrypt(keyId).
		WithCryptographicParameters(kmip.AES_GCM).
		WithIvCounterNonce(nonce).
		Data(plaintext).
		Exec()

	if err != nil {
		t.Error(err)
		return
	}

	decrypted, err := client.Decrypt(keyId).
		WithCryptographicParameters(kmip.AES_GCM).
		WithIvCounterNonce(nonce).
		WithAuthTag(encrypted.AuthenticatedEncryptionTag).
		Data(encrypted.Data).
		Exec()

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Encrypted", encrypted, string(decrypted.Data))
}

func generateNounce(keySize int) ([]byte, error) {
	b := make([]byte, keySize)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

/*
func TestDestroyKmipKeys(t *testing.T) {
	defer client.Close()

	// Find keys
	keys, err := client.Locate().Exec()
	if err != nil {
		t.Error(err)
		return
	}

	for _, keyID := range keys.UniqueIdentifier {
		fmt.Printf("Found key: %s\n", keyID)

		// Revoke the key if it is active
		client.Revoke(keyID).WithRevocationReasonCode(kmip.RevocationReasonCodeCessationOfOperation).Exec()

		// Destroy the key
		_, err = client.Destroy(keyID).Exec()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
*/

func TestListKmipKeys(t *testing.T) {
	defer client.Close()

	// Find keys
	keys, err := client.Locate().Exec()
	if err != nil {
		t.Error(err)
		return
	}

	for _, keyID := range keys.UniqueIdentifier {
		fmt.Printf("Found key: %s\n", keyID)
	}
}
