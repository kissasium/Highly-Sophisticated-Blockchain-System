package main

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"sync"
	"time"
)

// Node represents a participant in the blockchain network with advanced authentication
type Node struct {
	ID                 string            // Unique identifier for the node
	PublicKey          *ecdsa.PublicKey  // Primary authentication credential
	PrivateKey         *ecdsa.PrivateKey // Private key (only stored locally)
	IP                 net.IP            // Network address
	Port               int               // Network port
	Role               string            // Role in the network (validator, observer, etc.)
	JoinTimestamp      time.Time         // When the node joined the network
	LastActiveTime     time.Time         // Last recorded activity
	ChallengeNonce     string            // Current challenge nonce for continuous auth
	TrustScore         float64           // Dynamic trust score (0.0 to 1.0)
	TrustScoreHistory  []TrustScoreEntry // Historical trust scores
	AuthFactors        []AuthFactor      // Multiple authentication factors
	ValidatorPower     int               // Voting power (for validators)
	BehaviorMetrics    BehaviorMetrics   // Metrics on node behavior
	ConnectionMetrics  ConnectionMetrics // Network connection metrics
	Reputation         ReputationScore   // Comprehensive reputation metrics
	AuthStateValid     bool              // Whether authentication is currently valid
	ContinuousAuthLock sync.Mutex        // Lock for concurrent authentication updates
}

// TrustScoreEntry records changes to a node's trust score over time
type TrustScoreEntry struct {
	Timestamp time.Time
	Score     float64
	Reason    string
}

// AuthFactor represents a single authentication factor
type AuthFactor struct {
	Type       string // Type of factor (e.g., "pubkey", "hmac", "geolocation", "behavior")
	Value      []byte // Value of the factor (encoding depends on Type)
	LastVerify time.Time
	IsVerified bool
}

// BehaviorMetrics tracks the behavior of a node over time
type BehaviorMetrics struct {
	ProposedBlocks         int     // Number of blocks proposed
	ValidBlocks            int     // Number of valid blocks
	InvalidBlocks          int     // Number of invalid blocks
	TotalVotes             int     // Total number of votes cast
	CorrectVotes           int     // Number of correct votes (matched consensus)
	ResponseTimeAvg        float64 // Average time to respond to requests (ms)
	UpTime                 float64 // Percentage of time the node is available
	AnomalyScore           float64 // Score indicating abnormal behavior (0.0 to 1.0)
	ConsecutiveValidations int     // Number of consecutive successful validations
}

// ConnectionMetrics tracks network connection patterns
type ConnectionMetrics struct {
	LastSeenIPs         []net.IP      // Recent IP addresses
	GeolocationsHistory []Geolocation // History of geolocations
	ConnectionTimes     []time.Time   // Connection timestamps
	PingTimes           []int         // Recent ping times (ms)
	PacketLoss          float64       // Recent packet loss percentage
	ConnectionCount     int           // Number of connections established
	DisconnectCount     int           // Number of unexpected disconnects
}

// Geolocation represents a physical location estimate for a node
type Geolocation struct {
	Latitude  float64
	Longitude float64
	Country   string
	Region    string
	Timestamp time.Time
}

// ReputationScore combines various metrics into a comprehensive score
type ReputationScore struct {
	Overall            float64 // Combined reputation score (0.0 to 1.0)
	Behavioral         float64 // Based on behavior metrics
	Network            float64 // Based on network metrics
	Consensus          float64 // Based on consensus participation
	Age                float64 // Based on account age/history
	StakeWeight        float64 // Based on stake in the network
	CommunityTrust     float64 // Based on trust from other reputable nodes
	LastCalculated     time.Time
	CalculationVersion int // Version of calculation algorithm
}

// NewNode creates a new node with default values and initializes authentication
func NewNode(id string, pubKey *ecdsa.PublicKey, privKey *ecdsa.PrivateKey, role string) (*Node, error) {
	if pubKey == nil || (role == "validator" && privKey == nil) {
		return nil, errors.New("invalid credentials for node creation")
	}

	// Generate initial nonce for authentication challenges
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}
	nonce := hex.EncodeToString(nonceBytes)

	node := &Node{
		ID:             id,
		PublicKey:      pubKey,
		PrivateKey:     privKey,
		Role:           role,
		JoinTimestamp:  time.Now(),
		LastActiveTime: time.Now(),
		ChallengeNonce: nonce,
		TrustScore:     0.5, // Start with neutral trust
		AuthFactors: []AuthFactor{
			{
				Type:       "pubkey",
				Value:      []byte(fmt.Sprintf("%v", pubKey)),
				LastVerify: time.Now(),
				IsVerified: true,
			},
		},
		ValidatorPower: 1,
		BehaviorMetrics: BehaviorMetrics{
			UpTime: 100.0, // Start at 100%
		},
		Reputation: ReputationScore{
			Overall:            0.5,
			Behavioral:         0.5,
			Network:            0.5,
			Consensus:          0.5,
			Age:                0.1, // New node starts low
			StakeWeight:        0.5,
			CommunityTrust:     0.5,
			LastCalculated:     time.Now(),
			CalculationVersion: 1,
		},
		AuthStateValid: true,
	}

	return node, nil
}

// GenerateAuthChallenge creates a new challenge for continuous authentication
func (n *Node) GenerateAuthChallenge() (string, error) {
	challengeBytes := make([]byte, 32)
	if _, err := rand.Read(challengeBytes); err != nil {
		return "", fmt.Errorf("failed to generate auth challenge: %v", err)
	}

	n.ContinuousAuthLock.Lock()
	defer n.ContinuousAuthLock.Unlock()

	challenge := hex.EncodeToString(challengeBytes)
	n.ChallengeNonce = challenge
	return challenge, nil
}

// VerifyAuthResponse verifies the response to an authentication challenge
func (n *Node) VerifyAuthResponse(challenge string, signature []byte) bool {
	if challenge != n.ChallengeNonce {
		return false
	}

	// In a real implementation, would verify ECDSA signature
	// For example:
	// hash := sha256.Sum256([]byte(challenge))
	// var r, s big.Int
	// r.SetBytes(signature[:32])
	// s.SetBytes(signature[32:])
	// return ecdsa.Verify(n.PublicKey, hash[:], &r, &s)

	// Placeholder for signature verification
	return len(signature) > 0
}

// AddAuthFactor adds a new authentication factor for the node
func (n *Node) AddAuthFactor(factorType string, value []byte) error {
	for _, factor := range n.AuthFactors {
		if factor.Type == factorType {
			return fmt.Errorf("auth factor type '%s' already exists", factorType)
		}
	}

	n.AuthFactors = append(n.AuthFactors, AuthFactor{
		Type:       factorType,
		Value:      value,
		LastVerify: time.Now(),
		IsVerified: true,
	})

	return nil
}

// VerifyMultiFactorAuth verifies all available authentication factors
func (n *Node) VerifyMultiFactorAuth(factorResponses map[string][]byte) (bool, error) {
	if len(factorResponses) == 0 {
		return false, errors.New("no authentication factors provided")
	}

	validFactors := 0
	requiredFactors := len(n.AuthFactors)

	for _, factor := range n.AuthFactors {
		response, exists := factorResponses[factor.Type]
		if !exists {
			continue
		}

		valid := false

		switch factor.Type {
		case "pubkey":
			// Verify ECDSA signature (simplified)
			valid = len(response) > 0
		case "hmac":
			// Verify HMAC
			expectedHMAC := factor.Value
			valid = subtle.ConstantTimeCompare(expectedHMAC, response) == 1
		case "geolocation":
			// Verify if geolocation is within acceptable range
			valid = verifyGeolocation(factor.Value, response)
		case "behavior":
			// Verify behavioral patterns
			valid = verifyBehaviorPattern(factor.Value, response)
		}

		if valid {
			validFactors++
		}
	}

	// Require at least 2 factors or all available factors if less than 2
	requiredValid := int(math.Min(float64(requiredFactors), 2))
	return validFactors >= requiredValid, nil
}

// verifyGeolocation checks if a geolocation is within acceptable parameters
func verifyGeolocation(expected, actual []byte) bool {
	// Simplified implementation - would parse coordinates and check distance
	// This is a placeholder for actual geo-verification logic
	return true
}

// verifyBehaviorPattern checks if behavior matches expected patterns
func verifyBehaviorPattern(expected, actual []byte) bool {
	// Simplified implementation - would compare behavioral metrics
	// This is a placeholder for actual behavioral verification
	return true
}

// UpdateTrustScore adjusts a node's trust score based on behavior
func (n *Node) UpdateTrustScore(delta float64, reason string) {
	n.TrustScore += delta

	// Ensure trust score stays between 0 and 1
	if n.TrustScore < 0 {
		n.TrustScore = 0
	} else if n.TrustScore > 1 {
		n.TrustScore = 1
	}

	// Record the change
	n.TrustScoreHistory = append(n.TrustScoreHistory, TrustScoreEntry{
		Timestamp: time.Now(),
		Score:     n.TrustScore,
		Reason:    reason,
	})
}

// RecalculateReputation performs a comprehensive recalculation of node reputation
func (n *Node) RecalculateReputation() {
	// Calculate behavioral score
	behavioralMetrics := n.BehaviorMetrics

	// Behavior score is influenced by valid vs invalid blocks and votes
	var behaviorScore float64 = 0.5 // Default neutral
	if behavioralMetrics.TotalVotes > 0 {
		correctVoteRatio := float64(behavioralMetrics.CorrectVotes) / float64(behavioralMetrics.TotalVotes)
		behaviorScore = correctVoteRatio
	}

	// Adjust for block validation
	totalBlocks := behavioralMetrics.ValidBlocks + behavioralMetrics.InvalidBlocks
	if totalBlocks > 0 {
		validBlockRatio := float64(behavioralMetrics.ValidBlocks) / float64(totalBlocks)
		behaviorScore = (behaviorScore + validBlockRatio) / 2
	}

	// Network score calculation
	networkScore := 0.5 // Default
	if len(n.ConnectionMetrics.PingTimes) > 0 {
		// Lower ping times are better
		avgPing := 0
		for _, ping := range n.ConnectionMetrics.PingTimes {
			avgPing += ping
		}
		avgPing /= len(n.ConnectionMetrics.PingTimes)

		// Scale ping score - lower is better (e.g., 10ms = excellent, 500ms = poor)
		pingScore := math.Max(0, 1-(float64(avgPing)/500))

		// Packet loss penalty - 0% loss = 1.0, 100% loss = 0.0
		packetScore := 1 - n.ConnectionMetrics.PacketLoss

		networkScore = (pingScore + packetScore) / 2
	}

	// Age score calculation - older nodes are more trusted
	ageInDays := time.Since(n.JoinTimestamp).Hours() / 24
	ageScore := math.Min(1, ageInDays/30) // Max trust after 30 days

	// Update reputation components
	n.Reputation.Behavioral = behaviorScore
	n.Reputation.Network = networkScore
	n.Reputation.Age = ageScore

	// Calculate overall score with weighted components
	weights := map[string]float64{
		"behavioral": 0.4,
		"network":    0.2,
		"consensus":  0.2,
		"age":        0.1,
		"community":  0.1,
	}

	n.Reputation.Overall = (weights["behavioral"]*n.Reputation.Behavioral +
		weights["network"]*n.Reputation.Network +
		weights["consensus"]*n.Reputation.Consensus +
		weights["age"]*n.Reputation.Age +
		weights["community"]*n.Reputation.CommunityTrust)

	n.Reputation.LastCalculated = time.Now()
}

// GenerateHMACToken creates an HMAC token for authentication
func (n *Node) GenerateHMACToken(secret []byte) (string, error) {
	if len(secret) < 16 {
		return "", errors.New("secret is too short for secure HMAC")
	}

	mac := hmac.New(sha256.New, secret)

	// Use a combination of node ID and timestamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	message := n.ID + ":" + timestamp

	mac.Write([]byte(message))
	tokenBytes := mac.Sum(nil)

	// Encode as Base64 for transmission
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	// Add timestamp to token for verification
	return token + ":" + timestamp, nil
}

// VerifyHMACToken verifies an HMAC token for authentication
func VerifyHMACToken(nodeID string, token string, secret []byte) bool {
	parts := []string{}
	for _, part := range parts {
		if part == "" {
			return false
		}
	}

	// Extract token and timestamp
	tokenParts := []string{}
	hmacToken := ""
	timestamp := ""

	if len(tokenParts) != 2 {
		return false
	}

	hmacToken = tokenParts[0]
	timestamp = tokenParts[1]

	// Verify timestamp is recent (within 5 minutes)
	ts, err := timeAsInt64(timestamp)
	if err != nil {
		return false
	}

	now := time.Now().Unix()
	if now-ts > 300 { // 5 minutes
		return false
	}

	// Regenerate HMAC for comparison
	mac := hmac.New(sha256.New, secret)
	message := nodeID + ":" + timestamp
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	expectedToken := base64.StdEncoding.EncodeToString(expectedMAC)

	// Compare tokens
	return subtle.ConstantTimeCompare([]byte(hmacToken), []byte(expectedToken)) == 1
}

// Helper function to convert timestamp string to int64
func timeAsInt64(timestamp string) (int64, error) {
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, err
	}
	return ts.Unix(), nil
}

// ContinuousAuthenticationCheck performs continuous authentication checks
func (n *Node) ContinuousAuthenticationCheck() bool {
	n.ContinuousAuthLock.Lock()
	defer n.ContinuousAuthLock.Unlock()

	// If last activity was too long ago, require re-authentication
	if time.Since(n.LastActiveTime) > 10*time.Minute {
		n.AuthStateValid = false
		return false
	}

	// If reputation has fallen below threshold, require re-authentication
	if n.Reputation.Overall < 0.3 {
		n.AuthStateValid = false
		return false
	}

	return n.AuthStateValid
}

// SignMessage signs a message using the node's private key
func (n *Node) SignMessage(message []byte) ([]byte, error) {
	if n.PrivateKey == nil {
		return nil, errors.New("node doesn't have a private key")
	}

	hash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, n.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Combine r and s into a single signature
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	signature := make([]byte, 64)

	// Ensure proper padding
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):], sBytes)

	return signature, nil
}

// VerifyNodeSignature verifies a signature from another node
func VerifyNodeSignature(pubKey *ecdsa.PublicKey, message, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}

	var r, s big.Int
	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:])

	hash := sha256.Sum256(message)
	return ecdsa.Verify(pubKey, hash[:], &r, &s)
}

// NodeHeartbeat updates the node's active status and performs routine checks
func (n *Node) NodeHeartbeat() {
	n.LastActiveTime = time.Now()

	// Perform continuous authentication check
	n.ContinuousAuthenticationCheck()

	// Every 10th heartbeat, recalculate reputation
	if n.ConnectionMetrics.ConnectionCount%10 == 0 {
		n.RecalculateReputation()
	}

	n.ConnectionMetrics.ConnectionCount++
}

// VerifyAuthenticationState checks if a node's authentication is currently valid
func (n *Node) VerifyAuthenticationState() bool {
	return n.ContinuousAuthenticationCheck()
}

// ResetAuthenticationState forces a node to re-authenticate
func (n *Node) ResetAuthenticationState() {
	n.AuthStateValid = false
	// Generate a new challenge for re-authentication
	challenge, _ := n.GenerateAuthChallenge()
	n.ChallengeNonce = challenge
}
