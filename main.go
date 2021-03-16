package main

import (
	"fmt"
	"medusa/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// key := encrypt.NewAesEncryptionKey()
	// text := []byte("Hello from Medusa")
	// encrypted, err := encrypt.AesEncrypt(text, key)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println("Text: " + string(text))
	// fmt.Println("Encrypted: " + string(encrypted))

	// privKey := encrypt.ReadRsaPrivateKey("./private-key.pem")
	// pubKey := encrypt.ReadRsaPublicKey("./public-key.pem")

	// encrptedRsa, _ := encrypt.RsaEncrypt(key, pubKey)
	// decryptedKey, _ := encrypt.RsaDecrypt(encrptedRsa, privKey)

	// decrypted, err := encrypt.AesDecrypt(encrypted, decryptedKey)
	// fmt.Println("Decrypted: " + string(decrypted))

	// fmt.Println("encrypt: " + encrptedRsa)
}
