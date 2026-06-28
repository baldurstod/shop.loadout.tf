package encryption

import "context"

type KeyEncryptionWrapper interface {
	EncryptKey(ctx context.Context, keyPlain []byte) ([]byte, error)
	DecryptKey(ctx context.Context, keyCipher []byte) ([]byte, error)
}

type Enveloped struct {
	kms KeyEncryptionWrapper
}

func NewEnveloped(kms KeyEncryptionWrapper) *Enveloped {
	return &Enveloped{
		kms: kms,
	}
}

var kek = []byte("32-byte-key-for-AES-256-!!!!!!!!")

type Kms struct {
}

func (k Kms) EncryptKey(ctx context.Context, keyPlain []byte) ([]byte, error) {
	return EncryptAES(keyPlain, kek)
}

func (k Kms) DecryptKey(ctx context.Context, keyCipher []byte) ([]byte, error) {
	return DecryptAES(keyCipher, kek)
}
