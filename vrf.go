package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
)

type VRFProof struct {
	Gamma      *big.Int
	GammaX     *big.Int
	GammaY     *big.Int
	C          *big.Int
	S          *big.Int
	PublicKey  *ecdsa.PublicKey
	InputHash  []byte
	OutputHash []byte
	Proof      []byte
}

func NewVRFManager() *VRFManager {
	return &VRFManager{}
}

type VRFManager struct {
	curve elliptic.Curve
}

func (v *VRFManager) GenerateVRFKeys() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func (v *VRFManager) Compute(privateKey *ecdsa.PrivateKey, seed []byte) (*VRFProof, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("nil private key")
	}

	h := sha256.Sum256(seed)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, h[:])
	if err != nil {
		return nil, err
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()
	proof := append(rBytes, sBytes...)

	outputData := append(h[:], proof...)
	outputHash := sha256.Sum256(outputData)

	return &VRFProof{
		OutputHash: outputHash[:],
		Proof:      proof,
	}, nil
}

func (v *VRFManager) Verify(publicKey *ecdsa.PublicKey, seed []byte, proof *VRFProof) (bool, error) {
	if publicKey == nil || proof == nil {
		return false, fmt.Errorf("nil public key or proof")
	}

	if len(proof.Proof) < 64 {
		return false, fmt.Errorf("invalid proof format")
	}

	halfLen := len(proof.Proof) / 2
	r := new(big.Int).SetBytes(proof.Proof[:halfLen])
	s := new(big.Int).SetBytes(proof.Proof[halfLen:])

	h := sha256.Sum256(seed)
	return ecdsa.Verify(publicKey, h[:], r, s), nil
}

func (v *VRFManager) ConvertOutputToRandomFloat(output []byte) float64 {
	if len(output) < 8 {
		return 0.5
	}

	value := binary.BigEndian.Uint64(output[:8])

	maxUint64 := float64(^uint64(0))
	return float64(value) / maxUint64
}
func (vm *VRFManager) GetOutput(proof *VRFProof) ([]byte, error) {
	if proof == nil {
		return nil, errors.New("proof cannot be nil")
	}
	return proof.OutputHash, nil
}

func (vm *VRFManager) hashToPoint(hash []byte) (*big.Int, *big.Int) {
	params := vm.curve.Params()
	counter := 0

	for {

		h := sha256.New()
		h.Write(hash)
		h.Write([]byte{byte(counter)})
		counterHash := h.Sum(nil)

		x := new(big.Int).SetBytes(counterHash)
		x.Mod(x, params.P)

		y2 := new(big.Int).Mul(x, x)
		y2.Mul(y2, x)
		y2.Add(y2, params.B)
		y2.Mod(y2, params.P)

		y := new(big.Int).ModSqrt(y2, params.P)

		if y != nil && vm.curve.IsOnCurve(x, y) {
			return x, y
		}

		counter++
	}
}

func (vm *VRFManager) computeOutputHash(gammaX, gammaY *big.Int) []byte {
	h := sha256.New()
	h.Write(gammaX.Bytes())
	h.Write(gammaY.Bytes())
	return h.Sum(nil)
}

func (vm *VRFManager) ConvertOutputToRandomInt(output []byte, max *big.Int) *big.Int {
	if max == nil || max.Sign() <= 0 {
		return big.NewInt(0)
	}

	randomValue := new(big.Int).SetBytes(output)

	return randomValue.Mod(randomValue, max)
}
