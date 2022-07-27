package cobs

import (
	"fmt"
	"testing"
)

func TestBasicFeatures(t *testing.T) {
	config := Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		Type:        Reduced,
	}
	message := "AAAAAAAAAAAAA"

	raw := []byte(message)
	fmt.Println("Config | Special", config.SpecialByte, "Delimiter", config.Delimiter, "Type", config.Type, "|")
	fmt.Println("Message:", message)
	fmt.Println("Message Bytes:", raw)

	encoded := Encode(raw, config)
	fmt.Println("Encoded:", encoded)

	if !Verify(encoded, config) {
		fmt.Println("Status: CORRUPTED")
	}
	fmt.Println("Status: VALID")

	decoded := Decode(encoded, config)
	fmt.Println("Decoded:", decoded)
	fmt.Println("Message:", string(decoded))
}
