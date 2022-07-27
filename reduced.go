package cobs

import (
	"errors"
)

func reducedEncode(src []byte, config Config) (dst []byte) {
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
	if int(src[srcLen-1]) > (len(dst)-codePtr) && src[srcLen-1] != config.SpecialByte {
		code = src[srcLen-1]
		dst = dst[:len(dst)-1]
		dst[codePtr] = code
	}else{
		if code == config.SpecialByte {
			dst[codePtr] = 0x00
		} else {
			dst[codePtr] = code
		}
	}
	if config.Delimiter {
		dst = append(dst, config.SpecialByte)
	}
	return dst
}

func reducedDecode(src []byte, config Config) (dst []byte) {
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
	return append(dst, code)
}

func reducedVerify(src []byte, config Config) (err error) {
	loopLen := len(src)
	if config.Delimiter {
		if loopLen < 2 {
			return errors.New("COBS[Reduced]: Encoded slice is too short to be valid.")
		}
		if src[loopLen-1] != config.SpecialByte {
			return errors.New("COBS[Reduced]: Encoded slice's delimiter is not special byte.")
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return errors.New("COBS[Reduced]: Encoded slice is too short to be valid.")
		}
	}
	for _, b := range src[:loopLen] {
		if b == config.SpecialByte {
			return errors.New("COBS[Native]: Encoded slice's byte (not the delimter) is special byte.")
		}
	}
	return nil
}