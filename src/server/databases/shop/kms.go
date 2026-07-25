package shop

import (
	"sync"

	"shop.loadout.tf/src/server/kmip"
)

var maxKeyUse = 1000

type KeyEncryptionWrapper interface {
	EncryptKey(keyPlain []byte) ([]byte, int64, error)
	DecryptKey(keyCipher []byte, kekId int64) ([]byte, error)
}

type Enveloped struct {
	kms KeyEncryptionWrapper
}

func newEnveloped(kms KeyEncryptionWrapper) *Enveloped {
	return &Enveloped{
		kms: kms,
	}
}

type kms struct{}

func (k kms) EncryptKey(keyPlain []byte) (keyCipher []byte, kekId int64, err error) {
	key, kekId, err := getEncryptKey()
	if err != nil {
		return nil, 0, err
	}
	keyCipher, err = EncryptAES(keyPlain, key)
	if err != nil {
		return nil, 0, err
	}
	return
}

func (k kms) DecryptKey(dekCipher []byte, kekId int64) ([]byte, error) {
	kekCipher, nonce, encryptionTag, err := getKek(kekId)
	if err != nil {
		return nil, err
	}

	kekPlain, err := decryptKek(kekCipher, nonce, encryptionTag)
	if err != nil {
		return nil, err
	}

	return DecryptAES(dekCipher, kekPlain)
}

var getEncryptKey = func() func() (kekPlain []byte, kekId int64, err error) {
	var currentKey []byte
	var currentKekId int64
	var use int
	var mutex sync.Mutex
	return func() (kekPlain []byte, kekId int64, err error) {
		mutex.Lock()
		defer mutex.Unlock()
		if currentKey == nil || use > maxKeyUse {
			kekPlain, err = GenerateKey(32)
			if err != nil {
				return nil, 0, err
			}

			var kekCipher []byte
			var nonce []byte
			var encryptionTag []byte
			kekCipher, nonce, encryptionTag, err = kmip.KmipEncryptKek(kekPlain)
			if err != nil {
				return nil, 0, err
			}

			//var kekId int64
			kekId, err = insertKek(kekCipher, nonce, encryptionTag)
			if err != nil {
				return nil, 0, err
			}

			currentKey = kekPlain
			currentKekId = kekId
			use = 0

			return
		}

		use++
		return currentKey, currentKekId, nil
	}
}()

var keks = struct {
	sync.RWMutex
	m map[string][]byte
}{m: map[string][]byte{}}

func decryptKek(kekCipher []byte, nonce []byte, encryptionTag []byte) (kekPlain []byte, err error) {
	var read = func() ([]byte, bool) {
		keks.RLock()
		defer keks.RUnlock()
		kekPlain, found := keks.m[string(kekCipher)]
		return kekPlain, found
	}

	var write = func(kekCipher []byte, kekPlain []byte) {
		keks.Lock()
		defer keks.Unlock()
		keks.m[string(kekCipher)] = kekPlain
	}

	kekPlain, found := read()

	// If we have the key, return
	if found {
		return
	}

	// Decrypt the key encryption key
	kekPlain, err = kmip.KmipDecryptKey(kekCipher, nonce, encryptionTag)
	if err != nil {
		return nil, err
	}

	// cache the key
	write(kekCipher, kekPlain)

	return kekPlain, nil
}
