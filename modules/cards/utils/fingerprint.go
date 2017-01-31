package utils

import "encoding/hex"

type CalculateFingerprint interface {
	CalculateFingerprint(data []byte) []byte
}

type Fingerprint struct {
	Crypto CalculateFingerprint
}

func (fp *Fingerprint) Calculate(data []byte) string {
	f := fp.Crypto.CalculateFingerprint(data)
	return hex.EncodeToString(f)
}
