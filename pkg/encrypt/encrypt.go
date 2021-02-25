package encrypt

import (
	"bufio"
	b64 "encoding/base64"
	"log"
	"os"
)

// Encrypt will incrypt the data
func Encrypt(publicKeyPath string, output string, data []byte) (encryptedKey, encryptedData string) {
	// Public key
	publicKey := ReadRsaPublicKey(publicKeyPath)

	// AEs key
	key := NewAesEncryptionKey()
	// fmt.Printf("AES key: %v \n", *key)
	text := []byte(string(data))
	encrypted, _ := AesEncrypt(text, key)

	// Base64 encode
	encryptedRsa, _ := RsaEncrypt(key, publicKey)
	b64key := b64.StdEncoding.EncodeToString([]byte(encryptedRsa))
	sEnc := b64.StdEncoding.EncodeToString(encrypted)

	return b64key, sEnc
}

func Decrypt(privateKeyPath string, output string) string {
	// Decrypt
	file, err := os.Open(output)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	decryptData, _ := b64.StdEncoding.DecodeString(txtlines[0])
	decryptKey, _ := b64.StdEncoding.DecodeString(txtlines[1])

	// Decrypt rsa
	privateKey := ReadRsaPrivateKey(privateKeyPath)
	decryptedKey, _ := RsaDecrypt(string(decryptKey), privateKey)
	decrypted, _ := AesDecrypt(decryptData, decryptedKey)

	return string(decrypted)
}
