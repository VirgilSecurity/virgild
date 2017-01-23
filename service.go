package main

import (
	"encoding/hex"

	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilcrypto"
)

type ImpFingerprint struct {
	Crypto CalculateFingerprint
}

func (fp *ImpFingerprint) Calculate(data []byte) string {
	f := fp.Crypto.CalculateFingerprint(data)
	return hex.EncodeToString(f)
}

type Crypto interface {
	Sign(data []byte, signer virgilcrypto.PrivateKey) ([]byte, error)
	CalculateFingerprint(data []byte) []byte
}

type ImpRequestSigner struct {
	CardId     string
	PrivateKey virgilcrypto.PrivateKey
	Crypto     Crypto
}

func (s *ImpRequestSigner) Sign(req *virgil.SignableRequest) error {
	sign, err := s.Crypto.Sign(s.Crypto.CalculateFingerprint(req.Snapshot), s.PrivateKey)
	if err != nil {
		return err
	}
	req.AppendSignature(s.CardId, sign)
	return nil
}
