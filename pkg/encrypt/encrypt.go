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
	text := []byte(string(data))
	encrypted, _ := AesEncrypt(text, key)

	// Base64 encode
	encryptedRsa, _ := RsaEncrypt(key, publicKey)
	b64key := b64.StdEncoding.EncodeToString([]byte(encryptedRsa))
	sEnc := b64.StdEncoding.EncodeToString(encrypted)

	return b64key, sEnc
}

func Decrypt(privateKeyPath string, output string) (string, error) {
	file, err := os.Open(output)
	if err != nil {
		return "", err
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	// We don't know the length of the longesst line in the file yet.
	// Let's just find the total size of the file and set that as
	// the maximum size of the buffer.
	var maxCapacity int = int(fi.Size())
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	decryptData, err := b64.StdEncoding.DecodeString(txtlines[0])
	if err != nil {
		return "", err
	}
	decryptKey, err := b64.StdEncoding.DecodeString(txtlines[1])
	if err != nil {
		return "", err
	}

	// Decrypt rsa
	privateKey := ReadRsaPrivateKey(privateKeyPath)
	decryptedKey, err := RsaDecrypt(string(decryptKey), privateKey)
	if err != nil {
		return "", err
	}
	decrypted, err := AesDecrypt(decryptData, decryptedKey)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
