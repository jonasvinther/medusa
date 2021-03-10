package main

import "medusa/cmd"

func main() {
	cmd.Execute()

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
