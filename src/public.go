package cobs

import (
	"errors"
)

type EncoderType int

const (
	Native                  EncoderType = 0
	ReversedNative          EncoderType = 1
	Reduced                 EncoderType = 2
	PairElimination         EncoderType = 3
	ReversedPairElimination EncoderType = 4
)

type Config struct {
	Type        EncoderType
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

type Encoder interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
	Verify([]byte) error
	FlagCount([]byte) int
	MaxOverhead(int) int
	MinOverhead(int) int
}

func NewEncoder(c Config) (Encoder, error) {
	switch c.Type {
	case Native:
		return nativeEncoder{
			SpecialByte: c.SpecialByte,
			Delimiter:   c.Delimiter,
			EndingSave:  c.EndingSave}, nil
	case ReversedNative:
		return reversedNativeEncoder{
			SpecialByte: c.SpecialByte,
			Delimiter:   c.Delimiter,
			EndingSave:  c.EndingSave}, nil
	case Reduced:
		return reducedEncoder{
			SpecialByte: c.SpecialByte,
			Delimiter:   c.Delimiter,
			EndingSave:  c.EndingSave}, nil
	case PairElimination:
		return pairElimEncoder{
			SpecialByte: c.SpecialByte,
			Delimiter:   c.Delimiter,
			EndingSave:  c.EndingSave}, nil
	case ReversedPairElimination:
		return reversedPairElimEncoder{
			SpecialByte: c.SpecialByte,
			Delimiter:   c.Delimiter,
			EndingSave:  c.EndingSave}, nil
	default:
		return nil, errors.New("Config Not Recongizable via Type.")
	}
}
