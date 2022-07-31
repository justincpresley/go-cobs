package cobs

import (
	"errors"
)

func reducedEncode(src []byte, config Config) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 1, loopLen+1)
	codePtr := 0
	code := byte(0x01)
	ptr := 0
	for ptr < loopLen {
		if src[ptr] == config.SpecialByte {
			if code == config.SpecialByte {
				dst[codePtr] = 0x00
			} else {
				dst[codePtr] = code
			}
			codePtr = len(dst)
			dst = append(dst, 0)
			code = 0x01
			ptr++
			continue
		}
		dst = append(dst, src[ptr])
		code++
		if code == 0xFF && (!config.EndingSave || ptr != loopLen-1) {
			if code == config.SpecialByte {
				dst[codePtr] = 0x00
			} else {
				dst[codePtr] = code
			}
			codePtr = len(dst)
			dst = append(dst, 0)
			code = 0x01
		}
		ptr++
	}
	if loopLen != 0 && int(src[loopLen-1]) > (len(dst)-codePtr) && src[loopLen-1] != config.SpecialByte {
		code = src[loopLen-1]
		dst = dst[:len(dst)-1]
		dst[codePtr] = code
	} else {
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
		if code < 0xFF || (config.EndingSave && ptr == loopLen) {
			dst = append(dst, config.SpecialByte)
		}
	}
	dst = dst[:len(dst)-1]
	if jumpLen != int(code) {
		dst = append(dst, code)
	}
	return dst
}

func reducedVerify(src []byte, config Config) (err error) {
	nextFlag := 0
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
			return errors.New("COBS[Reduced]: Encoded slice's byte (not the delimter) is special byte.")
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
	if nextFlag < 0 {
		return errors.New("COBS[Reduced]: Encoded slice's flags do not lead to end.")
	}
	return nil
}

func reducedFlagCount(src []byte, config Config) (flags int) {
	numFlags := 0
	ptr := 0
	for ptr < len(src) {
		if src[ptr] == 0 {
			ptr += int(config.SpecialByte)
		} else {
			ptr += int(src[ptr])
		}
		numFlags++
	}
	if ptr != len(src) && config.Delimiter {
		numFlags++
	}
	return numFlags
}
