package cobs

import (
	"errors"
)

type pairElimEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (e pairElimEncoder) Encode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 1, loopLen+1)
	codePtr := 0
	code := byte(0x01)
	pairable := false
	ptr := 0
	for ptr < loopLen {
		if pairable {
			pairable = false
			if src[ptr] == e.SpecialByte {
				code |= 0xE0
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
			if code == e.SpecialByte {
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
		if src[ptr] == e.SpecialByte {
			if code <= 0x1F {
				pairable = true
				ptr++
				continue
			}
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
		if code == 0xE0 && (!e.EndingSave || ptr != loopLen-1) {
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
	if pairable {
		code |= 0xE0
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

func (e pairElimEncoder) Decode(src []byte) []byte {
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
			dst = append(dst, e.SpecialByte)
			dst = append(dst, e.SpecialByte)
		} else if code < 0xE0 || (e.EndingSave && ptr == loopLen) {
			dst = append(dst, e.SpecialByte)
		}
	}
	return dst[:len(dst)-1]
}

func (e pairElimEncoder) Verify(src []byte) error {
	nextFlag := 0
	loopLen := len(src)
	if e.Delimiter {
		if loopLen < 2 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
		if src[loopLen-1] != e.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's delimiter is not special byte.")
		}
		loopLen--
	} else {
		if loopLen < 1 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
	}
	for _, b := range src[:loopLen] {
		if b == e.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if b == 0x00 {
				if e.SpecialByte > 0xE0 {
					nextFlag = int(e.SpecialByte & 0x1F)
				} else {
					nextFlag = int(e.SpecialByte)
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

func (e pairElimEncoder) FlagCount(src []byte) int {
	numFlags := 0
	ptr := 0
	loopLen := len(src)
	if e.Delimiter {
		numFlags++
		loopLen--
	}
	for ptr < loopLen {
		if src[ptr] == 0x00 {
			if e.SpecialByte > 0xE0 {
				ptr += int(e.SpecialByte & 0x1F)
			} else {
				ptr += int(e.SpecialByte)
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

func (e pairElimEncoder) MaxOverhead(len int) int {
	ret := len + 1 + (len / 223)
	if e.Delimiter {
		ret++
	}
	return ret
}

func (e pairElimEncoder) MinOverhead(len int) int {
	ret := (len / 2) + 1
	if e.Delimiter {
		ret++
	}
	return ret
}
