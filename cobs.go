package cobs

// Config is a struct that holds configuration variables on how
// this COBS library will function. It can be customized according
// to the use case.
type Config struct {
	SpecialByte byte
	Delimiter   bool
	Reduced     bool
}

// Encode takes raw data and a configuration and produces the COBS-encoded
// byte slice.
func Encode(src []byte, config Config) (dst []byte) {
	srcLen := len(src)
	dst = make([]byte, 1, srcLen+1)
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
	if config.Reduced {
		if int(src[srcLen-1]) > (len(dst)-codePtr) && src[srcLen-1] != config.SpecialByte {
			code = src[srcLen-1]
			dst = dst[:len(dst)-1]
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
	jumpLen := 0
	code := byte(0x00)
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			code = config.SpecialByte
		} else {
			code = src[ptr]
		}
		ptr++
		if int(code) > (loopLen - ptr) {
			jumpLen = loopLen - ptr + 1
		} else {
			jumpLen = int(code)
		}
		for i := 1; i < jumpLen; i++ {
			dst = append(dst, src[ptr])
			ptr++
		}
		if code < 0xFF && ptr < loopLen {
			dst = append(dst, config.SpecialByte)
		}
	}
	if config.Reduced {
		dst = append(dst, code)
	}
	return dst
}

// Verify checks whether the given raw data can be a valid COBS-encoded byte slice
// based on the configuration. It checks to see if the special byte appears and
// whether the flags -lead- toward the end of the slice.
func Verify(src []byte, config Config) (success bool) {
	nextFlag := 0
	loopLen := len(src)
	if config.Delimiter {
		if loopLen < 2 || src[loopLen-1] != config.SpecialByte {
			return false
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return false
		}
	}
	for _, b := range src[:loopLen] {
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
	if nextFlag != 0 && !config.Reduced {
		return false
	}
	return true
}

// Worse Case calculates the worse case for the COBS overhead when given
// a raw length and an appropiate configuration.
func WorseCase(dLen int, config Config) (eLen int) {
	eLen = dLen + 1 + (dLen / 254)
	if config.Delimiter {
		eLen++
	}
	return eLen
}

// Best Case calculates the best case for the COBS overhead when given
// a raw length and an appropiate configuration.
func BestCase(dLen int, config Config) (eLen int) {
	eLen = dLen + 1
	if config.Delimiter {
		eLen++
	}
	if config.Reduced {
		eLen--
	}
	return eLen
}
