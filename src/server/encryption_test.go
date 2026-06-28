// #cgo windows CFLAGS: -I D:\Divers\libmongocrypt\include\mongocrypt

package main_test

import (
	"context"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	"shop.loadout.tf/src/server/encryption"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_, filename, _, _ := runtime.Caller(0)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	err = initConfig()
	if err != nil {
		panic(err)
	}
}

func TestEncryption(t *testing.T) {
	enveloped := encryption.NewEnveloped(encryption.Kms{})
	dekPlain, dekCipher, err := enveloped.GenerateDek(context.Background())
	if err != nil {
		t.Error(err)
	}

	dek, err := enveloped.DecryptDek(context.Background(), dekCipher)
	if err != nil {
		t.Error(err)
	}

	log.Println(dekPlain)
	log.Println(dekCipher)
	log.Println(dek)
}
