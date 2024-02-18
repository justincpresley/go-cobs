package encoders

import (
	"errors"
)

type ReducedEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (enc ReducedEncoder) Encode(src []byte) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 1, loopLen+1)
	codePtr := 0
	code := byte(0x01)
	ptr := 0
	for ptr < loopLen {
		if src[ptr] == enc.SpecialByte {
			if code == enc.SpecialByte {
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
		if code == 0xFF && (!enc.EndingSave || ptr != loopLen-1) {
			if code == enc.SpecialByte {
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
	if loopLen != 0 && int(src[loopLen-1]) > (len(dst)-codePtr) && src[loopLen-1] != enc.SpecialByte {
		code = src[loopLen-1]
		dst = dst[:len(dst)-1]
		dst[codePtr] = code
	} else {
		if code == enc.SpecialByte {
			dst[codePtr] = 0x00
		} else {
			dst[codePtr] = code
		}
	}
	if enc.Delimiter {
		dst = append(dst, enc.SpecialByte)
	}
	return dst
}

func (enc ReducedEncoder) Decode(src []byte) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 0, loopLen-1-(loopLen/254))
	if enc.Delimiter {
		loopLen--
	}
	ptr := 0
	jumpLen := 0
	code := byte(0x00)
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			code = enc.SpecialByte
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
		if code < 0xFF || (enc.EndingSave && ptr == loopLen) {
			dst = append(dst, enc.SpecialByte)
		}
	}
	dst = dst[:len(dst)-1]
	if jumpLen != int(code) {
		dst = append(dst, code)
	}
	return dst
}

func (enc ReducedEncoder) Verify(src []byte) (err error) {
	nextFlag := 0
	loopLen := len(src)
	if enc.Delimiter {
		if loopLen < 2 {
			return errors.New("COBS[Reduced]: Encoded slice is too short to be valid.")
		}
		if src[loopLen-1] != enc.SpecialByte {
			return errors.New("COBS[Reduced]: Encoded slice's delimiter is not special byte.")
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return errors.New("COBS[Reduced]: Encoded slice is too short to be valid.")
		}
	}
	for _, b := range src[:loopLen] {
		if b == enc.SpecialByte {
			return errors.New("COBS[Reduced]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if b == 0 {
				nextFlag = int(enc.SpecialByte)
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

func (enc ReducedEncoder) FlagCount(src []byte) (flags int) {
	numFlags := 0
	ptr := 0
	for ptr < len(src) {
		if src[ptr] == 0 {
			ptr += int(enc.SpecialByte)
		} else {
			ptr += int(src[ptr])
		}
		numFlags++
	}
	if ptr != len(src) && enc.Delimiter {
		numFlags++
	}
	return numFlags
}

func (enc ReducedEncoder) WorseCase(dLen int) (eLen int) {
	eLen = dLen + 1 + (dLen / 254)
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc ReducedEncoder) MaxOverhead(dLen int) (eLen int) {
	return enc.WorseCase(dLen)
}

func (enc ReducedEncoder) BestCase(dLen int) (eLen int) {
	eLen = dLen
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc ReducedEncoder) MinOverhead(dLen int) (eLen int) {
	return enc.BestCase(dLen)
}
