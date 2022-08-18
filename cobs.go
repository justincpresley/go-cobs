package cobs

import (
	"errors"
	enc "github.com/justincpresley/go-cobs/encoders"
)

type Type uint8

// The following are a list of COBS-Types. Each has different pros/cons.
// To learn more about types, see USE.md on the github repository.
const (
	Native          Type = 0
	Reduced         Type = 1
	PairElimination Type = 2
)

// Config is a struct that holds configuration variables on how
// this COBS library will function. It can be customized according
// to the use case. To learn more about the config parameters and what they do,
// see USE.md on the github repository.
type Config struct {
	Type        Type
	Reverse     bool
	SpecialByte byte
	Delimiter   bool
	EndingSave  bool
}

// Encoder is a interface outlining all avaliable functions.
type Encoder interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
	Verify([]byte) error
	FlagCount([]byte) int
	MaxOverhead(int) int
	MinOverhead(int) int
	BestCase(int) int
	WorseCase(int) int
}

// NewEncoder creates a specialized Encoder based on the given Config.
func NewEncoder(c Config) (Encoder, error) {
	switch c.Type {
	case Native:
		if !c.Reverse {
			return enc.NativeEncoder{
				SpecialByte: c.SpecialByte,
				Delimiter:   c.Delimiter,
				EndingSave:  c.EndingSave}, nil
		} else {
			return enc.R_NativeEncoder{
				SpecialByte: c.SpecialByte,
				Delimiter:   c.Delimiter,
				EndingSave:  c.EndingSave}, nil
		}
	case Reduced:
		if !c.Reverse {
			return enc.ReducedEncoder{
				SpecialByte: c.SpecialByte,
				Delimiter:   c.Delimiter,
				EndingSave:  c.EndingSave}, nil
		} else {
			return nil, errors.New("Reverse not avaliable for Reduced yet.")
		}
	case PairElimination:
		if !c.Reverse {
			return enc.PairelimEncoder{
				SpecialByte: c.SpecialByte,
				Delimiter:   c.Delimiter,
				EndingSave:  c.EndingSave}, nil
		} else {
			return enc.R_PairelimEncoder{
				SpecialByte: c.SpecialByte,
				Delimiter:   c.Delimiter,
				EndingSave:  c.EndingSave}, nil
		}
	default:
		return nil, errors.New("Config Not Recongizable via Type.")
	}
}
