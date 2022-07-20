package encoder

import (
	"encoding/base64"
	"encoding/hex"
)

type Decoder interface {
	DecodeString(string) ([]byte, error)
}

type EncodeDecoder interface {
	EncodeToString([]byte) string
	DecodeString(string) ([]byte, error)
}

type B64 struct {
	e *base64.Encoding
}

type Hex struct {
}

func NewB64() *B64 {
	return &B64{
		e: base64.StdEncoding,
	}
}

func NewHex() *Hex {
	return &Hex{}
}

func (b *B64) DecodeString(s string) ([]byte, error) {
	return b.e.DecodeString(s)
}

func (h *Hex) DecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (h *Hex) EncodeToString(data []byte) string {
	return hex.EncodeToString(data)
}
