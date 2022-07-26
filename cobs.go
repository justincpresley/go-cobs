package cobs

// Config is a struct that holds configuration variables on how
// this COBS library will function. It can be customized according
// to the use case.
type Config struct {
	SpecialByte byte
	Delimiter   bool
}

// Encode takes raw data and a configuration and produces the COBS-encoded
// byte slice.
func Encode(src []byte, config Config) (dst []byte) {
	srcLen := len(src) + 1
	dst = make([]byte, 1, srcLen)
	codePtr := 0
	code := byte(0x01)
	for _, b := range src {
		if b == config.SpecialByte {
			if code == config.SpecialByte {
				dst[codePtr] = 0x00
			} else {
				dst[codePtr] = code
			}
			codePtr = len(dst)
			dst = append(dst, 0)
			code = 0x01
			continue
		}
		dst = append(dst, b)
		code++
		if code == 0xFF {
			if code == config.SpecialByte {
				dst[codePtr] = 0x00
			} else {
				dst[codePtr] = code
			}
			codePtr = len(dst)
			dst = append(dst, 0)
			code = 0x01
		}
	}
	if code == config.SpecialByte {
		dst[codePtr] = 0x00
	} else {
		dst[codePtr] = code
	}
	if config.Delimiter {
		dst = append(dst, config.SpecialByte)
	}
	return dst
}

// Decode takes encoded data and a configuration and produces the raw COBS-decoded
// byte slice.
func Decode(src []byte, config Config) (dst []byte) {
	loopLen := len(src)
	// the cap needs optimization
	dst = make([]byte, 0, loopLen-1-(loopLen/254))
	if config.Delimiter {
		loopLen--
	}
	ptr := 0
	code := byte(0x00)
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			code = config.SpecialByte
		} else {
			code = src[ptr]
		}
		ptr++
		for i := 1; i < int(code); i++ {
			dst = append(dst, src[ptr])
			ptr++
		}
		if code < 0xFF && ptr < loopLen {
			dst = append(dst, config.SpecialByte)
		}
	}
	return dst
}

// Verify checks whether the given raw data can be a valid COBS-encoded byte slice
// based on the configuration. It checks to see if the special byte appears and
// whether the flags -lead- toward the end of the slice.
func Verify(src []byte, config Config) (success bool) {
	nextFlag := 0
	srcLen := len(src)
	if config.Delimiter {
		if srcLen < 2 {
			return false
		}
		for _, b := range src[:srcLen-1] {
			if b == config.SpecialByte {
				return false
			}
			if nextFlag == 0 {
				if b == 0 {
					nextFlag = int(config.SpecialByte)
				} else {
					nextFlag = int(b)
				}
			}
			nextFlag--
		}
		if nextFlag != 0 || src[srcLen-1] != config.SpecialByte {
			return false
		}
	} else {
		if srcLen < 1 {
			return false
		}
		for _, b := range src[:srcLen] {
			if b == config.SpecialByte {
				return false
			}
			if nextFlag == 0 {
				if b == 0 {
					nextFlag = int(config.SpecialByte)
				} else {
					nextFlag = int(b)
				}
			}
			nextFlag--
		}
		if nextFlag != 0 {
			return false
		}
	}
	return true
}

// Worse Case calculates the worse case for the COBS overhead when given
// a raw length and an appropiate configuration.
func WorseCase(dLen int, config Config) (eLen int) {
	if config.Delimiter {
		return dLen + 2 + (dLen / 254)
	} else {
		return dLen + 1 + (dLen / 254)
	}
}

// Best Case calculates the best case for the COBS overhead when given
// a raw length and an appropiate configuration.
func BestCase(dLen int, config Config) (eLen int) {
	if config.Delimiter {
		return dLen + 2
	} else {
		return dLen + 1
	}
}
