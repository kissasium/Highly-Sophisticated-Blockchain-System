package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

// VRFProof contains the proof data that can be used to verify
// that the output was correctly computed
type VRFProof struct {
	Gamma      *big.Int // ECDSA point representing the VRF output
	GammaX     *big.Int // X coordinate of Gamma
	GammaY     *big.Int // Y coordinate of Gamma
	C          *big.Int // Challenge value for the VRF proof
	S          *big.Int // Response value for the VRF proof
	PublicKey  *ecdsa.PublicKey
	InputHash  []byte
	OutputHash []byte
}

// VRFManager handles VRF operations for the blockchain system
type VRFManager struct {
	curve elliptic.Curve // Elliptic curve to use for VRF operations
}

// NewVRFManager creates a new VRF manager using the specified curve
func NewVRFManager() *VRFManager {
	return &VRFManager{
		curve: elliptic.P256(), // Using P-256 curve, same as in your ZKP system
	}
}

// GenerateVRFKeys generates a new key pair for VRF operations
func (vm *VRFManager) GenerateVRFKeys() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(vm.curve, rand.Reader)
}

// Compute calculates a VRF value and proof from the given input using the secret key
func (vm *VRFManager) Compute(sk *ecdsa.PrivateKey, input []byte) (*VRFProof, error) {
	if sk == nil {
		return nil, errors.New("private key cannot be nil")
	}

	// Hash the input to ensure uniform distribution
	h := sha256.New()
	h.Write(input)
	inputHash := h.Sum(nil)

	// 1. Compute ECDSA point Gamma = sk * H(input), where H maps input to a curve point
	hx, hy := vm.hashToPoint(inputHash)
	gammaX, gammaY := vm.curve.ScalarMult(hx, hy, sk.D.Bytes())

	// 2. Compute challenge value c = H(gamma, input)
	cHash := sha256.New()
	cHash.Write(gammaX.Bytes())
	cHash.Write(gammaY.Bytes())
	cHash.Write(inputHash)
	cBytes := cHash.Sum(nil)
	c := new(big.Int).SetBytes(cBytes)
	c.Mod(c, vm.curve.Params().N)

	// 3. Compute response value s = k - c * sk mod n
	// First, generate a random k (nonce)
	k, err := rand.Int(rand.Reader, vm.curve.Params().N)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random nonce: %v", err)
	}

	// Calculate s = k - c * sk mod n
	s := new(big.Int).Mul(c, sk.D)
	s.Sub(k, s)
	s.Mod(s, vm.curve.Params().N)

	// 4. Generate output hash (VRF value) from gamma
	outputHash := vm.computeOutputHash(gammaX, gammaY)

	// Create and return the proof
	return &VRFProof{
		Gamma:      new(big.Int).SetBytes(append(gammaX.Bytes(), gammaY.Bytes()...)),
		GammaX:     gammaX,
		GammaY:     gammaY,
		C:          c,
		S:          s,
		PublicKey:  &sk.PublicKey,
		InputHash:  inputHash,
		OutputHash: outputHash,
	}, nil
}

// Verify checks if a VRF proof is valid for the given input and public key
func (vm *VRFManager) Verify(pk *ecdsa.PublicKey, input []byte, proof *VRFProof) (bool, error) {
	if pk == nil || proof == nil {
		return false, errors.New("public key and proof cannot be nil")
	}

	// Hash the input
	h := sha256.New()
	h.Write(input)
	inputHash := h.Sum(nil)

	// Verify that proof.InputHash matches the hash of the input
	if hex.EncodeToString(inputHash) != hex.EncodeToString(proof.InputHash) {
		return false, errors.New("input hash mismatch")
	}

	// 1. Compute the hash point H(input)
	// hx, hy := vm.hashToPoint(inputHash)
	_, _ = vm.hashToPoint(inputHash)

	// 2. Verify the ECDSA challenge-response
	// u1 = s * G
	u1x, u1y := vm.curve.ScalarBaseMult(proof.S.Bytes())

	// u2 = c * Y (where Y is the public key)
	u2x, u2y := vm.curve.ScalarMult(pk.X, pk.Y, proof.C.Bytes())

	// gamma' = u1 + u2
	gammaX, gammaY := vm.curve.Add(u1x, u1y, u2x, u2y)

	// 3. Recompute challenge c' = H(gamma', input)
	cHash := sha256.New()
	cHash.Write(gammaX.Bytes())
	cHash.Write(gammaY.Bytes())
	cHash.Write(inputHash)
	cBytes := cHash.Sum(nil)
	c := new(big.Int).SetBytes(cBytes)
	c.Mod(c, vm.curve.Params().N)

	// 4. Check if the computed challenge matches the one in the proof
	if c.Cmp(proof.C) != 0 {
		return false, errors.New("challenge verification failed")
	}

	// 5. Verify the output hash
	expectedOutputHash := vm.computeOutputHash(gammaX, gammaY)
	if hex.EncodeToString(expectedOutputHash) != hex.EncodeToString(proof.OutputHash) {
		return false, errors.New("output verification failed")
	}

	return true, nil
}

// GetOutput retrieves the VRF output value from a proof
func (vm *VRFManager) GetOutput(proof *VRFProof) ([]byte, error) {
	if proof == nil {
		return nil, errors.New("proof cannot be nil")
	}
	return proof.OutputHash, nil
}

// // hashToPoint converts a hash to a point on the elliptic curve
// // This is a simplified implementation - in production, use a proper hash-to-curve function
// func (vm *VRFManager) hashToPoint(hash []byte) (*big.Int, *big.Int) {
// 	params := vm.curve.Params()

// 	// Try incrementing counter until we get a valid point
// 	counter := 0
// 	for {
// 		// Combine hash with counter
// 		h := sha256.New()
// 		h.Write(hash)
// 		h.Write([]byte{byte(counter)})
// 		counterHash := h.Sum(nil)

// 		// Convert to a potential x-coordinate
// 		x := new(big.Int).SetBytes(counterHash)
// 		x.Mod(x, params.P)

// 		// Attempt to find corresponding y-coordinate
// 		y2 := new(big.Int).Mul(x, x)
// 		y2.Mul(y2, x)
// 		y2.Add(y2, params.B)
// 		y2.Mod(y2, params.P)

// 		y := new(big.Int).ModSqrt(y2, params.P)
// 		if y != nil {
// 			// Found a valid point
// 			return x, y
// 		}

// 		counter++
// 	}
// }

func (vm *VRFManager) hashToPoint(hash []byte) (*big.Int, *big.Int) {
	params := vm.curve.Params()
	counter := 0

	for {
		// Combine hash with counter
		h := sha256.New()
		h.Write(hash)
		h.Write([]byte{byte(counter)})
		counterHash := h.Sum(nil)

		x := new(big.Int).SetBytes(counterHash)
		x.Mod(x, params.P)

		// Calculate y² = x³ + B mod p
		y2 := new(big.Int).Mul(x, x)
		y2.Mul(y2, x)
		y2.Add(y2, params.B)
		y2.Mod(y2, params.P)

		y := new(big.Int).ModSqrt(y2, params.P)

		// ✅ Extra safety check
		if y != nil && vm.curve.IsOnCurve(x, y) {
			return x, y
		}

		counter++
	}
}

// computeOutputHash generates the final VRF output value from gamma
func (vm *VRFManager) computeOutputHash(gammaX, gammaY *big.Int) []byte {
	h := sha256.New()
	h.Write(gammaX.Bytes())
	h.Write(gammaY.Bytes())
	return h.Sum(nil)
}

// ConvertOutputToRandomInt converts a VRF output to a random integer in [0, max)
func (vm *VRFManager) ConvertOutputToRandomInt(output []byte, max *big.Int) *big.Int {
	if max == nil || max.Sign() <= 0 {
		return big.NewInt(0)
	}

	// Convert output to big.Int
	randomValue := new(big.Int).SetBytes(output)

	// Mod by max to get a value in the range [0, max)
	return randomValue.Mod(randomValue, max)
}

// ConvertOutputToRandomFloat converts a VRF output to a random float in [0, 1)
func (vm *VRFManager) ConvertOutputToRandomFloat(output []byte) float64 {
	// Convert output to big.Int
	randomValue := new(big.Int).SetBytes(output)

	// Ensure positive value
	randomValue.Abs(randomValue)

	// Divide by 2^256 to get a value in [0, 1)
	maxVal := new(big.Int).Lsh(big.NewInt(1), 256)
	ratio := new(big.Float).Quo(
		new(big.Float).SetInt(randomValue),
		new(big.Float).SetInt(maxVal),
	)

	// Convert to float64
	result, _ := ratio.Float64()
	return result
}
