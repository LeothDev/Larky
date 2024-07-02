package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/zenazn/pkcs7pad"
	"io"
)

func Encrypt(encryptKey, eventContent string) string {
	// Convert encryptKey to SHA256 hashed key
	key := sha256.Sum256([]byte(encryptKey))

	// Convert eventContent to []byte
	plaintext := []byte(eventContent)

	// Pad the plaintext using PKCS7
	paddedPlainText := pkcs7pad.Pad(plaintext, aes.BlockSize)

	// Generate a random 16-byte IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Create AES cipher using the hashed key
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	// Encrypt the padded plaintext using CBC mode
	ciphertext := make([]byte, len(paddedPlainText))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, paddedPlainText)

	// Combine iv and encrypted ciphertext
	combined := append(iv, ciphertext...)

	// Encode combined data in base64
	encryptedBase64 := base64.StdEncoding.EncodeToString(combined)

	return encryptedBase64
}

func Decrypt(encryptKey, eventContent string) string {
	// Convert encryptedKey to a SHA256 hashed key of 32 bytes
	key := sha256.Sum256([]byte(encryptKey))

	// Decode the base64 encrypted data
	combined, err := base64.StdEncoding.DecodeString(eventContent)
	if err != nil {
		panic("base64 decoding error: " + err.Error())
	}

	// Split the combined data into iv and ciphertext
	iv := combined[:aes.BlockSize]
	ciphertext := combined[aes.BlockSize:]

	// Create AES cipher block using the hashed key
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic("AES cipher creation error: " + err.Error())
	}

	// Decrypt the ciphertext using CBC mode
	plaintext := make([]byte, len(ciphertext))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	plaintext, err = pkcs7pad.Unpad(plaintext)
	if err != nil {
		fmt.Printf("Pkcs7pad Error: %v", err)
	}
	return string(plaintext)
}
