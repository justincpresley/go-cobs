package cobs

func pairelimEncode(src []byte, config Config) (dst []byte) {
	srcLen := len(src)
	dst = make([]byte, 1, srcLen+1)
	codePtr := 0
	code := byte(0x01)
	pairable := false
	for _, b := range src {
		if pairable {
			pairable = false
			if b == config.SpecialByte {
				code |= 0xE0
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
			if code == config.SpecialByte {
				dst[codePtr] = 0x00
			} else {
				dst[codePtr] = code
			}
			codePtr = len(dst)
			dst = append(dst, 0)
			dst = append(dst, b)
			code = 0x01
			code++
			continue
		}
		if b == config.SpecialByte {
			if code <= 0x1F {
				pairable = true
				continue
			}
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
		if code == 0xE0 {
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
	if pairable {
		code |= 0xE0
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

func pairelimDecode(src []byte, config Config) (dst []byte) {
	loopLen := len(src)
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
		if code > 0xE0 {
			jumpLen = int(code & 0x1F)
		}else{
			jumpLen = int(code)
		}
		ptr++
		for i:=1; i<jumpLen; i++ {
			dst = append(dst, src[ptr])
			ptr++
		}
		switch {
		case code > 0xE0:
			dst = append(dst, config.SpecialByte)
			dst = append(dst, config.SpecialByte)
		case code < 0xE0:
			dst = append(dst, config.SpecialByte)
		}
	}
	return dst[:len(dst)-1]
}

func pairelimVerify(src []byte, config Config) (err error) {
	return
}

func pairelimFLagCount(src []byte, config Config) (flags int) {
	numFlags := 0
	ptr := 0
	loopLen := len(src)
	if config.Delimiter {
		numFlags++
		loopLen--
	}
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			if config.SpecialByte > 0xE0 {
				ptr += int(config.SpecialByte & 0x1F)
			}else{
				ptr += int(config.SpecialByte)
			}
		} else {
			if src[ptr] > 0xE0 {
				ptr += int(src[ptr] & 0x1F)
			}else{
				ptr += int(src[ptr])
			}
		}
		numFlags++
	}
	return numFlags
}