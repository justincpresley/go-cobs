# Consistent Overhead Byte Stuffing (COBS) Library In Go

This Go library provides an API which allows all COBS-related functionality to the programmer.

*What is Consistent Overhead Byte Stuffing (COBS)?*

Links: [Wiki](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) - [Technical Paper](http://www.stuartcheshire.org/papers/cobsforton.pdf)

Effectively, the goal of COBS is to remove a special byte within given data by replacing all special bytes with "flags", a byte telling where the next special byte is. There is minimal overhead in COBS as indicated by the paper itself.

Comparatively, this library offers combined functionality not found within any other COBS library on the market as of writing this. To name a few:

 - **Choose whether the Delimiter (Special byte) is added.** Libraries tend to choose this for you.
 - **Choose what the Special Byte is.** Every other library uses the NULL (0x00) byte only.
 - **Use COBS as a Layer of Integrity.** By ensuring that the special byte does not occur (expect with a delimiter) and by ensuring that the flags lead to the end of the data, COBS can provide a small layer of integrity.
 - **Use a TONS of different extensions (types).** I have included as many types as I could find and have included some of my own.
 - **Additional API Commands.** Not only is there functions to calculate the worst/best case for COBS, but there is also flag-related functions.

## Documentation / Usage

The full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/justincpresley/go-cobs).

Additionally, [use](https://github.com/justincpresley/go-cobs/blob/master/USE.md) outlines everything you need to know about types, config, and anything that may be unclear outside of the API.

## Installation

```
go get -u github.com/justincpresley/go-cobs
```

Or use your favorite golang vendoring tool.

## Example

It should be stated that the `Config` should be the same for a given use-case or layer/grouping. For example the go-to use case in networking, a frame on layer 2 might want the Special Byte be `0x00` and include a delimiter to mark where a frame ends.

```go
package main

import (
	"fmt"

	"github.com/justincpresley/go-cobs"
)

func main() {
  var ok error

	config := cobs.Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		Type:        Reduced,
		EndingSave:  true,
    Reverse:     false,
	}
  if ok = config.Validate(); ok != nil {
		fmt.Println("Error:", ok)
		return
	}

	message := "AAAAAAAAAAAAAAAAA"
	raw := []byte(message)

	fmt.Println("Config | Special", config.SpecialByte, "Delimiter", config.Delimiter, "Type", config.Type, "|")
	fmt.Println("Message:", message)
	fmt.Println("Message Bytes:", raw)

	encoded := cobs.Encode(raw, config)
	fmt.Println("Encoded:", encoded)

	if ok = cobs.Verify(encoded, config); ok != nil {
		fmt.Println("Status: CORRUPTED")
		fmt.Println("Error:", ok)
		return
	}
	fmt.Println("Status: VALID")

	decoded := cobs.Decode(encoded, config)
	fmt.Println("Decoded:", decoded)
	fmt.Println("Message:", string(decoded))
}
```

## Notes

There is a lot of research potential out of this library. For example, It would be an interesting to find what the optimal amount of flags is to provide the best integrity checking for given data.

The API is subject to change based on future additions.
Speaking of additions, I am more than happy to review pull requests if anyone wants to contribute. See [future](https://github.com/justincpresley/go-cobs/blob/master/FUTURE.md) for more details on what I plan to implement and for specific ways you can help out!
