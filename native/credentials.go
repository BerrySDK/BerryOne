package native

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"time"
)

type X25519KeyPair struct {
	Public  []byte `json:"public"`
	Private []byte `json:"private"`
}

type AuthCredentials struct {
	ClientID       string        `json:"client_id"`
	RegistrationID uint32        `json:"registration_id"`
	NoiseKey       X25519KeyPair `json:"noise_key"`
	IdentityKey    X25519KeyPair `json:"identity_key"`
	ADVSecret      []byte        `json:"adv_secret"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func GenerateAuthCredentials() (AuthCredentials, error) {
	noiseKey, err := generateX25519KeyPair()
	if err != nil {
		return AuthCredentials{}, err
	}

	identityKey, err := generateX25519KeyPair()
	if err != nil {
		return AuthCredentials{}, err
	}

	advSecret := make([]byte, 32)
	if _, err := rand.Read(advSecret); err != nil {
		return AuthCredentials{}, err
	}

	now := time.Now().UTC()
	return AuthCredentials{
		ClientID:       randomHex(16),
		RegistrationID: randomUint32(),
		NoiseKey:       noiseKey,
		IdentityKey:    identityKey,
		ADVSecret:      advSecret,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (c AuthCredentials) NoisePublicB64() string {
	return base64.StdEncoding.EncodeToString(c.NoiseKey.Public)
}

func (c AuthCredentials) IdentityPublicB64() string {
	return base64.StdEncoding.EncodeToString(c.IdentityKey.Public)
}

func (c AuthCredentials) ADVSecretB64() string {
	return base64.StdEncoding.EncodeToString(c.ADVSecret)
}

func generateX25519KeyPair() (X25519KeyPair, error) {
	curve := ecdh.X25519()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return X25519KeyPair{}, err
	}

	return X25519KeyPair{
		Public:  privateKey.PublicKey().Bytes(),
		Private: privateKey.Bytes(),
	}, nil
}

func GenerateX25519KeyPair() (X25519KeyPair, error) {
	return generateX25519KeyPair()
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405")
	}
	return hex.EncodeToString(bytes)
}

func randomUint32() uint32 {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return uint32(time.Now().Unix())
	}
	return uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
}
