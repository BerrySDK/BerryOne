package native

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
)

var errShortCiphertext = errors.New("berryone native runtime: ciphertext is shorter than the GCM tag length")

func sha256Bytes(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func hkdfSHA256(secret, salt []byte, size int) []byte {
	if size <= 0 {
		return nil
	}

	prkMac := hmac.New(sha256.New, salt)
	_, _ = prkMac.Write(secret)
	prk := prkMac.Sum(nil)

	var (
		result []byte
		prev   []byte
		index  byte = 1
	)

	for len(result) < size {
		blockMac := hmac.New(sha256.New, prk)
		if len(prev) > 0 {
			_, _ = blockMac.Write(prev)
		}
		blockMac.Write([]byte{})
		_, _ = blockMac.Write([]byte{index})
		prev = blockMac.Sum(nil)
		result = append(result, prev...)
		index++
	}

	return result[:size]
}

func aesEncryptGCM(plaintext, key, iv, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aead.Seal(nil, iv, plaintext, additionalData), nil
}

func aesDecryptGCM(ciphertext, key, iv, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aead.Overhead() {
		return nil, errShortCiphertext
	}

	return aead.Open(nil, iv, ciphertext, additionalData)
}

func sharedX25519(privateKey, publicKey []byte) ([]byte, error) {
	curve := ecdh.X25519()

	priv, err := curve.NewPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	pub, err := curve.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return priv.ECDH(pub)
}
