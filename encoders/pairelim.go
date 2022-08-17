package encoders

import (
	"errors"
)

type PairelimEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (enc PairelimEncoder) Encode(src []byte) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 1, loopLen+1)
	codePtr := 0
	code := byte(0x01)
	pairable := false
	ptr := 0
	for ptr < loopLen {
		if pairable {
			pairable = false
			if src[ptr] == enc.SpecialByte {
				code |= 0xE0
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
			if code == enc.SpecialByte {
				dst[codePtr] = 0x00
			} else {
				dst[codePtr] = code
			}
			codePtr = len(dst)
			dst = append(dst, 0)
			dst = append(dst, src[ptr])
			code = 0x01
			code++
			ptr++
			continue
		}
		if src[ptr] == enc.SpecialByte {
			if code <= 0x1F {
				pairable = true
				ptr++
				continue
			}
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
		if code == 0xE0 && (!enc.EndingSave || ptr != loopLen-1) {
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
	if pairable {
		code |= 0xE0
	}
	if code == enc.SpecialByte {
		dst[codePtr] = 0x00
	} else {
		dst[codePtr] = code
	}
	if enc.Delimiter {
		dst = append(dst, enc.SpecialByte)
	}
	return dst
}

func (enc PairelimEncoder) Decode(src []byte) (dst []byte) {
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
		if code > 0xE0 {
			jumpLen = int(code & 0x1F)
		} else {
			jumpLen = int(code)
		}
		ptr++
		for i := 1; i < jumpLen; i++ {
			dst = append(dst, src[ptr])
			ptr++
		}
		if code > 0xE0 {
			dst = append(dst, enc.SpecialByte)
			dst = append(dst, enc.SpecialByte)
		} else if code < 0xE0 || (enc.EndingSave && ptr == loopLen) {
			dst = append(dst, enc.SpecialByte)
		}
	}
	return dst[:len(dst)-1]
}

func (enc PairelimEncoder) Verify(src []byte) (err error) {
	nextFlag := 0
	loopLen := len(src)
	if enc.Delimiter {
		if loopLen < 2 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
		if src[loopLen-1] != enc.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's delimiter is not special byte.")
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
	}
	for _, b := range src[:loopLen] {
		if b == enc.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if b == 0x00 {
				if enc.SpecialByte > 0xE0 {
					nextFlag = int(enc.SpecialByte & 0x1F)
				} else {
					nextFlag = int(enc.SpecialByte)
				}
			} else {
				if b > 0xE0 {
					nextFlag = int(b & 0x1F)
				} else {
					nextFlag = int(b)
				}
			}
		}
		nextFlag--
	}
	if nextFlag != 0 {
		return errors.New("COBS[PairElimination]: Encoded slice's flags do not lead to end.")
	}
	return nil
}

func (enc PairelimEncoder) FlagCount(src []byte) (flags int) {
	numFlags := 0
	ptr := 0
	loopLen := len(src)
	if enc.Delimiter {
		numFlags++
		loopLen--
	}
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			if enc.SpecialByte > 0xE0 {
				ptr += int(enc.SpecialByte & 0x1F)
			} else {
				ptr += int(enc.SpecialByte)
			}
		} else {
			if src[ptr] > 0xE0 {
				ptr += int(src[ptr] & 0x1F)
			} else {
				ptr += int(src[ptr])
			}
		}
		numFlags++
	}
	return numFlags
}

func (enc PairelimEncoder) WorseCase(dLen int) (eLen int) {
	eLen = dLen + 1 + (dLen / 223)
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc PairelimEncoder) MaxOverhead(dLen int) (eLen int) {
	return enc.WorseCase(dLen)
}

func (enc PairelimEncoder) BestCase(dLen int) (eLen int) {
	eLen = (dLen / 2) + 1
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc PairelimEncoder) MinOverhead(dLen int) (eLen int) {
	return enc.BestCase(dLen)
}