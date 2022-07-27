# Consistent Overhead Byte Stuffing (COBS) Library In Go

This Go library provides an API which allows all COBS-related functionality to the programmer.

*What is Consistent Overhead Byte Stuffing (COBS)?*

Links: [Wiki](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) - [Technical Paper](http://www.stuartcheshire.org/papers/cobsforton.pdf)

Effectively, the goal of COBS is to remove a special byte within given data by replacing all special bytes with "flags", a byte telling where the next special byte is. There is minimal overhead in COBS as indicated by the paper itself.

Comparatively, this library offers combined functionality not found within any other COBS library on the market as of writing this:

 - **Be able to choose whether the Delimiter (Special byte) is added.** Libraries tend to only given one option.
 - **Be able to choose what the Special Byte is.** Every other library uses the NULL (0x00) byte only.
 - **Be able to use COBS as a Layer of Integrity.** By ensuring that the special byte does not occur (expect with a delimiter) and by ensuring that the flags lead to the end of the data, COBS can provide a small layer of integrity.
 - **Be able to use COBS/R (Reduced overhead).** This method is created and described [here](https://github.com/cmcqueen/cobs-c). Note by reducing, the flag-based check method for integrity is not applicable.
 - **Additional API Commands.** Not only is there functions to calculate the worst/best case for COBS, but there is also flag-related functions.

## Documentation

[The full documentation is available on pkg.go.dev](https://pkg.go.dev/github.com/justincpresley/go-cobs).

## Installation

```
go get -u github.com/justincpresley/go-cobs
```

Or use your favorite golang vendoring tool.

## Example / usage

It should be stated that the `Config` should be the same for a given use-case or layer/grouping. For example the go-to use case in networking, a frame on layer 2 might want the Special Byte be `0x00` and include a delimiter to mark where a frame ends.

```go
package main

import (
	"fmt"

	"github.com/justincpresley/go-cobs"
)

func main() {
	var config cobs.Config
	config.SpecialByte = 0x00
	config.Delimiter= true

	message := "AAAAAAAAAAAAAAAAA"
	raw := []byte(message)

	fmt.Println("Config Special", config.SpecialByte, "Delimiter", config.Delimiter)
	fmt.Println("Message:", message)
	fmt.Println("Message Bytes:", raw)

	encoded := cobs.Encode(raw, config)
	fmt.Println("Encoded:", encoded)

	if !cobs.Verify(encoded, config) {
		fmt.Println("Status: CORRUPTED")
	}
	fmt.Println("Status: VALID")

	decoded := cobs.Decode(encoded, config)
	fmt.Println("Decoded:", decoded)
	fmt.Println("Message:", string(decoded))
}
```

## Notes

It would be an interesting to find what the optimal amount of flags is to provide the best integrity checking for given data.

The API is subject to check as plan to incorporate the following features in future releases:

 - Create `FlagCount()` function.
 - A function that combines `Decode()` and `Verify()` in one function to eliminate the need to loop through the data twice.
 - Implement `COBS/ZPE` and/or `COBS/ZRE`.