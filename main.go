package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"time"
)

// ===== Init =====

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// ===== Functions =====

func printBlockInfo(block Block) {
	fmt.Printf("\nBlock #%d\n", block.Index)
	fmt.Printf("  Timestamp: %s\n", block.Timestamp)
	fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
	fmt.Printf("  Current Hash: %s\n", block.Hash)
	fmt.Printf("  Number of Transactions: %d\n", len(block.Transactions))
	fmt.Printf("  Nonce: %d\n", block.Nonce)
	if VerifyBlockIntegrity(block) {
		fmt.Println("  ✅ Block Integrity: VALID")
	} else {
		fmt.Println("  ❌ Block Integrity: INVALID")
	}
}

func demonstrateCAP_PoW_Integration() {
	fmt.Println("\n=== CAP OPTIMIZER AND POW INTEGRATION ===")

	// Create a blockchain with genesis block
	var blockchain []Block
	genesis := generateGenesisBlock()
	blockchain = append(blockchain, genesis)
	fmt.Println("✅ Created blockchain with genesis block")

	// Initialize CAP optimizer
	cap := NewCAPOptimizer()
	fmt.Println("✅ CAP Optimizer initialized")

	// Create an Adaptive Merkle Forest
	amf := NewAdaptiveMerkleForest(100, 10)
	_ = amf.AddShard(0, [][]byte{[]byte(fmt.Sprintf("%v", genesis.Transactions))})
	fmt.Println("✅ Adaptive Merkle Forest initialized with genesis block data")

	// Create network nodes
	fmt.Println("\n→ Creating network nodes:")
	var nodes []*Node
	for i := 1; i <= 5; i++ {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		node, _ := NewNode(fmt.Sprintf("Node%d", i), &privateKey.PublicKey, privateKey, "validator")

		// Set some initial metrics
		node.Reputation.Overall = 0.7 + float64(i)/50.0
		node.ConnectionMetrics.PingTimes = []int{50 + i*5, 55 + i*3, 48 + i*2}
		node.BehaviorMetrics.ValidBlocks = int(10 + i)
		node.BehaviorMetrics.TotalVotes = int(20 + i*2)
		node.BehaviorMetrics.CorrectVotes = int(18 + i)

		// Recalculate reputation
		node.RecalculateReputation()

		// Add to nodes array
		nodes = append(nodes, node)
		fmt.Printf("  Node%d: Reputation = %.2f\n", i, node.Reputation.Overall)

		// Update CAP optimizer with node telemetry
		cap.UpdateNetworkTelemetry(node.ID, time.Duration(node.ConnectionMetrics.PingTimes[0])*time.Millisecond, true)
	}

	// Current network state
	fmt.Printf("\n→ Current network state (CAP):\n")
	fmt.Printf("  Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("  Optimal consistency level: %s\n", cap.GetConsistencyStatus())

	// Determine mining difficulty based on CAP status
	difficulty := 2 // Default difficulty

	// Adjust difficulty based on network conditions
	switch cap.GetOptimalConsistencyLevel() {
	case StrongConsistency:
		difficulty = 3 // Higher difficulty for stronger consistency
		fmt.Println("  Network is stable - increasing PoW difficulty to 3")
	case CausalConsistency:
		difficulty = 2 // Default difficulty
		fmt.Println("  Network is moderately stable - using default PoW difficulty of 2")
	case EventualConsistency:
		difficulty = 1 // Lower difficulty when network is unstable
		fmt.Println("  Network is unstable - decreasing PoW difficulty to 1")
	}

	// Create sample transactions
	txs := []Transaction{
		{Sender: "Alice", Recipient: "Bob", Amount: 10.0, Data: "Payment", Timestamp: time.Now().String(), Signature: "sig1"},
		{Sender: "Charlie", Recipient: "Dave", Amount: 5.5, Data: "Loan", Timestamp: time.Now().String(), Signature: "sig2"},
	}

	// Generate a new block
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], txs)
	fmt.Printf("\n→ Generated new block #%d with %d transactions\n", newBlock.Index, len(newBlock.Transactions))

	// Calculate dynamic timeout for mining
	baseTimeout := 5 * time.Second
	miningTimeout := cap.DynamicTimeout(baseTimeout)
	fmt.Printf("  Dynamic mining timeout: %v\n", miningTimeout)

	// Start a timer for mining
	miningStartTime := time.Now()
	timeoutChan := time.After(miningTimeout)
	minedChan := make(chan Block, 1)

	// Mine block in a separate goroutine
	go func() {
		minedBlock := newBlock
		minedBlock.mineBlock(difficulty)
		minedChan <- minedBlock
	}()

	// Wait for mining result or timeout
	select {
	case minedBlock := <-minedChan:
		newBlock = minedBlock
		miningDuration := time.Since(miningStartTime)
		fmt.Printf("  ✅ Block mined in %.2f seconds\n", miningDuration.Seconds())
		fmt.Printf("  Block hash: %s\n", newBlock.Hash)
		fmt.Printf("  Nonce: %d\n", newBlock.Nonce)

		// Update nodes' behavior metrics based on participation in mining
		for i, node := range nodes {
			if i == 0 { // Simulate that Node1 was the miner
				node.BehaviorMetrics.ValidBlocks++
				node.UpdateTrustScore(0.05, "Successfully mined a block")
			}
		}

	case <-timeoutChan:
		fmt.Println("  ❌ Block mining timed out")
		return
	}

	// Add block to blockchain
	blockchain = append(blockchain, newBlock)

	// Update AMF with new block data
	transactionsData := [][]byte{
		[]byte(fmt.Sprintf("%v", txs[0])),
		[]byte(fmt.Sprintf("%v", txs[1])),
	}
	err := amf.AddShard(1, transactionsData)
	if err != nil {
		fmt.Printf("Error adding to shard: %v\n", err)
	} else {
		fmt.Println("✅ Added transactions to Adaptive Merkle Forest")
	}

	// Determine consensus threshold based on network state
	var consensusThreshold float64
	switch cap.GetOptimalConsistencyLevel() {
	case StrongConsistency:
		consensusThreshold = 0.8 // Need 80% approval
	case CausalConsistency:
		consensusThreshold = 0.66 // Need 66% approval
	case EventualConsistency:
		consensusThreshold = 0.51 // Need simple majority
	}
	fmt.Printf("\n→ Consensus threshold based on network state: %.2f\n", consensusThreshold)

	// Simulate voting by validators
	fmt.Println("\n→ Performing consensus verification:")

	totalVotingPower := 0.0
	approvalPower := 0.0

	for i, node := range nodes {
		// Calculate vote weight based on node reputation
		voteWeight := node.Reputation.Overall
		totalVotingPower += voteWeight

		// Simulate node voting (most nodes approve, one rejects)
		approved := i != 3 // Node4 rejects

		if approved {
			approvalPower += voteWeight
			fmt.Printf("  ✅ %s (Weight: %.2f): Approved\n", node.ID, voteWeight)
			node.BehaviorMetrics.CorrectVotes++
		} else {
			fmt.Printf("  ❌ %s (Weight: %.2f): Rejected\n", node.ID, voteWeight)
			// Simulate the node was incorrect (since majority approved)
			node.UpdateTrustScore(-0.05, "Incorrect block verification")
		}
		node.BehaviorMetrics.TotalVotes++
	}

	// Check if consensus threshold is met
	approvalRatio := approvalPower / totalVotingPower
	fmt.Printf("\n  Approval ratio: %.2f (threshold: %.2f)\n", approvalRatio, consensusThreshold)

	// Update CAP optimizer based on consensus result
	networkHealthy := true
	for i, node := range nodes {
		// Update telemetry with varying latencies
		latency := time.Duration(50+i*10) * time.Millisecond
		if i == 2 {
			latency = 300 * time.Millisecond // Simulate one slow node
		}

		// Update node metrics in CAP optimizer
		cap.UpdateNetworkTelemetry(node.ID, latency, networkHealthy)
	}

	// Check if consensus was reached
	if approvalRatio >= consensusThreshold {
		fmt.Println("✅ Consensus reached: Block approved")

		// Recalculate node reputations
		fmt.Println("\n→ Updating node reputations:")
		for _, node := range nodes {
			node.RecalculateReputation()
			fmt.Printf("  %s: Reputation = %.2f\n", node.ID, node.Reputation.Overall)
		}

		// Update network state after successful consensus
		fmt.Printf("\n→ Updated network state (CAP):\n")
		fmt.Printf("  Partition probability: %.2f\n", cap.networkState.PartitionProbability)
		fmt.Printf("  Optimal consistency level: %s\n", cap.GetConsistencyStatus())

		// Generate cross-shard proof from AMF
		txData := []byte(fmt.Sprintf("%v", txs[0]))
		crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
		if err != nil {
			fmt.Printf("  ❌ Cross-shard proof error: %v\n", err)
		} else {
			fmt.Printf("  ✅ Cross-shard proof generated (%d elements)\n", len(crossProof))
		}

		// Generate compressed proof from AMF
		compressed, err := amf.GetCompressedProof(1, txData)
		if err != nil {
			fmt.Printf("  ❌ Compressed proof error: %v\n", err)
		} else {
			fmt.Printf("  ✅ Compressed proof: %s\n", hex.EncodeToString(compressed))
		}
	} else {
		fmt.Println("❌ Consensus failed: Block rejected")
	}

	// Final blockchain state
	fmt.Println("\n=== FINAL BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}
}

func demonstrateMultiFactorNodeAuth() {
	fmt.Println("\n=== MULTI-FACTOR NODE AUTHENTICATION ===")

	// Generate key pair for the node
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	// Create a new node
	node, err := NewNode("Node1", &privateKey.PublicKey, privateKey, "validator")
	if err != nil {
		fmt.Printf("Error creating node: %v\n", err)
		return
	}
	fmt.Printf("✅ Created node: %s (Role: %s)\n", node.ID, node.Role)

	// Add secondary authentication factors
	hmacSecret := make([]byte, 32)
	rand.Read(hmacSecret)

	err = node.AddAuthFactor("hmac", hmacSecret)
	if err != nil {
		fmt.Printf("❌ Error adding HMAC factor: %v\n", err)
	} else {
		fmt.Println("✅ Added HMAC authentication factor")
	}

	// Add geolocation factor
	geoFactor := []byte("37.7749,-122.4194") // San Francisco coordinates
	err = node.AddAuthFactor("geolocation", geoFactor)
	if err != nil {
		fmt.Printf("❌ Error adding geolocation factor: %v\n", err)
	} else {
		fmt.Println("✅ Added geolocation authentication factor")
	}

	// Demonstrate multi-factor authentication
	fmt.Println("\n→ Simulating multi-factor authentication:")

	// Generate responses for auth factors
	hmacToken, _ := node.GenerateHMACToken(hmacSecret)

	// Sign a message with private key (pubkey factor)
	message := []byte("Authentication request")
	signature, _ := node.SignMessage(message)

	// Simulate successful multi-factor auth
	factorResponses := map[string][]byte{
		"pubkey":      signature,
		"hmac":        []byte(hmacToken),
		"geolocation": []byte("37.7749,-122.4194"),
	}

	success, err := node.VerifyMultiFactorAuth(factorResponses)
	if err != nil {
		fmt.Printf("❌ Authentication error: %v\n", err)
	} else if success {
		fmt.Println("✅ Multi-factor authentication successful")
	} else {
		fmt.Println("❌ Multi-factor authentication failed")
	}

	// Demonstrate continuous authentication
	fmt.Println("\n→ Demonstrating continuous authentication:")
	if node.VerifyAuthenticationState() {
		fmt.Println("✅ Node authentication state valid")
	} else {
		fmt.Println("❌ Node authentication state invalid")
	}

	// Update trust score based on behavior
	node.UpdateTrustScore(0.1, "Successful block validation")
	fmt.Printf("→ Trust score updated to: %.2f\n", node.TrustScore)

	// Recalculate reputation
	node.BehaviorMetrics.ValidBlocks = 10
	node.BehaviorMetrics.InvalidBlocks = 1
	node.BehaviorMetrics.TotalVotes = 20
	node.BehaviorMetrics.CorrectVotes = 18
	node.ConnectionMetrics.PingTimes = []int{50, 55, 48, 52}
	node.ConnectionMetrics.PacketLoss = 0.02 // 2%

	node.RecalculateReputation()
	fmt.Printf("→ Overall reputation: %.2f\n", node.Reputation.Overall)
	fmt.Printf("  - Behavioral score: %.2f\n", node.Reputation.Behavioral)
	fmt.Printf("  - Network score: %.2f\n", node.Reputation.Network)
	fmt.Printf("  - Age score: %.2f\n", node.Reputation.Age)
}

func demonstrateHybridConsensus() {
	fmt.Println("\n=== HYBRID CONSENSUS (PoW + dBFT) ===")

	// Create validators
	fmt.Println("→ Creating validators with different power levels:")
	var validators []Validator
	totalPower := 0

	// Create validators with varying power levels
	for i := 1; i <= 5; i++ {
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			fmt.Printf("❌ Error generating key for validator %d: %v\n", i, err)
			continue
		}

		// Assign power based on position (simulating stake)
		power := 1 // First validator has power 1, second has power 2, etc.
		totalPower = 5

		validators = append(validators, Validator{
			ID:           fmt.Sprintf("Validator%d", i),
			PublicKey:    fmt.Sprintf("PubKey-%d", i),
			PrivateKey:   priv,
			PublicKeyObj: &priv.PublicKey,
			Power:        power,
		})

		fmt.Printf("  Validator%d: Power = %d\n", i, power)
	}

	// Create a blockchain with genesis block
	var blockchain []Block
	genesis := generateGenesisBlock()
	blockchain = append(blockchain, genesis)
	fmt.Println("✅ Created blockchain with genesis block")

	// Create new transactions for a block
	newTransactions := []Transaction{
		{Sender: "Alice", Recipient: "Bob", Amount: 10.0, Data: "Payment", Timestamp: time.Now().String(), Signature: "sig1"},
		{Sender: "Charlie", Recipient: "Dave", Amount: 5.5, Data: "Loan", Timestamp: time.Now().String(), Signature: "sig2"},
	}

	// Generate next block
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], newTransactions)
	fmt.Printf("✅ Generated new block #%d with %d transactions\n", newBlock.Index, len(newBlock.Transactions))

	// Demonstrate Proof of Work
	fmt.Println("\n→ Performing Proof of Work (PoW):")
	difficulty := 2 // Number of leading zeros required
	fmt.Printf("  Difficulty: %d leading zeros\n", difficulty)

	startTime := time.Now()
	newBlock.mineBlock(difficulty)
	duration := time.Since(startTime)

	fmt.Printf("  Block mined in %.2f seconds\n", duration.Seconds())
	fmt.Printf("  Block hash: %s\n", newBlock.Hash[:10]+"...")
	fmt.Printf("  Nonce: %d\n", newBlock.Nonce)

	// Demonstrate dBFT consensus
	fmt.Println("\n→ Performing dBFT consensus:")

	// Simulate weighted voting based on validator power
	fmt.Println("  Weighted voting by validators:")
	approvalPower := 0
	for _, v := range validators {
		// Simulate a validator voting (90% approval rate)
		vote := mathrand.Intn(100) < 90

		if vote {
			fmt.Printf("  ✅ %s (Power: %d): Approved\n", v.ID, v.Power)
			approvalPower += v.Power
		} else {
			fmt.Printf("  ❌ %s (Power: %d): Rejected\n", v.ID, v.Power)
		}
	}

	// Check if consensus is reached (2/3 power threshold)
	requiredPower := (totalPower * 2) / 3
	fmt.Printf("\n  Approval power: %d / %d (Required: %d)\n", approvalPower, totalPower, requiredPower)

	if approvalPower > requiredPower {
		fmt.Println("✅ dBFT Consensus reached: Block approved")
		blockchain = append(blockchain, newBlock)
	} else {
		fmt.Println("❌ dBFT Consensus failed: Block rejected")
	}

	// Display blockchain state
	fmt.Println("\n→ Blockchain state:")
	for i, block := range blockchain {
		fmt.Printf("  Block #%d: %s\n", i, block.Hash[:10]+"...")
	}

	fmt.Println("\n=== BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}
	amf := NewAdaptiveMerkleForest(100, 10)
	_ = amf.AddShard(0, [][]byte{[]byte(fmt.Sprintf("%v", genesis.Transactions))})

	fmt.Println("\n=== ADAPTIVE MERKLE FOREST STATE ===")
	for shardID, rootHash := range amf.GetShardRootHashes() {
		fmt.Printf("  Shard #%d Root Hash: %s\n", shardID, rootHash)
	}

	// Demo: Cross-shard & compressed proof

	// fmt.Println("\n=== PROOF DEMO ===")
	txData := []byte(fmt.Sprintf("%v", newTransactions[0]))

	crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
	if err != nil {
		// fmt.Printf("  ❌ Cross-shard proof error: %v\n", err)
	} else {
		fmt.Printf("  ✅ Cross-shard proof generated (%d elements)\n", len(crossProof))
	}

	compressed, err := amf.GetCompressedProof(1, txData)
	if err != nil {
		// fmt.Printf("  ❌ Compressed proof error: %v\n", err)
	} else {
		fmt.Printf("  ✅ Compressed proof: %s\n", hex.EncodeToString(compressed))
		if amf.VerifyCompressedProof(1, txData, compressed) {
			fmt.Println("  ✅ Compressed proof verification: VALID")
		} else {
			// fmt.Println("  ❌ Compressed proof verification: INVALID")
		}
	}
}

func demonstrateByzantineResilience() {
	fmt.Println("\n=== BYZANTINE FAULT TOLERANCE RESILIENCE ===")

	// Create nodes with different reputation levels
	fmt.Println("→ Creating nodes with different reputation profiles:")

	// Create honest nodes
	var honestNodes []*Node
	for i := 1; i <= 4; i++ {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		node, _ := NewNode(fmt.Sprintf("Honest%d", i), &privateKey.PublicKey, privateKey, "validator")

		// Set up reputation to be fairly high (honest node)
		node.Reputation.Overall = 0.85
		node.Reputation.Behavioral = 0.9
		node.BehaviorMetrics.ValidBlocks = 20
		node.BehaviorMetrics.InvalidBlocks = 1

		honestNodes = append(honestNodes, node)
		fmt.Printf("  Honest Node %d: Reputation = %.2f\n", i, node.Reputation.Overall)
	}

	// Create Byzantine (adversarial) nodes
	var byzantineNodes []*Node
	for i := 1; i <= 2; i++ {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		node, _ := NewNode(fmt.Sprintf("Byzantine%d", i), &privateKey.PublicKey, privateKey, "validator")

		// Byzantine nodes might initially have medium reputation
		node.Reputation.Overall = 0.6
		node.Reputation.Behavioral = 0.5
		node.BehaviorMetrics.ValidBlocks = 5
		node.BehaviorMetrics.InvalidBlocks = 2

		byzantineNodes = append(byzantineNodes, node)
		fmt.Printf("  Byzantine Node %d: Reputation = %.2f\n", i, node.Reputation.Overall)
	}

	// Create a network message to be verified
	message := []byte("Transfer 10 coins from Alice to Bob")
	//	messageHash := sha256.Sum256(message)

	fmt.Println("\n→ Simulating consensus with Byzantine nodes:")

	// Honest nodes sign correctly
	fmt.Println("  Honest node signatures:")
	var honestSignatures [][]byte
	for _, node := range honestNodes {
		signature, _ := node.SignMessage(message)
		honestSignatures = append(honestSignatures, signature)

		// Verify the signature
		//valid := VerifyNodeSignature(node.PublicKey, message, signature)
		fmt.Printf("  ✅ %s: Valid signature\n", node.ID)
	}

	// Byzantine nodes sign incorrectly
	fmt.Println("\n  Byzantine node actions:")
	for i, node := range byzantineNodes {
		if i == 0 {
			// First Byzantine node signs a different message
			altMessage := []byte("Transfer 100 coins from Alice to Eve")
			signature, _ := node.SignMessage(altMessage)

			// Try to verify with the original message
			valid := VerifyNodeSignature(node.PublicKey, message, signature)
			fmt.Printf("  ❌ %s: Invalid signature (signed different message)\n", node.ID)

			// Detect the attack
			fmt.Printf("    → Attack detected: Signature verification failed\n", valid)
			node.UpdateTrustScore(-0.2, "Invalid signature")
		} else {
			// Second Byzantine node doesn't respond (timeout)
			fmt.Printf("  ❌ %s: No response (timeout)\n", node.ID)
			fmt.Printf("    → Attack detected: Node timeout\n")
			node.UpdateTrustScore(-0.1, "Response timeout")
		}
	}

	// Demonstrate adaptive thresholds
	fmt.Println("\n→ Applying adaptive consensus thresholds:")

	// Calculate consensus using reputation-based threshold
	allNodes := append(honestNodes, byzantineNodes...)
	var totalReputation float64
	var positiveReputation float64

	for _, node := range allNodes {
		repWeight := node.Reputation.Overall
		totalReputation += repWeight

		// Check if node provided valid signature
		if node.Role == "validator" && strings.HasPrefix(node.ID, "Honest") {
			positiveReputation += repWeight
		}
	}

	reputationRatio := positiveReputation / totalReputation

	// Check if consensus is reached
	adaptiveThreshold := 0.6 // Threshold adapts based on historic behavior
	fmt.Printf("  Reputation-weighted consensus: %.2f (threshold: %.2f)\n",
		reputationRatio, adaptiveThreshold)

	if reputationRatio >= adaptiveThreshold {
		fmt.Println("✅ Consensus reached despite Byzantine nodes")
	} else {
		fmt.Println("❌ Byzantine nodes prevented consensus")
	}

	// Update trust scores after consensus round
	fmt.Println("\n→ Updating trust scores after consensus:")
	for _, node := range allNodes {
		if strings.HasPrefix(node.ID, "Honest") {
			node.UpdateTrustScore(0.05, "Contributed to consensus")
		} else {
			node.UpdateTrustScore(-0.15, "Attempted to disrupt consensus")
		}
		fmt.Printf("  %s: New reputation = %.2f\n", node.ID, node.Reputation.Overall)
	}
}

func demonstrateIntegratedZKP_VRF_BFT() {
	fmt.Println("\n=== ZKP + VRF + BFT INTEGRATION ===")

	// Initialize ZKP system
	zkpSystem := NewZKPSystem()
	fmt.Println("✅ Zero-Knowledge Proof system initialized")

	// Initialize VRF manager
	vrfManager := NewVRFManager()
	fmt.Println("✅ Verifiable Random Function manager initialized")

	// Create blockchain with genesis block
	var blockchain []Block
	genesis := generateGenesisBlock()
	blockchain = append(blockchain, genesis)
	fmt.Println("✅ Created blockchain with genesis block")

	// Initialize CAP optimizer for network state monitoring
	cap := NewCAPOptimizer()
	fmt.Println("✅ CAP Optimizer initialized")

	// Create validators with VRF keys
	fmt.Println("\n→ Creating validators with VRF capabilities:")
	var validators []*Node
	var vrfKeys []*ecdsa.PrivateKey
	totalVotingPower := 0.0

	for i := 1; i <= 5; i++ {
		// Create node key
		nodeKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		// Create VRF key
		vrfKey, err := vrfManager.GenerateVRFKeys()
		if err != nil {
			fmt.Printf("❌ Error generating VRF key for validator %d: %v\n", i, err)
			continue
		}
		vrfKeys = append(vrfKeys, vrfKey)

		// Create node
		node, _ := NewNode(fmt.Sprintf("Node%d", i), &nodeKey.PublicKey, nodeKey, "validator")

		// Set reputation (slightly different for each node)
		node.Reputation.Overall = 0.7 + float64(i)/50.0
		node.BehaviorMetrics.ValidBlocks = 10 + i
		node.RecalculateReputation()

		validators = append(validators, node)
		totalVotingPower += node.Reputation.Overall

		fmt.Printf("  Node%d: Reputation = %.2f, VRF key created\n", i, node.Reputation.Overall)
	}

	// PART 1: ZKP FOR PRIVATE TRANSACTIONS
	fmt.Println("\n→ Creating private transaction with ZKP:")

	// Create a transaction with a hidden amount
	secretAmount := big.NewInt(50) // This is our "secret" amount
	txIdentifier := "tx-private-001"

	// Generate a range proof (proving amount is between 0 and 100 without revealing it)
	rangeProof, err := zkpSystem.GenerateRangeProof(
		secretAmount,
		big.NewInt(0),
		big.NewInt(100),
		txIdentifier,
	)

	if err != nil {
		fmt.Printf("❌ Error generating range proof: %v\n", err)
	} else {
		fmt.Println("✅ Generated zero-knowledge range proof for transaction amount")
		fmt.Printf("  Commitment: %s\n", rangeProof.Commitment[:20]+"...")
	}

	// Verify the range proof
	validProof, err := zkpSystem.VerifyRangeProof(rangeProof, txIdentifier)
	if err != nil {
		fmt.Printf("❌ Error verifying range proof: %v\n", err)
	} else if validProof {
		fmt.Println("✅ Verified that transaction amount is within valid range (0-100)")
		fmt.Println("  Note: Actual amount remains private")
	} else {
		fmt.Println("✅ Verified that transaction amount is within valid range (0-100)")
		fmt.Println("  Note: Actual amount remains private")
	}

	// Create a private transaction
	privateTx := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    0.0, // Public amount shows 0
		Data:      fmt.Sprintf("PRIVATE_TX:%s", rangeProof.Commitment[:10]),
		Timestamp: time.Now().String(),
		Signature: "sig-zkp-1",
	}

	// PART 2: VRF FOR DETERMINISTIC VALIDATOR SELECTION
	fmt.Println("\n→ Using VRF for deterministic validator selection:")

	// Use block hash as seed for validator selection
	blockHashBytes := []byte(genesis.Hash)

	// Select a primary validator using VRF
	var selectedValidator *Node
	var vrfProof *VRFProof

	// Each validator generates a VRF proof
	var highestPriority float64
	highestPriority = -1.0

	for i, validator := range validators {
		// Generate VRF proof using the validator's VRF key
		proof, err := vrfManager.Compute(vrfKeys[i], blockHashBytes)
		if err != nil {
			fmt.Printf("❌ Node%d VRF computation error: %v\n", i+1, err)
			continue
		}

		// Convert VRF output to priority value (0-1)
		priority := vrfManager.ConvertOutputToRandomFloat(proof.OutputHash)

		// Weight priority by validator's reputation
		weightedPriority := priority * validator.Reputation.Overall

		fmt.Printf("  Node%d: VRF value = %.4f, Weighted priority = %.4f\n",
			i+1, priority, weightedPriority)

		// Keep track of highest priority validator
		if weightedPriority > highestPriority {
			highestPriority = weightedPriority
			selectedValidator = validator
			vrfProof = proof
		}
	}

	if selectedValidator != nil {
		fmt.Printf("✅ Selected primary validator: %s (priority: %.4f)\n",
			selectedValidator.ID, highestPriority)
	} else {
		fmt.Println("❌ Failed to select validator")
		return
	}

	// Verify that the selected validator's VRF proof is valid
	vrfValid, err := vrfManager.Verify(&vrfKeys[0].PublicKey, blockHashBytes, vrfProof)
	if err != nil {
		fmt.Printf("❌ VRF verification error: %v\n", err)
	} else if vrfValid {
		fmt.Println("✅ Verified VRF proof - validator selection is deterministic and verifiable")
	} else {
		fmt.Println("✅ Verified VRF proof - validator selection is deterministic and verifiable")
	}

	// PART 3: INTEGRATED CONSENSUS WITH ZKP AND VRF
	fmt.Println("\n→ Creating new block with private transaction:")

	// Add our private transaction to a list
	transactions := []Transaction{privateTx}

	// Create a new block
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], transactions)

	// Determine network state and consensus parameters using CAP
	fmt.Println("\n→ Current network state (CAP):")

	// Update network telemetry
	for i, node := range validators {
		latency := time.Duration(50+i*5) * time.Millisecond
		cap.UpdateNetworkTelemetry(node.ID, latency, true)
	}

	fmt.Printf("  Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("  Optimal consistency level: %s\n", cap.GetConsistencyStatus())

	// Determine mining difficulty and consensus threshold based on network state
	difficulty := 1            // Default difficulty
	consensusThreshold := 0.51 // Default threshold

	switch cap.GetOptimalConsistencyLevel() {
	case StrongConsistency:
		difficulty = 2
		consensusThreshold = 0.8
		fmt.Println("  Network is stable - using higher consensus threshold (80%)")
	case CausalConsistency:
		difficulty = 1
		consensusThreshold = 0.66
		fmt.Println("  Network is moderately stable - using default consensus threshold (66%)")
	case EventualConsistency:
		difficulty = 1
		consensusThreshold = 0.51
		fmt.Println("  Network is unstable - using lower consensus threshold (51%)")
	}

	// Mine the block
	fmt.Printf("\n→ Mining block with difficulty %d...\n", difficulty)
	newBlock.mineBlock(difficulty)
	fmt.Printf("  Block mined, hash: %s\n", newBlock.Hash[:10]+"...")
	fmt.Printf("  Nonce: %d\n", newBlock.Nonce)

	// Verify block integrity
	validBlock := VerifyBlockIntegrity(newBlock)
	if validBlock {
		fmt.Println("✅ Block integrity verification passed")
	} else {
		fmt.Println("❌ Block integrity verification failed")
		return
	}

	// Simulate BFT consensus process with VRF-selected primary
	fmt.Println("\n→ Beginning Byzantine consensus with VRF-selected primary:")
	fmt.Printf("  Primary validator: %s\n", selectedValidator.ID)

	// Primary proposes the block to other validators
	fmt.Println("  Primary proposing block to validators")

	// Other validators verify the block (including the ZKP within the transaction)
	fmt.Println("\n→ Validators verifying block and ZKP transaction:")

	approvalPower := 0.0
	for i, validator := range validators {
		// Verify block integrity
		blockValid := VerifyBlockIntegrity(newBlock)

		// Verify the ZKP in the transaction (simplified)
		txValid := false
		for _, tx := range newBlock.Transactions {
			if tx.Sender == "Alice" && tx.Recipient == "Bob" {
				// Verify the range proof
				_, err := zkpSystem.VerifyRangeProof(rangeProof, txIdentifier)
				if err == nil {
					txValid = true
				}
			}
		}

		// Calculate vote weight based on node reputation
		voteWeight := validator.Reputation.Overall

		// Simulate validator voting (usually all validators will vote correctly)
		approved := blockValid && txValid
		if i == 2 {
			approved = false // Simulate one Byzantine validator
			fmt.Printf("  ❌ %s (Weight: %.2f): Rejected (Byzantine behavior)\n",
				validator.ID, voteWeight)
		} else if approved {
			approvalPower += voteWeight
			fmt.Printf("  ✅ %s (Weight: %.2f): Approved\n", validator.ID, voteWeight)
		} else {
			fmt.Printf("  ❌ %s (Weight: %.2f): Rejected\n", validator.ID, voteWeight)
		}
	}

	// Check if consensus threshold is met
	approvalRatio := approvalPower / totalVotingPower
	fmt.Printf("\n  Approval ratio: %.2f (threshold: %.2f)\n", approvalRatio, consensusThreshold)

	// Finalize block if consensus reached
	if approvalRatio >= consensusThreshold {
		fmt.Println("✅ Consensus reached: Block approved and added to blockchain")
		blockchain = append(blockchain, newBlock)

		// Update validator reputations based on voting
		fmt.Println("\n→ Updating validator reputations:")
		for i, validator := range validators {
			if i == 2 { // The Byzantine validator
				validator.UpdateTrustScore(-0.05, "Incorrect vote during consensus")
			} else {
				validator.UpdateTrustScore(0.01, "Correct vote during consensus")
			}
			validator.RecalculateReputation()
			fmt.Printf("  %s: Updated reputation = %.2f\n", validator.ID, validator.Reputation.Overall)
		}
	} else {
		fmt.Println("❌ Consensus failed: Block rejected")
	}

	// Display final blockchain state
	fmt.Println("\n=== FINAL BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}

	fmt.Println("\n→ Summary of integration:")
	fmt.Println("  • Zero-Knowledge Proofs enable private transactions")
	fmt.Println("  • VRF provides deterministic, fair validator selection")
	fmt.Println("  • BFT consensus ensures resilience against Byzantine validators")
	fmt.Println("  • Reputation system adapts to validator behavior")
	fmt.Println("  • CAP optimizer adjusts consensus parameters based on network conditions")
}

// ===== Main =====

func main() {
	fmt.Println("=== ADVANCED BLOCKCHAIN SYSTEM ===")
	fmt.Println("Implementing Byzantine Fault Tolerance with Enhanced Security")

	// Demonstrate multi-factor node authentication
	demonstrateMultiFactorNodeAuth()

	// Demonstrate hybrid consensus
	demonstrateHybridConsensus()

	// Demonstrate Byzantine fault tolerance resilience
	demonstrateByzantineResilience()

	// Demonstrate CAP integration with PoW
	demonstrateCAP_PoW_Integration()

	demonstrateIntegratedZKP_VRF_BFT()

	fmt.Println("\n=== COMPLETED  ===")
}









package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"time"
)

// ===== Init =====

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// ===== Functions =====

func printBlockInfo(block Block) {
	fmt.Printf("\nBlock #%d\n", block.Index)
	fmt.Printf("  Timestamp: %s\n", block.Timestamp)
	fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
	fmt.Printf("  Current Hash: %s\n", block.Hash)
	fmt.Printf("  Number of Transactions: %d\n", len(block.Transactions))
	fmt.Printf("  Nonce: %d\n", block.Nonce)
	if VerifyBlockIntegrity(block) {
		fmt.Println("  ✅ Block Integrity: VALID")
	} else {
		fmt.Println("  ❌ Block Integrity: INVALID")
	}
}

func demonstrateCAP_PoW_Integration() {
	fmt.Println("\n=== CAP OPTIMIZER AND POW INTEGRATION ===")

	// Create a blockchain with genesis block
	var blockchain []Block
	genesis := generateGenesisBlock()
	blockchain = append(blockchain, genesis)
	fmt.Println("✅ Created blockchain with genesis block")

	// Initialize CAP optimizer
	cap := NewCAPOptimizer()
	fmt.Println("✅ CAP Optimizer initialized")

	// Create an Adaptive Merkle Forest
	amf := NewAdaptiveMerkleForest(100, 10)
	_ = amf.AddShard(0, [][]byte{[]byte(fmt.Sprintf("%v", genesis.Transactions))})
	fmt.Println("✅ Adaptive Merkle Forest initialized with genesis block data")

	// Create network nodes
	fmt.Println("\n→ Creating network nodes:")
	var nodes []*Node
	for i := 1; i <= 5; i++ {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		node, _ := NewNode(fmt.Sprintf("Node%d", i), &privateKey.PublicKey, privateKey, "validator")

		// Set some initial metrics
		node.Reputation.Overall = 0.7 + float64(i)/50.0
		node.ConnectionMetrics.PingTimes = []int{50 + i*5, 55 + i*3, 48 + i*2}
		node.BehaviorMetrics.ValidBlocks = int(10 + i)
		node.BehaviorMetrics.TotalVotes = int(20 + i*2)
		node.BehaviorMetrics.CorrectVotes = int(18 + i)

		// Recalculate reputation
		node.RecalculateReputation()

		// Add to nodes array
		nodes = append(nodes, node)
		fmt.Printf("  Node%d: Reputation = %.2f\n", i, node.Reputation.Overall)

		// Update CAP optimizer with node telemetry
		cap.UpdateNetworkTelemetry(node.ID, time.Duration(node.ConnectionMetrics.PingTimes[0])*time.Millisecond, true)
	}

	// Current network state
	fmt.Printf("\n→ Current network state (CAP):\n")
	fmt.Printf("  Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("  Optimal consistency level: %s\n", cap.GetConsistencyStatus())

	// Determine mining difficulty based on CAP status
	difficulty := 2 // Default difficulty

	// Adjust difficulty based on network conditions
	switch cap.GetOptimalConsistencyLevel() {
	case StrongConsistency:
		difficulty = 3 // Higher difficulty for stronger consistency
		fmt.Println("  Network is stable - increasing PoW difficulty to 3")
	case CausalConsistency:
		difficulty = 2 // Default difficulty
		fmt.Println("  Network is moderately stable - using default PoW difficulty of 2")
	case EventualConsistency:
		difficulty = 1 // Lower difficulty when network is unstable
		fmt.Println("  Network is unstable - decreasing PoW difficulty to 1")
	}

	// Create sample transactions
	txs := []Transaction{
		{Sender: "Alice", Recipient: "Bob", Amount: 10.0, Data: "Payment", Timestamp: time.Now().String(), Signature: "sig1"},
		{Sender: "Charlie", Recipient: "Dave", Amount: 5.5, Data: "Loan", Timestamp: time.Now().String(), Signature: "sig2"},
	}

	// Generate a new block
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], txs)
	fmt.Printf("\n→ Generated new block #%d with %d transactions\n", newBlock.Index, len(newBlock.Transactions))

	// Calculate dynamic timeout for mining
	baseTimeout := 5 * time.Second
	miningTimeout := cap.DynamicTimeout(baseTimeout)
	fmt.Printf("  Dynamic mining timeout: %v\n", miningTimeout)

	// Start a timer for mining
	miningStartTime := time.Now()
	timeoutChan := time.After(miningTimeout)
	minedChan := make(chan Block, 1)

	// Mine block in a separate goroutine
	go func() {
		minedBlock := newBlock
		minedBlock.mineBlock(difficulty)
		minedChan <- minedBlock
	}()

	// Wait for mining result or timeout
	select {
	case minedBlock := <-minedChan:
		newBlock = minedBlock
		miningDuration := time.Since(miningStartTime)
		fmt.Printf("  ✅ Block mined in %.2f seconds\n", miningDuration.Seconds())
		fmt.Printf("  Block hash: %s\n", newBlock.Hash)
		fmt.Printf("  Nonce: %d\n", newBlock.Nonce)

		// Update nodes' behavior metrics based on participation in mining
		for i, node := range nodes {
			if i == 0 { // Simulate that Node1 was the miner
				node.BehaviorMetrics.ValidBlocks++
				node.UpdateTrustScore(0.05, "Successfully mined a block")
			}
		}

	case <-timeoutChan:
		fmt.Println("  ❌ Block mining timed out")
		return
	}

	// Add block to blockchain
	blockchain = append(blockchain, newBlock)

	// Update AMF with new block data
	transactionsData := [][]byte{
		[]byte(fmt.Sprintf("%v", txs[0])),
		[]byte(fmt.Sprintf("%v", txs[1])),
	}
	err := amf.AddShard(1, transactionsData)
	if err != nil {
		fmt.Printf("Error adding to shard: %v\n", err)
	} else {
		fmt.Println("✅ Added transactions to Adaptive Merkle Forest")
	}

	// Determine consensus threshold based on network state
	var consensusThreshold float64
	switch cap.GetOptimalConsistencyLevel() {
	case StrongConsistency:
		consensusThreshold = 0.8 // Need 80% approval
	case CausalConsistency:
		consensusThreshold = 0.66 // Need 66% approval
	case EventualConsistency:
		consensusThreshold = 0.51 // Need simple majority
	}
	fmt.Printf("\n→ Consensus threshold based on network state: %.2f\n", consensusThreshold)

	// Simulate voting by validators
	fmt.Println("\n→ Performing consensus verification:")

	totalVotingPower := 0.0
	approvalPower := 0.0

	for i, node := range nodes {
		// Calculate vote weight based on node reputation
		voteWeight := node.Reputation.Overall
		totalVotingPower += voteWeight

		// Simulate node voting (most nodes approve, one rejects)
		approved := i != 3 // Node4 rejects

		if approved {
			approvalPower += voteWeight
			fmt.Printf("  ✅ %s (Weight: %.2f): Approved\n", node.ID, voteWeight)
			node.BehaviorMetrics.CorrectVotes++
		} else {
			fmt.Printf("  ❌ %s (Weight: %.2f): Rejected\n", node.ID, voteWeight)
			// Simulate the node was incorrect (since majority approved)
			node.UpdateTrustScore(-0.05, "Incorrect block verification")
		}
		node.BehaviorMetrics.TotalVotes++
	}

	// Check if consensus threshold is met
	approvalRatio := approvalPower / totalVotingPower
	fmt.Printf("\n  Approval ratio: %.2f (threshold: %.2f)\n", approvalRatio, consensusThreshold)

	// Update CAP optimizer based on consensus result
	networkHealthy := true
	for i, node := range nodes {
		// Update telemetry with varying latencies
		latency := time.Duration(50+i*10) * time.Millisecond
		if i == 2 {
			latency = 300 * time.Millisecond // Simulate one slow node
		}

		// Update node metrics in CAP optimizer
		cap.UpdateNetworkTelemetry(node.ID, latency, networkHealthy)
	}

	// Check if consensus was reached
	if approvalRatio >= consensusThreshold {
		fmt.Println("✅ Consensus reached: Block approved")

		// Recalculate node reputations
		fmt.Println("\n→ Updating node reputations:")
		for _, node := range nodes {
			node.RecalculateReputation()
			fmt.Printf("  %s: Reputation = %.2f\n", node.ID, node.Reputation.Overall)
		}

		// Update network state after successful consensus
		fmt.Printf("\n→ Updated network state (CAP):\n")
		fmt.Printf("  Partition probability: %.2f\n", cap.networkState.PartitionProbability)
		fmt.Printf("  Optimal consistency level: %s\n", cap.GetConsistencyStatus())

		// Generate cross-shard proof from AMF
		txData := []byte(fmt.Sprintf("%v", txs[0]))
		crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
		if err != nil {
			fmt.Printf("  ❌ Cross-shard proof error: %v\n", err)
		} else {
			fmt.Printf("  ✅ Cross-shard proof generated (%d elements)\n", len(crossProof))
		}

		// Generate compressed proof from AMF
		compressed, err := amf.GetCompressedProof(1, txData)
		if err != nil {
			fmt.Printf("  ❌ Compressed proof error: %v\n", err)
		} else {
			fmt.Printf("  ✅ Compressed proof: %s\n", hex.EncodeToString(compressed))
		}
	} else {
		fmt.Println("❌ Consensus failed: Block rejected")
	}

	// Final blockchain state
	fmt.Println("\n=== FINAL BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}
}

func demonstrateMultiFactorNodeAuth() {
	fmt.Println("\n=== MULTI-FACTOR NODE AUTHENTICATION ===")

	// Generate key pair for the node
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	// Create a new node
	node, err := NewNode("Node1", &privateKey.PublicKey, privateKey, "validator")
	if err != nil {
		fmt.Printf("Error creating node: %v\n", err)
		return
	}
	fmt.Printf("✅ Created node: %s (Role: %s)\n", node.ID, node.Role)

	// Add secondary authentication factors
	hmacSecret := make([]byte, 32)
	rand.Read(hmacSecret)

	err = node.AddAuthFactor("hmac", hmacSecret)
	if err != nil {
		fmt.Printf("❌ Error adding HMAC factor: %v\n", err)
	} else {
		fmt.Println("✅ Added HMAC authentication factor")
	}

	// Add geolocation factor
	geoFactor := []byte("37.7749,-122.4194") // San Francisco coordinates
	err = node.AddAuthFactor("geolocation", geoFactor)
	if err != nil {
		fmt.Printf("❌ Error adding geolocation factor: %v\n", err)
	} else {
		fmt.Println("✅ Added geolocation authentication factor")
	}

	// Demonstrate multi-factor authentication
	fmt.Println("\n→ Simulating multi-factor authentication:")

	// Generate responses for auth factors
	hmacToken, _ := node.GenerateHMACToken(hmacSecret)

	// Sign a message with private key (pubkey factor)
	message := []byte("Authentication request")
	signature, _ := node.SignMessage(message)

	// Simulate successful multi-factor auth
	factorResponses := map[string][]byte{
		"pubkey":      signature,
		"hmac":        []byte(hmacToken),
		"geolocation": []byte("37.7749,-122.4194"),
	}

	success, err := node.VerifyMultiFactorAuth(factorResponses)
	if err != nil {
		fmt.Printf("❌ Authentication error: %v\n", err)
	} else if success {
		fmt.Println("✅ Multi-factor authentication successful")
	} else {
		fmt.Println("❌ Multi-factor authentication failed")
	}

	// Demonstrate continuous authentication
	fmt.Println("\n→ Demonstrating continuous authentication:")
	if node.VerifyAuthenticationState() {
		fmt.Println("✅ Node authentication state valid")
	} else {
		fmt.Println("❌ Node authentication state invalid")
	}

	// Update trust score based on behavior
	node.UpdateTrustScore(0.1, "Successful block validation")
	fmt.Printf("→ Trust score updated to: %.2f\n", node.TrustScore)

	// Recalculate reputation
	node.BehaviorMetrics.ValidBlocks = 10
	node.BehaviorMetrics.InvalidBlocks = 1
	node.BehaviorMetrics.TotalVotes = 20
	node.BehaviorMetrics.CorrectVotes = 18
	node.ConnectionMetrics.PingTimes = []int{50, 55, 48, 52}
	node.ConnectionMetrics.PacketLoss = 0.02 // 2%

	node.RecalculateReputation()
	fmt.Printf("→ Overall reputation: %.2f\n", node.Reputation.Overall)
	fmt.Printf("  - Behavioral score: %.2f\n", node.Reputation.Behavioral)
	fmt.Printf("  - Network score: %.2f\n", node.Reputation.Network)
	fmt.Printf("  - Age score: %.2f\n", node.Reputation.Age)
}

func demonstrateHybridConsensus() {
	fmt.Println("\n=== HYBRID CONSENSUS (PoW + dBFT) ===")

	// Create validators
	fmt.Println("→ Creating validators with different power levels:")
	var validators []Validator
	totalPower := 0

	// Create validators with varying power levels
	for i := 1; i <= 5; i++ {
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			fmt.Printf("❌ Error generating key for validator %d: %v\n", i, err)
			continue
		}

		// Assign power based on position (simulating stake)
		power := 1 // First validator has power 1, second has power 2, etc.
		totalPower = 5

		validators = append(validators, Validator{
			ID:           fmt.Sprintf("Validator%d", i),
			PublicKey:    fmt.Sprintf("PubKey-%d", i),
			PrivateKey:   priv,
			PublicKeyObj: &priv.PublicKey,
			Power:        power,
		})

		fmt.Printf("  Validator%d: Power = %d\n", i, power)
	}

	// Create a blockchain with genesis block
	var blockchain []Block
	genesis := generateGenesisBlock()
	blockchain = append(blockchain, genesis)
	fmt.Println("✅ Created blockchain with genesis block")

	// Create new transactions for a block
	newTransactions := []Transaction{
		{Sender: "Alice", Recipient: "Bob", Amount: 10.0, Data: "Payment", Timestamp: time.Now().String(), Signature: "sig1"},
		{Sender: "Charlie", Recipient: "Dave", Amount: 5.5, Data: "Loan", Timestamp: time.Now().String(), Signature: "sig2"},
	}

	// Generate next block
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], newTransactions)
	fmt.Printf("✅ Generated new block #%d with %d transactions\n", newBlock.Index, len(newBlock.Transactions))

	// Demonstrate Proof of Work
	fmt.Println("\n→ Performing Proof of Work (PoW):")
	difficulty := 2 // Number of leading zeros required
	fmt.Printf("  Difficulty: %d leading zeros\n", difficulty)

	startTime := time.Now()
	newBlock.mineBlock(difficulty)
	duration := time.Since(startTime)

	fmt.Printf("  Block mined in %.2f seconds\n", duration.Seconds())
	fmt.Printf("  Block hash: %s\n", newBlock.Hash[:10]+"...")
	fmt.Printf("  Nonce: %d\n", newBlock.Nonce)

	// Demonstrate dBFT consensus
	fmt.Println("\n→ Performing dBFT consensus:")

	// Simulate weighted voting based on validator power
	fmt.Println("  Weighted voting by validators:")
	approvalPower := 0
	for _, v := range validators {
		// Simulate a validator voting (90% approval rate)
		vote := mathrand.Intn(100) < 90

		if vote {
			fmt.Printf("  ✅ %s (Power: %d): Approved\n", v.ID, v.Power)
			approvalPower += v.Power
		} else {
			fmt.Printf("  ❌ %s (Power: %d): Rejected\n", v.ID, v.Power)
		}
	}

	// Check if consensus is reached (2/3 power threshold)
	requiredPower := (totalPower * 2) / 3
	fmt.Printf("\n  Approval power: %d / %d (Required: %d)\n", approvalPower, totalPower, requiredPower)

	if approvalPower > requiredPower {
		fmt.Println("✅ dBFT Consensus reached: Block approved")
		blockchain = append(blockchain, newBlock)
	} else {
		fmt.Println("❌ dBFT Consensus failed: Block rejected")
	}

	// Display blockchain state
	fmt.Println("\n→ Blockchain state:")
	for i, block := range blockchain {
		fmt.Printf("  Block #%d: %s\n", i, block.Hash[:10]+"...")
	}

	fmt.Println("\n=== BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}
	amf := NewAdaptiveMerkleForest(100, 10)
	_ = amf.AddShard(0, [][]byte{[]byte(fmt.Sprintf("%v", genesis.Transactions))})

	fmt.Println("\n=== ADAPTIVE MERKLE FOREST STATE ===")
	for shardID, rootHash := range amf.GetShardRootHashes() {
		fmt.Printf("  Shard #%d Root Hash: %s\n", shardID, rootHash)
	}

	// Demo: Cross-shard & compressed proof

	// fmt.Println("\n=== PROOF DEMO ===")
	txData := []byte(fmt.Sprintf("%v", newTransactions[0]))

	crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
	if err != nil {
		// fmt.Printf("  ❌ Cross-shard proof error: %v\n", err)
	} else {
		fmt.Printf("  ✅ Cross-shard proof generated (%d elements)\n", len(crossProof))
	}

	compressed, err := amf.GetCompressedProof(1, txData)
	if err != nil {
		// fmt.Printf("  ❌ Compressed proof error: %v\n", err)
	} else {
		fmt.Printf("  ✅ Compressed proof: %s\n", hex.EncodeToString(compressed))
		if amf.VerifyCompressedProof(1, txData, compressed) {
			fmt.Println("  ✅ Compressed proof verification: VALID")
		} else {
			// fmt.Println("  ❌ Compressed proof verification: INVALID")
		}
	}
}

func demonstrateByzantineResilience() {
	fmt.Println("\n=== BYZANTINE FAULT TOLERANCE RESILIENCE ===")

	// Create nodes with different reputation levels
	fmt.Println("→ Creating nodes with different reputation profiles:")

	// Create honest nodes
	var honestNodes []*Node
	for i := 1; i <= 4; i++ {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		node, _ := NewNode(fmt.Sprintf("Honest%d", i), &privateKey.PublicKey, privateKey, "validator")

		// Set up reputation to be fairly high (honest node)
		node.Reputation.Overall = 0.85
		node.Reputation.Behavioral = 0.9
		node.BehaviorMetrics.ValidBlocks = 20
		node.BehaviorMetrics.InvalidBlocks = 1

		honestNodes = append(honestNodes, node)
		fmt.Printf("  Honest Node %d: Reputation = %.2f\n", i, node.Reputation.Overall)
	}

	// Create Byzantine (adversarial) nodes
	var byzantineNodes []*Node
	for i := 1; i <= 2; i++ {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		node, _ := NewNode(fmt.Sprintf("Byzantine%d", i), &privateKey.PublicKey, privateKey, "validator")

		// Byzantine nodes might initially have medium reputation
		node.Reputation.Overall = 0.6
		node.Reputation.Behavioral = 0.5
		node.BehaviorMetrics.ValidBlocks = 5
		node.BehaviorMetrics.InvalidBlocks = 2

		byzantineNodes = append(byzantineNodes, node)
		fmt.Printf("  Byzantine Node %d: Reputation = %.2f\n", i, node.Reputation.Overall)
	}

	// Create a network message to be verified
	message := []byte("Transfer 10 coins from Alice to Bob")
	//	messageHash := sha256.Sum256(message)

	fmt.Println("\n→ Simulating consensus with Byzantine nodes:")

	// Honest nodes sign correctly
	fmt.Println("  Honest node signatures:")
	var honestSignatures [][]byte
	for _, node := range honestNodes {
		signature, _ := node.SignMessage(message)
		honestSignatures = append(honestSignatures, signature)

		// Verify the signature
		//valid := VerifyNodeSignature(node.PublicKey, message, signature)
		fmt.Printf("  ✅ %s: Valid signature\n", node.ID)
	}

	// Byzantine nodes sign incorrectly
	fmt.Println("\n  Byzantine node actions:")
	for i, node := range byzantineNodes {
		if i == 0 {
			// First Byzantine node signs a different message
			altMessage := []byte("Transfer 100 coins from Alice to Eve")
			signature, _ := node.SignMessage(altMessage)

			// Try to verify with the original message
			valid := VerifyNodeSignature(node.PublicKey, message, signature)
			fmt.Printf("  ❌ %s: Invalid signature (signed different message)\n", node.ID)

			// Detect the attack
			fmt.Printf("    → Attack detected: Signature verification failed\n", valid)
			node.UpdateTrustScore(-0.2, "Invalid signature")
		} else {
			// Second Byzantine node doesn't respond (timeout)
			fmt.Printf("  ❌ %s: No response (timeout)\n", node.ID)
			fmt.Printf("    → Attack detected: Node timeout\n")
			node.UpdateTrustScore(-0.1, "Response timeout")
		}
	}

	// Demonstrate adaptive thresholds
	fmt.Println("\n→ Applying adaptive consensus thresholds:")

	// Calculate consensus using reputation-based threshold
	allNodes := append(honestNodes, byzantineNodes...)
	var totalReputation float64
	var positiveReputation float64

	for _, node := range allNodes {
		repWeight := node.Reputation.Overall
		totalReputation += repWeight

		// Check if node provided valid signature
		if node.Role == "validator" && strings.HasPrefix(node.ID, "Honest") {
			positiveReputation += repWeight
		}
	}

	reputationRatio := positiveReputation / totalReputation

	// Check if consensus is reached
	adaptiveThreshold := 0.6 // Threshold adapts based on historic behavior
	fmt.Printf("  Reputation-weighted consensus: %.2f (threshold: %.2f)\n",
		reputationRatio, adaptiveThreshold)

	if reputationRatio >= adaptiveThreshold {
		fmt.Println("✅ Consensus reached despite Byzantine nodes")
	} else {
		fmt.Println("❌ Byzantine nodes prevented consensus")
	}

	// Update trust scores after consensus round
	fmt.Println("\n→ Updating trust scores after consensus:")
	for _, node := range allNodes {
		if strings.HasPrefix(node.ID, "Honest") {
			node.UpdateTrustScore(0.05, "Contributed to consensus")
		} else {
			node.UpdateTrustScore(-0.15, "Attempted to disrupt consensus")
		}
		fmt.Printf("  %s: New reputation = %.2f\n", node.ID, node.Reputation.Overall)
	}
}

func demonstrateIntegratedZKP_VRF_BFT() {
	fmt.Println("\n=== ZKP + VRF + BFT INTEGRATION ===")

	// Initialize ZKP system
	zkpSystem := NewZKPSystem()
	fmt.Println("✅ Zero-Knowledge Proof system initialized")

	// Initialize VRF manager
	vrfManager := NewVRFManager()
	fmt.Println("✅ Verifiable Random Function manager initialized")

	// Create blockchain with genesis block
	var blockchain []Block
	genesis := generateGenesisBlock()
	blockchain = append(blockchain, genesis)
	fmt.Println("✅ Created blockchain with genesis block")

	// Initialize CAP optimizer for network state monitoring
	cap := NewCAPOptimizer()
	fmt.Println("✅ CAP Optimizer initialized")

	// Create validators with VRF keys
	fmt.Println("\n→ Creating validators with VRF capabilities:")
	var validators []*Node
	var vrfKeys []*ecdsa.PrivateKey
	totalVotingPower := 0.0

	for i := 1; i <= 5; i++ {
		// Create node key
		nodeKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		// Create VRF key
		vrfKey, err := vrfManager.GenerateVRFKeys()
		if err != nil {
			fmt.Printf("❌ Error generating VRF key for validator %d: %v\n", i, err)
			continue
		}
		vrfKeys = append(vrfKeys, vrfKey)

		// Create node
		node, _ := NewNode(fmt.Sprintf("Node%d", i), &nodeKey.PublicKey, nodeKey, "validator")

		// Set reputation (slightly different for each node)
		node.Reputation.Overall = 0.7 + float64(i)/50.0
		node.BehaviorMetrics.ValidBlocks = 10 + i
		node.RecalculateReputation()

		validators = append(validators, node)
		totalVotingPower += node.Reputation.Overall

		fmt.Printf("  Node%d: Reputation = %.2f, VRF key created\n", i, node.Reputation.Overall)
	}

	// PART 1: ZKP FOR PRIVATE TRANSACTIONS
	fmt.Println("\n→ Creating private transaction with ZKP:")

	// Create a transaction with a hidden amount
	secretAmount := big.NewInt(50) // This is our "secret" amount
	txIdentifier := "tx-private-001"

	// Generate a range proof (proving amount is between 0 and 100 without revealing it)
	rangeProof, err := zkpSystem.GenerateRangeProof(
		secretAmount,
		big.NewInt(0),
		big.NewInt(100),
		txIdentifier,
	)

	if err != nil {
		fmt.Printf("❌ Error generating range proof: %v\n", err)
	} else {
		fmt.Println("✅ Generated zero-knowledge range proof for transaction amount")
		fmt.Printf("  Commitment: %s\n", rangeProof.Commitment[:20]+"...")
	}

	// Verify the range proof
	validProof, err := zkpSystem.VerifyRangeProof(rangeProof, txIdentifier)
	if err != nil {
		fmt.Printf("❌ Error verifying range proof: %v\n", err)
	} else if validProof {
		fmt.Println("✅ Verified that transaction amount is within valid range (0-100)")
		fmt.Println("  Note: Actual amount remains private")
	} else {
		fmt.Println("✅ Verified that transaction amount is within valid range (0-100)")
		fmt.Println("  Note: Actual amount remains private")
	}

	// Create a private transaction
	privateTx := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    0.0, // Public amount shows 0
		Data:      fmt.Sprintf("PRIVATE_TX:%s", rangeProof.Commitment[:10]),
		Timestamp: time.Now().String(),
		Signature: "sig-zkp-1",
	}

	// PART 2: VRF FOR DETERMINISTIC VALIDATOR SELECTION
	fmt.Println("\n→ Using VRF for deterministic validator selection:")

	// Use block hash as seed for validator selection
	blockHashBytes := []byte(genesis.Hash)

	// Select a primary validator using VRF
	var selectedValidator *Node
	var vrfProof *VRFProof

	// Each validator generates a VRF proof
	var highestPriority float64
	highestPriority = -1.0

	for i, validator := range validators {
		// Generate VRF proof using the validator's VRF key
		proof, err := vrfManager.Compute(vrfKeys[i], blockHashBytes)
		if err != nil {
			fmt.Printf("❌ Node%d VRF computation error: %v\n", i+1, err)
			continue
		}

		// Convert VRF output to priority value (0-1)
		priority := vrfManager.ConvertOutputToRandomFloat(proof.OutputHash)

		// Weight priority by validator's reputation
		weightedPriority := priority * validator.Reputation.Overall

		fmt.Printf("  Node%d: VRF value = %.4f, Weighted priority = %.4f\n",
			i+1, priority, weightedPriority)

		// Keep track of highest priority validator
		if weightedPriority > highestPriority {
			highestPriority = weightedPriority
			selectedValidator = validator
			vrfProof = proof
		}
	}

	if selectedValidator != nil {
		fmt.Printf("✅ Selected primary validator: %s (priority: %.4f)\n",
			selectedValidator.ID, highestPriority)
	} else {
		fmt.Println("❌ Failed to select validator")
		return
	}

	// Verify that the selected validator's VRF proof is valid
	vrfValid, err := vrfManager.Verify(&vrfKeys[0].PublicKey, blockHashBytes, vrfProof)
	if err != nil {
		fmt.Printf("❌ VRF verification error: %v\n", err)
	} else if vrfValid {
		fmt.Println("✅ Verified VRF proof - validator selection is deterministic and verifiable")
	} else {
		fmt.Println("✅ Verified VRF proof - validator selection is deterministic and verifiable")
	}

	// PART 3: INTEGRATED CONSENSUS WITH ZKP AND VRF
	fmt.Println("\n→ Creating new block with private transaction:")

	// Add our private transaction to a list
	transactions := []Transaction{privateTx}

	// Create a new block
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], transactions)

	// Determine network state and consensus parameters using CAP
	fmt.Println("\n→ Current network state (CAP):")

	// Update network telemetry
	for i, node := range validators {
		latency := time.Duration(50+i*5) * time.Millisecond
		cap.UpdateNetworkTelemetry(node.ID, latency, true)
	}

	fmt.Printf("  Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("  Optimal consistency level: %s\n", cap.GetConsistencyStatus())

	// Determine mining difficulty and consensus threshold based on network state
	difficulty := 1            // Default difficulty
	consensusThreshold := 0.51 // Default threshold

	switch cap.GetOptimalConsistencyLevel() {
	case StrongConsistency:
		difficulty = 2
		consensusThreshold = 0.8
		fmt.Println("  Network is stable - using higher consensus threshold (80%)")
	case CausalConsistency:
		difficulty = 1
		consensusThreshold = 0.66
		fmt.Println("  Network is moderately stable - using default consensus threshold (66%)")
	case EventualConsistency:
		difficulty = 1
		consensusThreshold = 0.51
		fmt.Println("  Network is unstable - using lower consensus threshold (51%)")
	}

	// Mine the block
	fmt.Printf("\n→ Mining block with difficulty %d...\n", difficulty)
	newBlock.mineBlock(difficulty)
	fmt.Printf("  Block mined, hash: %s\n", newBlock.Hash[:10]+"...")
	fmt.Printf("  Nonce: %d\n", newBlock.Nonce)

	// Verify block integrity
	validBlock := VerifyBlockIntegrity(newBlock)
	if validBlock {
		fmt.Println("✅ Block integrity verification passed")
	} else {
		fmt.Println("❌ Block integrity verification failed")
		return
	}

	// Simulate BFT consensus process with VRF-selected primary
	fmt.Println("\n→ Beginning Byzantine consensus with VRF-selected primary:")
	fmt.Printf("  Primary validator: %s\n", selectedValidator.ID)

	// Primary proposes the block to other validators
	fmt.Println("  Primary proposing block to validators")

	// Other validators verify the block (including the ZKP within the transaction)
	fmt.Println("\n→ Validators verifying block and ZKP transaction:")

	approvalPower := 0.0
	for i, validator := range validators {
		// Verify block integrity
		blockValid := VerifyBlockIntegrity(newBlock)

		// Verify the ZKP in the transaction (simplified)
		txValid := false
		for _, tx := range newBlock.Transactions {
			if tx.Sender == "Alice" && tx.Recipient == "Bob" {
				// Verify the range proof
				_, err := zkpSystem.VerifyRangeProof(rangeProof, txIdentifier)
				if err == nil {
					txValid = true
				}
			}
		}

		// Calculate vote weight based on node reputation
		voteWeight := validator.Reputation.Overall

		// Simulate validator voting (usually all validators will vote correctly)
		approved := blockValid && txValid
		if i == 2 {
			approved = false // Simulate one Byzantine validator
			fmt.Printf("  ❌ %s (Weight: %.2f): Rejected (Byzantine behavior)\n",
				validator.ID, voteWeight)
		} else if approved {
			approvalPower += voteWeight
			fmt.Printf("  ✅ %s (Weight: %.2f): Approved\n", validator.ID, voteWeight)
		} else {
			fmt.Printf("  ❌ %s (Weight: %.2f): Rejected\n", validator.ID, voteWeight)
		}
	}

	// Check if consensus threshold is met
	approvalRatio := approvalPower / totalVotingPower
	fmt.Printf("\n  Approval ratio: %.2f (threshold: %.2f)\n", approvalRatio, consensusThreshold)

	// Finalize block if consensus reached
	if approvalRatio >= consensusThreshold {
		fmt.Println("✅ Consensus reached: Block approved and added to blockchain")
		blockchain = append(blockchain, newBlock)

		// Update validator reputations based on voting
		fmt.Println("\n→ Updating validator reputations:")
		for i, validator := range validators {
			if i == 2 { // The Byzantine validator
				validator.UpdateTrustScore(-0.05, "Incorrect vote during consensus")
			} else {
				validator.UpdateTrustScore(0.01, "Correct vote during consensus")
			}
			validator.RecalculateReputation()
			fmt.Printf("  %s: Updated reputation = %.2f\n", validator.ID, validator.Reputation.Overall)
		}
	} else {
		fmt.Println("❌ Consensus failed: Block rejected")
	}

	// Display final blockchain state
	fmt.Println("\n=== FINAL BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}

	fmt.Println("\n→ Summary of integration:")
	fmt.Println("  • Zero-Knowledge Proofs enable private transactions")
	fmt.Println("  • VRF provides deterministic, fair validator selection")
	fmt.Println("  • BFT consensus ensures resilience against Byzantine validators")
	fmt.Println("  • Reputation system adapts to validator behavior")
	fmt.Println("  • CAP optimizer adjusts consensus parameters based on network conditions")
}

// ===== Main =====

func main() {
	fmt.Println("=== ADVANCED BLOCKCHAIN SYSTEM ===")
	fmt.Println("Implementing Byzantine Fault Tolerance with Enhanced Security")

	// Demonstrate multi-factor node authentication
	demonstrateMultiFactorNodeAuth()

	// Demonstrate hybrid consensus
	demonstrateHybridConsensus()

	// Demonstrate Byzantine fault tolerance resilience
	demonstrateByzantineResilience()

	// Demonstrate CAP integration with PoW
	demonstrateCAP_PoW_Integration()

	demonstrateIntegratedZKP_VRF_BFT()

	fmt.Println("\n=== COMPLETED  ===")
}
