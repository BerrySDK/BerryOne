package native

import "encoding/hex"

var (
	NoiseMode     = []byte("Noise_XX_25519_AESGCM_SHA256\x00\x00\x00\x00")
	NoiseWAHeader = []byte{87, 65, 6, 3}
)

type WACertDetails struct {
	Serial    uint32
	Issuer    string
	PublicKey []byte
}

var DefaultWACertDetails = WACertDetails{
	Serial:    0,
	Issuer:    "WhatsAppLongTerm1",
	PublicKey: mustDecodeHex("142375574d0a587166aae71ebe516437c4a28b73e3695c6ce1f7f9545da8ee6b"),
}

func mustDecodeHex(value string) []byte {
	bytes, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return bytes
}
