package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func sha256Key(key string) []byte {
	h := sha256.New()
	h.Write([]byte(key))
	newKey := h.Sum(nil)
	return newKey
}

func pKCS7Padding(ciphertext []byte) []byte {
	bs := aes.BlockSize
	padding := bs - len(ciphertext)%bs
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key string) (string, error) {
	newKey := sha256Key(key)
	block, err := aes.NewCipher(newKey)
	if err != nil {
		return "", err
	}
	newOrigData := []byte(origData)
	newOrigData = pKCS7Padding(newOrigData)
	blockMode := cipher.NewCBCEncrypter(block, newKey[:16])
	crypted := make([]byte, len(newOrigData))
	blockMode.CryptBlocks(crypted, newOrigData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(crypted, key string) (string, error) {
	newKey := sha256Key(key)
	block, err := aes.NewCipher(newKey)
	if err != nil {
		return "", err
	}
	newCrypted, _ := base64.StdEncoding.DecodeString(crypted)
	blockMode := cipher.NewCBCDecrypter(block, newKey[:16])
	origData := make([]byte, len(newCrypted))
	blockMode.CryptBlocks(origData, newCrypted)
	origData = pKCS7UnPadding(origData)
	return string(origData), nil
}

func GenerateSaltPassword(password, salt string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}
