package cobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNativeBasicFeatures(t *testing.T) {
	config := Config{
		Type:        Native,
		Reverse:     false,
		SpecialByte: 0x00,
		Delimiter:   true,
		EndingSave:  false,
	}
	encoder, status := NewEncoder(config)
	assert.Equal(t, status, nil)

	required_message := "aaaaaaaaaaa"
	required_raw := []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_encode := []byte{12, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 0}
	required_decode := required_raw

	raw := []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode := encoder.Encode(raw)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, encoder.Verify(encode), nil)
	decode := encoder.Decode(encode)
	assert.Equal(t, decode, required_decode)
}