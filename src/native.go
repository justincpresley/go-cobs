package cobs

import (
	"errors"
)

type nativeEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (e nativeEncoder) Encode(src []byte) []byte {
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
	if code == e.SpecialByte {
		dst[codePtr] = 0x00
	} else {
		dst[codePtr] = code
	}
	if e.Delimiter {
		dst = append(dst, e.SpecialByte)
	}
	return dst
}

func (e nativeEncoder) Decode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 0, loopLen-1-(loopLen/254))
	if e.Delimiter {
		loopLen--
	}
	ptr := 0
	code := byte(0x00)
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			code = e.SpecialByte
		} else {
			code = src[ptr]
		}
		ptr++
		for i := 1; i < int(code); i++ {
			dst = append(dst, src[ptr])
			ptr++
		}
		if code < 0xFF || (e.EndingSave && ptr == loopLen) {
			dst = append(dst, e.SpecialByte)
		}
	}
	return dst[:len(dst)-1]
}

func (e nativeEncoder) Verify(src []byte) error {
	nextFlag := 0
	loopLen := len(src)
	if e.Delimiter {
		if loopLen < 2 {
			return errors.New("COBS[Native]: Encoded slice is too short to be valid.")
		}
		if src[loopLen-1] != e.SpecialByte {
			return errors.New("COBS[Native]: Encoded slice's delimiter is not special byte.")
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return errors.New("COBS[Native]: Encoded slice is too short to be valid.")
		}
	}
	for _, b := range src[:loopLen] {
		if b == e.SpecialByte {
			return errors.New("COBS[Native]: Encoded slice's byte (not the delimter) is special byte.")
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
	if nextFlag != 0 {
		return errors.New("COBS[Native]: Encoded slice's flags do not lead to end.")
	}
	return nil
}

func (e nativeEncoder) FlagCount(src []byte) int {
	ret := 0
	ptr := 0
	for ptr < len(src) {
		if src[ptr] == 0 {
			ptr += int(e.SpecialByte)
		} else {
			ptr += int(src[ptr])
		}
		ret++
	}
	return ret
}

func (e nativeEncoder) MaxOverhead(len int) int {
	ret := len + 1 + (len / 254)
	if e.Delimiter {
		ret++
	}
	return ret
}

func (e nativeEncoder) MinOverhead(len int) int {
	ret := len + 1
	if e.Delimiter {
		ret++
	}
	return ret
}
