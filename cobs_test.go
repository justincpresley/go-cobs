package cobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNativeBasicFeatures(t *testing.T) {
	// [Native] Delimiter | 0x00
	config := Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		Type:        Native,
	}
	required_message := "aaaaaaaaaaa"
	required_raw := []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_encode := []byte{12, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 0}
	required_decode := required_raw

	raw := []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode := Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode := Decode(encode, config)
	assert.Equal(t, decode, required_decode)

	// [Native] No Delimiter | 0x00
	config = Config{
		SpecialByte: 0x00,
		Delimiter:   false,
		Type:        Native,
	}
	required_message = "aaaaaaaaaaa"
	required_raw = []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_encode = []byte{12, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_decode = required_raw

	raw = []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode = Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode = Decode(encode, config)
	assert.Equal(t, decode, required_decode)

	// [Native] No Delimiter | 0x61
	config = Config{
		SpecialByte: 0x61,
		Delimiter:   false,
		Type:        Native,
	}
	required_message = "aabbbaabbabbabb"
	required_raw = []byte{97, 97, 98, 98, 98, 97, 97, 98, 98, 97, 98, 98, 97, 98, 98}
	required_encode = []byte{1, 1, 4, 98, 98, 98, 1, 3, 98, 98, 3, 98, 98, 3, 98, 98}
	required_decode = required_raw

	raw = []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode = Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode = Decode(encode, config)
	assert.Equal(t, decode, required_decode)
}

func TestReducedBasicFeatures(t *testing.T) {
	// [Reduced] Delimiter | 0x00
	config := Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		Type:        Reduced,
	}
	required_message := "aaaaaaaaaaa"
	required_raw := []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_encode := []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 0}
	required_decode := required_raw

	raw := []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode := Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode := Decode(encode, config)
	assert.Equal(t, decode, required_decode)

	// [Reduced] No Delimiter | 0x00
	config = Config{
		SpecialByte: 0x00,
		Delimiter:   false,
		Type:        Reduced,
	}
	required_message = "aaaaaaaaaaa"
	required_raw = []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_encode = []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	required_decode = required_raw

	raw = []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode = Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode = Decode(encode, config)
	assert.Equal(t, decode, required_decode)

	// [Reduced] No Delimiter | 0x61
	config = Config{
		SpecialByte: 0x61,
		Delimiter:   false,
		Type:        Reduced,
	}
	required_message = "aabbbaabbabbabb"
	required_raw = []byte{97, 97, 98, 98, 98, 97, 97, 98, 98, 97, 98, 98, 97, 98, 98}
	required_encode = []byte{1, 1, 4, 98, 98, 98, 1, 3, 98, 98, 3, 98, 98, 98, 98}
	required_decode = required_raw

	raw = []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode = Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode = Decode(encode, config)
	assert.Equal(t, decode, required_decode)
}

func TestFlagCounting(t *testing.T) {
	config := Config{
		SpecialByte: 0x61,
		Delimiter:   true,
		Type:        Reduced,
	}
	required_message := "aabbbaabbabbabb"
	required_raw := []byte{97, 97, 98, 98, 98, 97, 97, 98, 98, 97, 98, 98, 97, 98, 98}
	required_encode := []byte{1, 1, 4, 98, 98, 98, 1, 3, 98, 98, 3, 98, 98, 98, 98, 97}
	required_decode := required_raw

	raw := []byte(required_message)
	assert.Equal(t, raw, required_raw)
	encode := Encode(raw, config)
	assert.Equal(t, encode, required_encode)
	assert.Equal(t, Verify(encode, config), nil)
	decode := Decode(encode, config)
	assert.Equal(t, decode, required_decode)

	assert.Equal(t, FlagCount(encode, config), 8)
}
