package cobs

type Type uint8

// The Following are a list of COBS-Types. Each has differents pros/cons listed below.
// Native: supports all verification, default protocol.
// Reduced: flag-based verification is inapplicable, potenially reduces overhead by 1 byte.
const (
	Native  Type = 0
	Reduced Type = 1
)

// Config is a struct that holds configuration variables on how
// this COBS library will function. It can be customized according
// to the use case.
type Config struct {
	SpecialByte byte
	Delimiter   bool
	Type        Type
}

// Encode takes raw data and a configuration and produces the COBS-encoded
// byte slice.
func Encode(src []byte, config Config) (dst []byte) {
	switch config.Type {
	case Native:
	  return nativeEncode(src, config)
	case Reduced:
	  return reducedEncode(src, config)
	default:
	  return nil
	}
}

// Decode takes encoded data and a configuration and produces the raw COBS-decoded
// byte slice.
func Decode(src []byte, config Config) (dst []byte) {
	switch config.Type {
	case Native:
	  return nativeDecode(src, config)
	case Reduced:
	  return reducedDecode(src, config)
	default:
	  return
	}
}

// Verify checks whether the given raw data can be a valid COBS-encoded byte slice
// based on the configuration. It can not only check to see if the special byte appears
// but also can see if the flags -lead- towards the end of the slice.
func Verify(src []byte, config Config) (err error) {
	switch config.Type {
	case Native:
	  return nativeVerify(src, config)
	case Reduced:
	  return reducedVerify(src, config)
	default:
	  return
	}
}

// WorseCase calculates the worse case for the COBS overhead when given
// a raw length and an appropiate configuration.
func WorseCase(dLen int, config Config) (eLen int) {
	eLen = dLen + 1 + (dLen / 254)
	if config.Delimiter {
		eLen++
	}
	return eLen
}

// MaxOverhead is an alias for WorseCase.
func MaxOverhead(dLen int, config Config) (eLen int) {
	return WorseCase(dLen, config)
}

// BestCase calculates the best case for the COBS overhead when given
// a raw length and an appropiate configuration.
func BestCase(dLen int, config Config) (eLen int) {
	eLen = dLen + 1
	if config.Delimiter {
		eLen++
	}
	if config.Type == Reduced {
		eLen--
	}
	return eLen
}

// MinOverhead is an alias for BestCase.
func MinOverhead(dLen int, config Config) (eLen int) {
	return BestCase(dLen, config)
}
