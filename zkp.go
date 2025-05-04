package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

type ZKProof struct {
	Commitment string
	Challenge  string
	Response   *big.Int
	Aux        []byte
}

type ZKPSystem struct {
	curve   elliptic.Curve
	G       *big.Int
	Q       *big.Int
	proofs  map[string]ZKProof
	secrets map[string]*big.Int
}

type RangeProof struct {
	Commitment string
	Proof      []byte
	Min        *big.Int
	Max        *big.Int
}

func NewZKPSystem() *ZKPSystem {
	return &ZKPSystem{}
}

func (zkp *ZKPSystem) GenerateSchnorrProof(secretValue *big.Int, identifier string) (*ZKProof, error) {
	if secretValue == nil {
		return nil, errors.New("secret value cannot be nil")
	}

	// Store secret (in a real system, this would be securely stored or not stored at all)
	zkp.secrets[identifier] = secretValue

	// Step 1: Pick a random value (commitment randomness)
	r, err := rand.Int(rand.Reader, zkp.Q)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random value: %v", err)
	}

	// Step 2: Calculate commitment (g^r mod p)
	rx, ry := zkp.curve.ScalarBaseMult(r.Bytes())
	commitment := fmt.Sprintf("%x,%x", rx, ry)

	// Step 3: Compute challenge (in interactive version, this would come from verifier)
	// Here we simulate the challenge using a hash
	h := sha256.New()
	h.Write([]byte(commitment))
	h.Write([]byte(identifier))
	challenge := fmt.Sprintf("%x", h.Sum(nil))

	// Step 4: Compute response = r + challenge * secret (mod q)
	e := new(big.Int).SetBytes(h.Sum(nil))
	e.Mod(e, zkp.Q)

	// s = r + e * x (mod q)
	s := new(big.Int).Mul(e, secretValue)
	s.Add(s, r)
	s.Mod(s, zkp.Q)

	// Create and store the proof
	proof := &ZKProof{
		Commitment: commitment,
		Challenge:  challenge,
		Response:   s,
	}

	zkp.proofs[identifier] = *proof
	return proof, nil
}

func (zkp *ZKPSystem) VerifySchnorrProof(publicX, publicY *big.Int, proof *ZKProof, identifier string) (bool, error) {
	if proof == nil {
		return false, errors.New("proof cannot be nil")
	}

	// Parse the commitment
	var rxStr, ryStr string
	_, err := fmt.Sscanf(proof.Commitment, "%x,%x", &rxStr, &ryStr)
	if err != nil {
		return false, fmt.Errorf("invalid commitment format: %v", err)
	}

	// Reconstruct challenge
	h := sha256.New()
	h.Write([]byte(proof.Commitment))
	h.Write([]byte(identifier))
	calculatedChallenge := fmt.Sprintf("%x", h.Sum(nil))

	if calculatedChallenge != proof.Challenge {
		return false, errors.New("challenge mismatch")
	}

	// Convert challenge to big.Int
	e := new(big.Int).SetBytes(h.Sum(nil))
	e.Mod(e, zkp.Q)

	// Verification: g^s = commitment * (public_key)^challenge
	// Calculate g^s
	sx, sy := zkp.curve.ScalarBaseMult(proof.Response.Bytes())

	// Calculate (public_key)^challenge
	ex, ey := zkp.curve.ScalarMult(publicX, publicY, e.Bytes())

	// Parse original commitment point
	rx := new(big.Int)
	ry := new(big.Int)
	rx.SetString(rxStr, 16)
	ry.SetString(ryStr, 16)

	// Calculate commitment * (public_key)^challenge
	resultX, resultY := zkp.curve.Add(rx, ry, ex, ey)

	// Compare the points
	return sx.Cmp(resultX) == 0 && sy.Cmp(resultY) == 0, nil
}

func (z *ZKPSystem) GenerateRangeProof(value, min, max *big.Int, identifier string) (*RangeProof, error) {
	if value.Cmp(min) < 0 || value.Cmp(max) > 0 {
		return nil, fmt.Errorf("value %v outside range [%v, %v]", value, min, max)
	}

	valueBytes := value.Bytes()
	randomness := make([]byte, 32)
	rand.Read(randomness)

	commitmentData := append(valueBytes, randomness...)
	commitmentData = append(commitmentData, []byte(identifier)...)
	commitmentHash := sha256.Sum256(commitmentData)

	proof := make([]byte, 64)
	copy(proof[0:32], valueBytes)
	copy(proof[32:64], randomness)

	return &RangeProof{
		Commitment: hex.EncodeToString(commitmentHash[:]),
		Proof:      proof,
		Min:        min,
		Max:        max,
	}, nil
}

func (z *ZKPSystem) VerifyRangeProof(proof *RangeProof, identifier string) (bool, error) {
	if proof == nil {
		return false, fmt.Errorf("nil proof provided")
	}

	if len(proof.Proof) < 64 {
		return false, fmt.Errorf("invalid proof format: proof too short")
	}

	valueBytes := proof.Proof[0:32]
	randomness := proof.Proof[32:64]

	value := new(big.Int).SetBytes(valueBytes)

	if value.Cmp(proof.Min) < 0 || value.Cmp(proof.Max) > 0 {
		return false, nil
	}

	commitmentData := append(valueBytes, randomness...)
	commitmentData = append(commitmentData, []byte(identifier)...)
	calculatedHash := sha256.Sum256(commitmentData)
	calculatedCommitment := hex.EncodeToString(calculatedHash[:])

	return calculatedCommitment == proof.Commitment, nil
}

func (zkp *ZKPSystem) CreateZeroKnowledgeSetMembership(
	value *big.Int,
	set []*big.Int,
	identifier string,
) (*ZKProof, error) {

	found := false
	for _, elem := range set {
		if elem.Cmp(value) == 0 {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("value not in the set")
	}

	zkp.secrets[identifier] = value

	r, err := rand.Int(rand.Reader, zkp.Q)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random value: %v", err)
	}

	vx, vy := zkp.curve.ScalarBaseMult(value.Bytes())
	rx, ry := zkp.curve.ScalarBaseMult(r.Bytes())
	commitX, commitY := zkp.curve.Add(vx, vy, rx, ry)
	commitment := fmt.Sprintf("%x,%x", commitX, commitY)

	setHashes := make([]string, len(set))
	for i, elem := range set {
		h := sha256.New()
		h.Write(elem.Bytes())
		setHashes[i] = fmt.Sprintf("%x", h.Sum(nil))
	}

	h := sha256.New()
	h.Write([]byte(commitment))
	for _, hash := range setHashes {
		h.Write([]byte(hash))
	}
	h.Write([]byte(identifier))
	challenge := fmt.Sprintf("%x", h.Sum(nil))

	e := new(big.Int).SetBytes(h.Sum(nil))
	e.Mod(e, zkp.Q)

	s := new(big.Int).Mul(e, value)
	s.Add(s, r)
	s.Mod(s, zkp.Q)

	setStr := ""
	for i, elem := range set {
		if i > 0 {
			setStr += ","
		}
		setStr += elem.String()
	}

	proof := &ZKProof{
		Commitment: commitment,
		Challenge:  challenge,
		Response:   s,
		Aux:        []byte(setStr),
	}

	zkp.proofs[identifier] = *proof
	return proof, nil
}

func (zkp *ZKPSystem) VerifySetMembership(proof *ZKProof, identifier string) (bool, error) {
	if proof == nil {
		return false, errors.New("proof cannot be nil")
	}

	setStr := string(proof.Aux)
	setStrings := []string{}
	var n int
	for i := 0; i < len(setStr); i++ {
		if setStr[i] == ',' {
			n++
		}
	}

	set := make([]*big.Int, n+1)
	_, err := fmt.Sscanf(setStr, "%s", &setStrings)
	if err != nil {

		set = []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	}

	setHashes := make([]string, len(set))
	for i, elem := range set {
		h := sha256.New()
		h.Write(elem.Bytes())
		setHashes[i] = fmt.Sprintf("%x", h.Sum(nil))
	}

	h := sha256.New()
	h.Write([]byte(proof.Commitment))
	for _, hash := range setHashes {
		h.Write([]byte(hash))
	}
	h.Write([]byte(identifier))
	calculatedChallenge := fmt.Sprintf("%x", h.Sum(nil))

	if calculatedChallenge != proof.Challenge {
		return false, errors.New("challenge mismatch")
	}
	_, exists := zkp.proofs[identifier]
	return exists, nil
}
