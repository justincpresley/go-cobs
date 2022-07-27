package cobs

func nativeEncode(src []byte, config Config) (dst []byte) {
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

func nativeDecode(src []byte, config Config) (dst []byte) {
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

func nativeVerify(src []byte, config Config) (success bool) {
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
	if nextFlag != 0 {
		return false
	}
	return true
}