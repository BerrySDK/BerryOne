package native

import (
	"errors"

	"google.golang.org/protobuf/encoding/protowire"
)

var errInvalidHandshake = errors.New("berryone native runtime: invalid handshake message")

type HandshakeMessage struct {
	ClientHello  *HandshakeClientHello
	ServerHello  *HandshakeServerHello
	ClientFinish *HandshakeClientFinish
}

type HandshakeClientHello struct {
	Ephemeral []byte
}

type HandshakeServerHello struct {
	Ephemeral []byte
	Static    []byte
	Payload   []byte
}

type HandshakeClientFinish struct {
	Static  []byte
	Payload []byte
}

func EncodeHandshakeClientHello(ephemeral []byte) []byte {
	message := protowire.AppendTag(nil, 1, protowire.BytesType)
	message = protowire.AppendBytes(message, ephemeral)

	root := protowire.AppendTag(nil, 2, protowire.BytesType)
	root = protowire.AppendBytes(root, message)
	return root
}

func EncodeHandshakeClientFinish(staticValue, payload []byte) []byte {
	message := protowire.AppendTag(nil, 1, protowire.BytesType)
	message = protowire.AppendBytes(message, staticValue)
	message = protowire.AppendTag(message, 2, protowire.BytesType)
	message = protowire.AppendBytes(message, payload)

	root := protowire.AppendTag(nil, 4, protowire.BytesType)
	root = protowire.AppendBytes(root, message)
	return root
}

func DecodeHandshakeMessage(buffer []byte) (HandshakeMessage, error) {
	var message HandshakeMessage

	for len(buffer) > 0 {
		fieldNumber, wireType, fieldLength := protowire.ConsumeTag(buffer)
		if fieldLength < 0 {
			return HandshakeMessage{}, errInvalidHandshake
		}

		buffer = buffer[fieldLength:]

		switch fieldNumber {
		case 2:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeMessage{}, errInvalidHandshake
			}
			clientHello, err := decodeHandshakeClientHello(bytes)
			if err != nil {
				return HandshakeMessage{}, err
			}
			message.ClientHello = &clientHello
			buffer = buffer[consumed:]
		case 3:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeMessage{}, errInvalidHandshake
			}
			serverHello, err := decodeHandshakeServerHello(bytes)
			if err != nil {
				return HandshakeMessage{}, err
			}
			message.ServerHello = &serverHello
			buffer = buffer[consumed:]
		case 4:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeMessage{}, errInvalidHandshake
			}
			clientFinish, err := decodeHandshakeClientFinish(bytes)
			if err != nil {
				return HandshakeMessage{}, err
			}
			message.ClientFinish = &clientFinish
			buffer = buffer[consumed:]
		default:
			consumed := protowire.ConsumeFieldValue(fieldNumber, wireType, buffer)
			if consumed < 0 {
				return HandshakeMessage{}, errInvalidHandshake
			}
			buffer = buffer[consumed:]
		}
	}

	return message, nil
}

func decodeHandshakeClientHello(buffer []byte) (HandshakeClientHello, error) {
	var hello HandshakeClientHello

	for len(buffer) > 0 {
		fieldNumber, wireType, fieldLength := protowire.ConsumeTag(buffer)
		if fieldLength < 0 {
			return HandshakeClientHello{}, errInvalidHandshake
		}
		buffer = buffer[fieldLength:]

		switch fieldNumber {
		case 1:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeClientHello{}, errInvalidHandshake
			}
			hello.Ephemeral = append([]byte(nil), bytes...)
			buffer = buffer[consumed:]
		default:
			consumed := protowire.ConsumeFieldValue(fieldNumber, wireType, buffer)
			if consumed < 0 {
				return HandshakeClientHello{}, errInvalidHandshake
			}
			buffer = buffer[consumed:]
		}
	}

	return hello, nil
}

func decodeHandshakeServerHello(buffer []byte) (HandshakeServerHello, error) {
	var hello HandshakeServerHello

	for len(buffer) > 0 {
		fieldNumber, wireType, fieldLength := protowire.ConsumeTag(buffer)
		if fieldLength < 0 {
			return HandshakeServerHello{}, errInvalidHandshake
		}
		buffer = buffer[fieldLength:]

		switch fieldNumber {
		case 1:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeServerHello{}, errInvalidHandshake
			}
			hello.Ephemeral = append([]byte(nil), bytes...)
			buffer = buffer[consumed:]
		case 2:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeServerHello{}, errInvalidHandshake
			}
			hello.Static = append([]byte(nil), bytes...)
			buffer = buffer[consumed:]
		case 3:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeServerHello{}, errInvalidHandshake
			}
			hello.Payload = append([]byte(nil), bytes...)
			buffer = buffer[consumed:]
		default:
			consumed := protowire.ConsumeFieldValue(fieldNumber, wireType, buffer)
			if consumed < 0 {
				return HandshakeServerHello{}, errInvalidHandshake
			}
			buffer = buffer[consumed:]
		}
	}

	return hello, nil
}

func decodeHandshakeClientFinish(buffer []byte) (HandshakeClientFinish, error) {
	var finish HandshakeClientFinish

	for len(buffer) > 0 {
		fieldNumber, wireType, fieldLength := protowire.ConsumeTag(buffer)
		if fieldLength < 0 {
			return HandshakeClientFinish{}, errInvalidHandshake
		}
		buffer = buffer[fieldLength:]

		switch fieldNumber {
		case 1:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeClientFinish{}, errInvalidHandshake
			}
			finish.Static = append([]byte(nil), bytes...)
			buffer = buffer[consumed:]
		case 2:
			bytes, consumed := protowire.ConsumeBytes(buffer)
			if consumed < 0 {
				return HandshakeClientFinish{}, errInvalidHandshake
			}
			finish.Payload = append([]byte(nil), bytes...)
			buffer = buffer[consumed:]
		default:
			consumed := protowire.ConsumeFieldValue(fieldNumber, wireType, buffer)
			if consumed < 0 {
				return HandshakeClientFinish{}, errInvalidHandshake
			}
			buffer = buffer[consumed:]
		}
	}

	return finish, nil
}
