package shop

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateKey(keySize int) ([]byte, error) {
	b := make([]byte, keySize)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (e *Enveloped) EncryptAES(plainText []byte) (cipherText []byte, cipherKey []byte, kekId int64, err error) {
	plainKey, err := GenerateKey(32)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("could not generate data key: %w", err)
	}

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return nil, nil, 0, err
	}

	// Pad plaintext to block size
	padding := aes.BlockSize - len(plainText)%aes.BlockSize
	padtext := append(plainText, bytes.Repeat([]byte{byte(padding)}, padding)...)

	cipherText = make([]byte, aes.BlockSize+len(padtext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, 0, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], padtext)

	cipherKey, kekId, err = e.kms.EncryptKey(plainKey)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("could not encrypt data key: %w", err)
	}

	return cipherText, cipherKey, kekId, nil
}

func (e *Enveloped) DecryptAES(cipherText []byte, cipherKey []byte, kekId int64) ([]byte, error) {
	plainKey, err := e.kms.DecryptKey(cipherKey, kekId)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt data encryption key: %w", err)
	}

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	padding := int(cipherText[len(cipherText)-1])
	return cipherText[:len(cipherText)-padding], nil
}

func EncryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad plaintext to block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, aes.BlockSize+len(padtext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padtext)

	return ciphertext, nil
}

func DecryptAES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padding := int(ciphertext[len(ciphertext)-1])
	return ciphertext[:len(ciphertext)-padding], nil
}
