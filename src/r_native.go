package cobs

import (
	"errors"
)

type reversedNativeEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (e reversedNativeEncoder) Encode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 0, loopLen+1)
	code := byte(0x01)
	ptr := 0
	for ptr < loopLen {
		if src[ptr] == e.SpecialByte {
			if code == e.SpecialByte {
				dst = append(dst, 0)
			} else {
				dst = append(dst, code)
			}
			code = 0x01
			ptr++
			continue
		}
		dst = append(dst, src[ptr])
		code++
		if code == 0xFF && (!e.EndingSave || ptr != loopLen-1) {
			if code == e.SpecialByte {
				dst = append(dst, 0)
			} else {
				dst = append(dst, code)
			}
			code = 0x01
		}
		ptr++
	}
	if code == e.SpecialByte {
		dst = append(dst, 0)
	} else {
		dst = append(dst, code)
	}
	if e.Delimiter {
		dst = append(dst, e.SpecialByte)
	}
	return dst
}

func (e reversedNativeEncoder) Decode(src []byte) []byte {
	ptr := len(src) - 1
	dst := make([]byte, 0, ptr-(ptr/254))
	if e.Delimiter {
		ptr--
	}
	code := byte(0x00)
	for ptr >= 0 {
		if src[ptr] == 0x00 {
			code = e.SpecialByte
		} else {
			code = src[ptr]
		}
		if code < 0xFF || (e.EndingSave && ptr < 0) {
			dst = append([]byte{e.SpecialByte}, dst...)
		}
		ptr--
		for i := 1; i < int(code); i++ {
			dst = append([]byte{src[ptr]}, dst...)
			ptr--
		}
	}
	return dst[:len(dst)-1]
}

func (e reversedNativeEncoder) Verify(src []byte) error {
	nextFlag := 0
	ptr := len(src)
	if e.Delimiter {
		if ptr < 2 {
			return errors.New("COBS[Native]: Encoded slice is too short to be valid.")
		}
		if src[ptr-1] != e.SpecialByte {
			return errors.New("COBS[Native]: Encoded slice's delimiter is not special byte.")
		}
		ptr--
	} else {
		if ptr < 1 {
			return errors.New("COBS[Native]: Encoded slice is too short to be valid.")
		}
	}
	ptr--
	for ptr >= 0 {
		if src[ptr] == e.SpecialByte {
			return errors.New("COBS[Native]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if src[ptr] == 0 {
				nextFlag = int(e.SpecialByte)
			} else {
				nextFlag = int(src[ptr])
			}
		}
		ptr--
		nextFlag--
	}
	if nextFlag != 0 {
		return errors.New("COBS[Native]: Encoded slice's flags do not lead to end.")
	}
	return nil
}

func (e reversedNativeEncoder) FlagCount(src []byte) int {
	numFlags := 0
	ptr := len(src) - 1
	if e.Delimiter {
		ptr--
		numFlags++
	}
	for ptr >= 0 {
		if src[ptr] == 0 {
			ptr -= int(e.SpecialByte)
		} else {
			ptr -= int(src[ptr])
		}
		numFlags++
	}
	return numFlags
}

func (e reversedNativeEncoder) MaxOverhead(len int) int {
	ret := len + 1 + (len / 254)
	if e.Delimiter {
		ret++
	}
	return ret
}

func (e reversedNativeEncoder) MinOverhead(len int) int {
	ret := len + 1
	if e.Delimiter {
		ret++
	}
	return ret
}
