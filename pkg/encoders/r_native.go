package encoders

import (
	"errors"
)

type R_NativeEncoder struct {
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

func (enc R_NativeEncoder) Encode(src []byte) (dst []byte) {
	loopLen := len(src)
	dst = make([]byte, 0, loopLen+1)
	code := byte(0x01)
	ptr := 0
	for ptr < loopLen {
		if src[ptr] == enc.SpecialByte {
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
		if code == 0xFF && (!enc.EndingSave || ptr != loopLen-1) {
			if code == enc.SpecialByte {
				dst = append(dst, 0)
			} else {
				dst = append(dst, code)
			}
			code = 0x01
		}
		ptr++
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

func (enc R_NativeEncoder) Decode(src []byte) (dst []byte) {
	ptr := len(src)-1
	dst = make([]byte, 0, ptr-(ptr/254))
	if enc.Delimiter {
		ptr--
	}
	code := byte(0x00)
	for ptr >= 0 {
		if src[ptr] == 0x00 {
			code = enc.SpecialByte
		} else {
			code = src[ptr]
		}
		if code < 0xFF || (enc.EndingSave && ptr < 0) {
			dst = append([]byte{enc.SpecialByte}, dst...)
		}
		ptr--
		for i := 1; i < int(code); i++ {
			dst = append([]byte{src[ptr]}, dst...)
			ptr--
		}
	}
	return dst[:len(dst)-1]
}

func (enc R_NativeEncoder) Verify(src []byte) (err error) {
	nextFlag := 0
	ptr := len(src)
	if enc.Delimiter {
		if ptr < 2 {
			return errors.New("COBS[Native]: Encoded slice is too short to be valid.")
		}
		if src[ptr-1] != enc.SpecialByte {
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
		if src[ptr] == enc.SpecialByte {
			return errors.New("COBS[Native]: Encoded slice's byte (not the delimter) is special byte.")
		}
		if nextFlag == 0 {
			if src[ptr] == 0 {
				nextFlag = int(enc.SpecialByte)
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

func (enc R_NativeEncoder) FlagCount(src []byte) (flags int) {
	numFlags := 0
	ptr := len(src)-1
	if enc.Delimiter {
		ptr--
		numFlags++
	}
	for ptr >= 0 {
		if src[ptr] == 0 {
			ptr -= int(enc.SpecialByte)
		} else {
			ptr -= int(src[ptr])
		}
		numFlags++
	}
	return numFlags
}

func (enc R_NativeEncoder) WorseCase(dLen int) (eLen int) {
	eLen = dLen + 1 + (dLen / 254)
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc R_NativeEncoder) MaxOverhead(dLen int) (eLen int) {
	return enc.WorseCase(dLen)
}

func (enc R_NativeEncoder) BestCase(dLen int) (eLen int) {
	eLen = dLen + 1
	if enc.Delimiter {
		eLen++
	}
	return eLen
}

func (enc R_NativeEncoder) MinOverhead(dLen int) (eLen int) {
	return enc.BestCase(dLen)
}