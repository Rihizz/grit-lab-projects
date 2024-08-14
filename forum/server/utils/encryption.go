package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Encrypts data
//
// Example:
// []byte("Hello World") -> []byte("�[ĝ8�p���U��(��")
func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return cipherText
}

// Decrypts data
//
// Example:
// []byte("�[ĝ8�p���U��(��") -> []byte("Hello World")
func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, _ := gcm.Open(nil, nonce, cipherText, nil)
	return plainText
}

// Helper function for Encrypt and Decrypt that creates a Hash
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
