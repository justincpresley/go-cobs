package cobs

import (
	"errors"
)

type reducedEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (e reducedEncoder) Encode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 1, loopLen+1)
	codePtr := 0
	code := byte(0x01)
	ptr := 0
	for ptr < loopLen {
		if src[ptr] == e.SpecialByte {
			if code == e.SpecialByte {
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
		if code == 0xFF && (!e.EndingSave || ptr != loopLen-1) {
			if code == e.SpecialByte {
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
	if loopLen != 0 && int(src[loopLen-1]) > (len(dst)-codePtr) && src[loopLen-1] != e.SpecialByte {
		code = src[loopLen-1]
		dst = dst[:len(dst)-1]
		dst[codePtr] = code
	} else {
		if code == e.SpecialByte {
			dst[codePtr] = 0x00
		} else {
			dst[codePtr] = code
		}
	}
	if e.Delimiter {
		dst = append(dst, e.SpecialByte)
	}
	return dst
}

func (e reducedEncoder) Decode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 0, loopLen-1-(loopLen/254))
	if e.Delimiter {
		loopLen--
	}
	ptr := 0
	jumpLen := 0
	code := byte(0x00)
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			code = e.SpecialByte
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
		if code < 0xFF || (e.EndingSave && ptr == loopLen) {
			dst = append(dst, e.SpecialByte)
		}
	}
	dst = dst[:len(dst)-1]
	if jumpLen != int(code) {
		dst = append(dst, code)
	}
	return dst
}

func (e reducedEncoder) Verify(src []byte) error {
	nextFlag := 0
	loopLen := len(src)
	if e.Delimiter {
		if loopLen < 2 {
			return errors.New("COBS[Reduced]: Encoded slice is too short to be valid.")
		}
		if src[loopLen-1] != e.SpecialByte {
			return errors.New("COBS[Reduced]: Encoded slice's delimiter is not special byte.")
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return errors.New("COBS[Reduced]: Encoded slice is too short to be valid.")
		}
	}
	for _, b := range src[:loopLen] {
		if b == e.SpecialByte {
			return errors.New("COBS[Reduced]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if b == 0 {
				nextFlag = int(e.SpecialByte)
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

func (e reducedEncoder) FlagCount(src []byte) int {
	numFlags := 0
	ptr := 0
	for ptr < len(src) {
		if src[ptr] == 0 {
			ptr += int(e.SpecialByte)
		} else {
			ptr += int(src[ptr])
		}
		numFlags++
	}
	if ptr != len(src) && e.Delimiter {
		numFlags++
	}
	return numFlags
}

func (e reducedEncoder) MaxOverhead(len int) int {
	ret := len + 1 + (len / 254)
	if e.Delimiter {
		ret++
	}
	return ret
}

func (e reducedEncoder) MinOverhead(len int) int {
	ret := len
	if e.Delimiter {
		ret++
	}
	return ret
}
