package encoders

import (
	"errors"
)

type R_PairelimEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (enc R_PairelimEncoder) Encode(src []byte) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 0, loopLen+1)
	code := byte(0x01)
	pairable := false
	ptr := 0
	for ptr < loopLen {
		if pairable {
			pairable = false
			if src[ptr] == enc.SpecialByte {
				code |= 0xE0
				if code == enc.SpecialByte {
					dst = append(dst, 0)
				} else {
					dst = append(dst, code)
				}
				code = 0x01
				ptr++
				continue
			}
			if code == enc.SpecialByte {
				dst = append(dst, 0)
			} else {
				dst = append(dst, code)
			}
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
		if code == 0xE0 && (!enc.EndingSave || ptr != loopLen-1) {
			if code == enc.SpecialByte {
				dst = append(dst, 0)
			} else {
				dst = append(dst, code)
			}
			code = 0x01
		}
		ptr++
	}
	if pairable {
		code |= 0xE0
	}
	if code == enc.SpecialByte {
		dst = append(dst, 0)
	} else {
		dst = append(dst, code)
	}
	if enc.Delimiter {
		dst = append(dst, enc.SpecialByte)
	}
	return dst
}

func (enc R_PairelimEncoder) Decode(src []byte) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 0, loopLen-1-(loopLen/254))
	if enc.Delimiter {
		loopLen--
	}
	ptr := loopLen - 1
	jumpLen := 0
	code := byte(0x00)
	for ptr >= 0 {
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
		if code > 0xE0 {
			dst = append([]byte{enc.SpecialByte}, dst...)
			dst = append([]byte{enc.SpecialByte}, dst...)
		} else if code < 0xE0 || (enc.EndingSave && ptr == loopLen) {
			dst = append([]byte{enc.SpecialByte}, dst...)
		}
		ptr--
		for i := 1; i < jumpLen; i++ {
			dst = append([]byte{src[ptr]}, dst...)
			ptr--
		}
	}
	return dst[:len(dst)-1]
}

func (enc R_PairelimEncoder) Verify(src []byte) (err error) {
	nextFlag := 0
	ptr := len(src)
	if enc.Delimiter {
		if ptr < 2 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
		if src[ptr-1] != enc.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's delimiter is not special byte.")
		}
		ptr--
	} else {
		if ptr < 1 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
	}
	ptr--
	for ptr >= 0 {
		if src[ptr] == enc.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if src[ptr] == 0x00 {
				if enc.SpecialByte > 0xE0 {
					nextFlag = int(enc.SpecialByte & 0x1F)
				} else {
					nextFlag = int(enc.SpecialByte)
				}
			} else {
				if src[ptr] > 0xE0 {
					nextFlag = int(src[ptr] & 0x1F)
				} else {
					nextFlag = int(src[ptr])
				}
			}
		}
		ptr--
		nextFlag--
	}
	if nextFlag != 0 {
		return errors.New("COBS[PairElimination]: Encoded slice's flags do not lead to end.")
	}
	return nil
}

func (enc R_PairelimEncoder) FlagCount(src []byte) (flags int) {
	numFlags := 0
	ptr := len(src) - 1
	if enc.Delimiter {
		ptr--
		numFlags++
	}
	for ptr >= 0 {
		if src[ptr] == 0x00 {
			if enc.SpecialByte > 0xE0 {
				ptr -= int(enc.SpecialByte & 0x1F)
			} else {
				ptr -= int(enc.SpecialByte)
			}
		} else {
			if src[ptr] > 0xE0 {
				ptr -= int(src[ptr] & 0x1F)
			} else {
				ptr -= int(src[ptr])
			}
		}
		numFlags++
	}
	return numFlags
}

func (enc R_PairelimEncoder) WorseCase(dLen int) (eLen int) {
	eLen = dLen + 1 + (dLen / 223)
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc R_PairelimEncoder) MaxOverhead(dLen int) (eLen int) {
	return enc.WorseCase(dLen)
}

func (enc R_PairelimEncoder) BestCase(dLen int) (eLen int) {
	eLen = (dLen / 2) + 1
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc R_PairelimEncoder) MinOverhead(dLen int) (eLen int) {
	return enc.BestCase(dLen)
}
