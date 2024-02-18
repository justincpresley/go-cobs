package main

import (
	"fmt"

	cobs "github.com/justincpresley/go-cobs/pkg"
)

func main() {
	config := cobs.Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		Type:        cobs.Reduced,
		EndingSave:  true,
		Reverse:     false,
	}
	encoder, ok := cobs.NewEncoder(config)
	if ok != nil {
		fmt.Println("Error:", ok)
		return
	}

	message := "AAAAAAAAAAAAAAAAA"
	raw := []byte(message)

	fmt.Println("Config | Special", config.SpecialByte, "Delimiter", config.Delimiter, "Type", config.Type, "|")
	fmt.Println("Message:", message)
	fmt.Println("Message Bytes:", raw)

	encoded := encoder.Encode(raw)
	fmt.Println("Encoded:", encoded)

	if ok = encoder.Verify(encoded); ok != nil {
		fmt.Println("Status: CORRUPTED")
		fmt.Println("Error:", ok)
		return
	}
	fmt.Println("Status: VALID")

	decoded := encoder.Decode(encoded)
	fmt.Println("Decoded:", decoded)
	fmt.Println("Message:", string(decoded))
}
