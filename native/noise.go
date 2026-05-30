package native

import (
	"encoding/binary"
	"errors"
)

var (
	errInvalidNoiseFrame          = errors.New("berryone native runtime: invalid noise frame")
	ErrRegistrationPayloadPending = errors.New("berryone native runtime: server hello received, but registration payload is not ported yet")
)

type NoiseHandler struct {
	privateKey  []byte
	publicKey   []byte
	hash        []byte
	salt        []byte
	encKey      []byte
	decKey      []byte
	counter     uint32
	sentIntro   bool
	introHeader []byte
}

func NewNoiseHandler(privateKey, publicKey, routingInfo []byte) *NoiseHandler {
	data := append([]byte(nil), NoiseMode...)
	hash := data
	if len(data) != 32 {
		hash = sha256Bytes(data)
	}

	introHeader := append([]byte(nil), NoiseWAHeader...)
	if len(routingInfo) > 0 {
		header := make([]byte, 7+len(routingInfo)+len(NoiseWAHeader))
		copy(header[0:], []byte("ED"))
		header[2] = 0
		header[3] = 1
		header[4] = byte(len(routingInfo) >> 16)
		binary.BigEndian.PutUint16(header[5:7], uint16(len(routingInfo)&0xffff))
		copy(header[7:], routingInfo)
		copy(header[7+len(routingInfo):], NoiseWAHeader)
		introHeader = header
	}

	handler := &NoiseHandler{
		privateKey:  append([]byte(nil), privateKey...),
		publicKey:   append([]byte(nil), publicKey...),
		hash:        append([]byte(nil), hash...),
		salt:        append([]byte(nil), hash...),
		encKey:      append([]byte(nil), hash...),
		decKey:      append([]byte(nil), hash...),
		introHeader: introHeader,
	}

	handler.authenticate(NoiseWAHeader)
	handler.authenticate(publicKey)
	return handler
}

func (n *NoiseHandler) EncodeFrame(data []byte) []byte {
	introSize := 0
	if !n.sentIntro {
		introSize = len(n.introHeader)
	}

	frame := make([]byte, introSize+3+len(data))
	if !n.sentIntro {
		copy(frame, n.introHeader)
		n.sentIntro = true
	}

	offset := introSize
	frame[offset] = byte(len(data) >> 16)
	frame[offset+1] = byte(len(data) >> 8)
	frame[offset+2] = byte(len(data))
	copy(frame[offset+3:], data)
	return frame
}

func (n *NoiseHandler) DecodeHandshakeFrame(frame []byte) ([]byte, error) {
	if len(frame) < 3 {
		return nil, errInvalidNoiseFrame
	}

	size := int(frame[0])<<16 | int(frame[1])<<8 | int(frame[2])
	if size <= 0 || len(frame) < size+3 {
		return nil, errInvalidNoiseFrame
	}

	return append([]byte(nil), frame[3:3+size]...), nil
}

func (n *NoiseHandler) ProcessServerHello(serverHello HandshakeServerHello, noiseKey X25519KeyPair) ([]byte, error) {
	n.authenticate(serverHello.Ephemeral)

	shared, err := sharedX25519(n.privateKey, serverHello.Ephemeral)
	if err != nil {
		return nil, err
	}
	n.mixIntoKey(shared)

	decStatic, err := n.decrypt(serverHello.Static)
	if err != nil {
		return nil, err
	}

	sharedStatic, err := sharedX25519(n.privateKey, decStatic)
	if err != nil {
		return nil, err
	}
	n.mixIntoKey(sharedStatic)

	if _, err := n.decrypt(serverHello.Payload); err != nil {
		return nil, err
	}

	keyEnc, err := n.encrypt(noiseKey.Public)
	if err != nil {
		return nil, err
	}

	sharedNoise, err := sharedX25519(noiseKey.Private, serverHello.Ephemeral)
	if err != nil {
		return nil, err
	}
	n.mixIntoKey(sharedNoise)

	return keyEnc, nil
}

func (n *NoiseHandler) authenticate(data []byte) {
	hashInput := make([]byte, 0, len(n.hash)+len(data))
	hashInput = append(hashInput, n.hash...)
	hashInput = append(hashInput, data...)
	n.hash = sha256Bytes(hashInput)
}

func (n *NoiseHandler) encrypt(plaintext []byte) ([]byte, error) {
	iv := generateNoiseIV(n.counter)
	n.counter++

	result, err := aesEncryptGCM(plaintext, n.encKey, iv, n.hash)
	if err != nil {
		return nil, err
	}

	n.authenticate(result)
	return result, nil
}

func (n *NoiseHandler) decrypt(ciphertext []byte) ([]byte, error) {
	iv := generateNoiseIV(n.counter)
	n.counter++

	result, err := aesDecryptGCM(ciphertext, n.decKey, iv, n.hash)
	if err != nil {
		return nil, err
	}

	n.authenticate(ciphertext)
	return result, nil
}

func (n *NoiseHandler) mixIntoKey(data []byte) {
	key := hkdfSHA256(data, n.salt, 64)
	n.salt = append([]byte(nil), key[:32]...)
	n.encKey = append([]byte(nil), key[32:]...)
	n.decKey = append([]byte(nil), key[32:]...)
	n.counter = 0
}

func generateNoiseIV(counter uint32) []byte {
	iv := make([]byte, 12)
	binary.BigEndian.PutUint32(iv[8:], counter)
	return iv
}
