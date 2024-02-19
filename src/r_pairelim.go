package cobs

import (
	"errors"
)

type reversedPairElimEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (e reversedPairElimEncoder) Encode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 0, loopLen+1)
	code := byte(0x01)
	pairable := false
	ptr := 0
	for ptr < loopLen {
		if pairable {
			pairable = false
			if src[ptr] == e.SpecialByte {
				code |= 0xE0
				if code == e.SpecialByte {
					dst = append(dst, 0)
				} else {
					dst = append(dst, code)
				}
				code = 0x01
				ptr++
				continue
			}
			if code == e.SpecialByte {
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
		if src[ptr] == e.SpecialByte {
			if code <= 0x1F {
				pairable = true
				ptr++
				continue
			}
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
		if code == 0xE0 && (!e.EndingSave || ptr != loopLen-1) {
			if code == e.SpecialByte {
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

func (e reversedPairElimEncoder) Decode(src []byte) []byte {
	loopLen := len(src)
	dst := make([]byte, 0, loopLen-1-(loopLen/254))
	if e.Delimiter {
		loopLen--
	}
	ptr := loopLen - 1
	jumpLen := 0
	code := byte(0x00)
	for ptr >= 0 {
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
		if code > 0xE0 {
			dst = append([]byte{e.SpecialByte}, dst...)
			dst = append([]byte{e.SpecialByte}, dst...)
		} else if code < 0xE0 || (e.EndingSave && ptr == loopLen) {
			dst = append([]byte{e.SpecialByte}, dst...)
		}
		ptr--
		for i := 1; i < jumpLen; i++ {
			dst = append([]byte{src[ptr]}, dst...)
			ptr--
		}
	}
	return dst[:len(dst)-1]
}

func (e reversedPairElimEncoder) Verify(src []byte) error {
	nextFlag := 0
	ptr := len(src)
	if e.Delimiter {
		if ptr < 2 {
			return errors.New("COBS[PairElimination]: Encoded slice is too short to be valid.")
		}
		if src[ptr-1] != e.SpecialByte {
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
		if src[ptr] == e.SpecialByte {
			return errors.New("COBS[PairElimination]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if src[ptr] == 0x00 {
				if e.SpecialByte > 0xE0 {
					nextFlag = int(e.SpecialByte & 0x1F)
				} else {
					nextFlag = int(e.SpecialByte)
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

func (e reversedPairElimEncoder) FlagCount(src []byte) int {
	numFlags := 0
	ptr := len(src) - 1
	if e.Delimiter {
		ptr--
		numFlags++
	}
	for ptr >= 0 {
		if src[ptr] == 0x00 {
			if e.SpecialByte > 0xE0 {
				ptr -= int(e.SpecialByte & 0x1F)
			} else {
				ptr -= int(e.SpecialByte)
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

func (e reversedPairElimEncoder) MaxOverhead(len int) int {
	ret := len + 1 + (len / 223)
	if e.Delimiter {
		ret++
	}
	return ret
}

func (e reversedPairElimEncoder) MinOverhead(len int) int {
	ret := (len / 2) + 1
	if e.Delimiter {
		ret++
	}
	return ret
}
