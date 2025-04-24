// package main

// import "fmt"

// func main() {
// 	var Blockchain []Block

// 	// Create genesis block
// 	genesis := generateGenesisBlock()
// 	Blockchain = append(Blockchain, genesis)

// 	// Add a couple more blocks
// 	block1 := generateNextBlock(Blockchain[len(Blockchain)-1], "First real block")
// 	Blockchain = append(Blockchain, block1)

// 	block2 := generateNextBlock(Blockchain[len(Blockchain)-1], "Second real block")
// 	Blockchain = append(Blockchain, block2)

// 	// Print the blockchain
// 	for _, block := range Blockchain {
// 		fmt.Printf("\nIndex: %d\nTimestamp: %s\nData: %s\nPrevHash: %s\nHash: %s\n",
// 			block.Index, block.Timestamp, block.Data, block.PrevHash, block.Hash)
// 	}
// }

// package main

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"time"
// )

// func main() {
// 	var Blockchain []Block

// 	// Create genesis block
// 	genesis := generateGenesisBlock()
// 	Blockchain = append(Blockchain, genesis)

// 	// Create an Adaptive Merkle Forest
// 	amf := NewAdaptiveMerkleForest(100, 10)

// 	// Add first shard using genesis block
// 	transactionsData := [][]byte{[]byte(fmt.Sprintf("%v", Blockchain[0].Transactions[0]))}
// 	err := amf.AddShard(0, transactionsData)
// 	if err != nil {
// 		fmt.Printf("Error creating shard: %v\n", err)
// 	}

// 	// Create some transactions for block 1
// 	tx1 := Transaction{
// 		Sender:    "Alice",
// 		Recipient: "Bob",
// 		Amount:    10.0,
// 		Data:      "Payment for services",
// 		Timestamp: time.Now().String(),
// 		Signature: "signature1",
// 	}

// 	tx2 := Transaction{
// 		Sender:    "Charlie",
// 		Recipient: "Dave",
// 		Amount:    5.5,
// 		Data:      "Loan repayment",
// 		Timestamp: time.Now().String(),
// 		Signature: "signature2",
// 	}

// 	// Create block 1
// 	block1 := generateNextBlock(Blockchain[len(Blockchain)-1], []Transaction{tx1, tx2})
// 	Blockchain = append(Blockchain, block1)

// 	// Add transactions to AMF
// 	transactionsData = [][]byte{
// 		[]byte(fmt.Sprintf("%v", tx1)),
// 		[]byte(fmt.Sprintf("%v", tx2)),
// 	}
// 	err = amf.AddShard(1, transactionsData)
// 	if err != nil {
// 		fmt.Printf("Error adding to shard: %v\n", err)
// 	}

// 	// Create more transactions for block 2
// 	tx3 := Transaction{
// 		Sender:    "Eve",
// 		Recipient: "Frank",
// 		Amount:    25.0,
// 		Data:      "Investment",
// 		Timestamp: time.Now().String(),
// 		Signature: "signature3",
// 	}

// 	tx4 := Transaction{
// 		Sender:    "Grace",
// 		Recipient: "Heidi",
// 		Amount:    12.3,
// 		Data:      "Monthly subscription",
// 		Timestamp: time.Now().String(),
// 		Signature: "signature4",
// 	}

// 	tx5 := Transaction{
// 		Sender:    "Ivan",
// 		Recipient: "Julia",
// 		Amount:    7.8,
// 		Data:      "Birthday gift",
// 		Timestamp: time.Now().String(),
// 		Signature: "signature5",
// 	}

// 	// Create block 2
// 	block2 := generateNextBlock(Blockchain[len(Blockchain)-1], []Transaction{tx3, tx4, tx5})
// 	Blockchain = append(Blockchain, block2)

// 	// Add transactions to AMF
// 	transactionsData = [][]byte{
// 		[]byte(fmt.Sprintf("%v", tx3)),
// 		[]byte(fmt.Sprintf("%v", tx4)),
// 		[]byte(fmt.Sprintf("%v", tx5)),
// 	}
// 	err = amf.AddShard(2, transactionsData)
// 	if err != nil {
// 		fmt.Printf("Error adding to shard: %v\n", err)
// 	}

// 	// Print blockchain with Merkle Tree info
// 	fmt.Println("\n=== BLOCKCHAIN WITH MERKLE TREES ===")
// 	for i, block := range Blockchain {
// 		fmt.Printf("\nBlock #%d:\n", i)
// 		fmt.Printf("  Index: %d\n", block.Index)
// 		fmt.Printf("  Timestamp: %s\n", block.Timestamp)
// 		fmt.Printf("  PrevHash: %s\n", block.PrevHash)
// 		fmt.Printf("  Hash: %s\n", block.Hash)
// 		fmt.Printf("  Merkle Root: %s\n", block.MerkleRoot)
// 		fmt.Printf("  Transactions: %d\n", len(block.Transactions))

// 		// Verify block integrity
// 		if VerifyBlockIntegrity(block) {
// 			fmt.Println("  Block Integrity: VALID")
// 		} else {
// 			fmt.Println("  Block Integrity: INVALID")
// 		}
// 	}

// 	// Demonstrate AMF operations
// 	fmt.Println("\n=== ADAPTIVE MERKLE FOREST OPERATIONS ===")

// 	// Get all shard root hashes
// 	rootHashes := amf.GetShardRootHashes()
// 	fmt.Println("Shard Root Hashes:")
// 	for shardID, rootHash := range rootHashes {
// 		fmt.Printf("  Shard #%d: %s\n", shardID, rootHash)
// 	}

// 	// Demonstrate cross-shard proof
// 	fmt.Println("\nGenerating Cross-Shard Proof:")
// 	txData := []byte(fmt.Sprintf("%v", tx1))
// 	crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
// 	if err != nil {
// 		fmt.Printf("  Error: %v\n", err)
// 	} else {
// 		fmt.Printf("  Cross-shard proof generated: %d elements\n", len(crossProof))
// 	}

// 	// Demonstrate compressed proof
// 	fmt.Println("\nGenerating Compressed Proof:")
// 	compressedProof, err := amf.GetCompressedProof(1, txData)
// 	if err != nil {
// 		fmt.Printf("  Error: %v\n", err)
// 	} else {
// 		fmt.Printf("  Compressed proof: %s\n", hex.EncodeToString(compressedProof))

// 		// Verify compressed proof
// 		if amf.VerifyCompressedProof(1, txData, compressedProof) {
// 			fmt.Println("  Compressed proof verification: VALID")
// 		} else {
// 			fmt.Println("  Compressed proof verification: INVALID")
// 		}
// 	}

// 	// Demonstrate shard split
// 	fmt.Println("\nDemonstrating Shard Split:")
// 	shard1, shard2, err := amf.SplitShard(2)
// 	if err != nil {
// 		fmt.Printf("  Error: %v\n", err)
// 	} else {
// 		fmt.Printf("  Shard #2 split into Shards #%d and #%d\n", shard1, shard2)

// 		// Print updated shard info
// 		fmt.Println("  Updated Shard Info:")
// 		for shardID, info := range amf.ShardInfo {
// 			fmt.Printf("    Shard #%d: Size=%d, Load=%.2f\n",
// 				shardID, info.Size, info.ComputationalLoad)
// 		}
// 	}

// }

package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	var Blockchain []Block

	// Create genesis block
	genesis := generateGenesisBlock()
	Blockchain = append(Blockchain, genesis)

	// Create an Adaptive Merkle Forest
	amf := NewAdaptiveMerkleForest(100, 10)

	// Add first shard using genesis block
	transactionsData := [][]byte{[]byte(fmt.Sprintf("%v", Blockchain[0].Transactions[0]))}
	err := amf.AddShard(0, transactionsData)
	if err != nil {
		fmt.Printf("Error creating shard: %v\n", err)
	}

	// Create some transactions for block 1
	tx1 := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    10.0,
		Data:      "Payment for services",
		Timestamp: time.Now().String(),
		Signature: "signature1",
	}

	tx2 := Transaction{
		Sender:    "Charlie",
		Recipient: "Dave",
		Amount:    5.5,
		Data:      "Loan repayment",
		Timestamp: time.Now().String(),
		Signature: "signature2",
	}

	// Create block 1
	block1 := generateNextBlock(Blockchain[len(Blockchain)-1], []Transaction{tx1, tx2})
	Blockchain = append(Blockchain, block1)

	// Add transactions to AMF
	transactionsData = [][]byte{
		[]byte(fmt.Sprintf("%v", tx1)),
		[]byte(fmt.Sprintf("%v", tx2)),
	}
	err = amf.AddShard(1, transactionsData)
	if err != nil {
		fmt.Printf("Error adding to shard: %v\n", err)
	}

	// Create more transactions for block 2
	tx3 := Transaction{
		Sender:    "Eve",
		Recipient: "Frank",
		Amount:    25.0,
		Data:      "Investment",
		Timestamp: time.Now().String(),
		Signature: "signature3",
	}

	tx4 := Transaction{
		Sender:    "Grace",
		Recipient: "Heidi",
		Amount:    12.3,
		Data:      "Monthly subscription",
		Timestamp: time.Now().String(),
		Signature: "signature4",
	}

	tx5 := Transaction{
		Sender:    "Ivan",
		Recipient: "Julia",
		Amount:    7.8,
		Data:      "Birthday gift",
		Timestamp: time.Now().String(),
		Signature: "signature5",
	}

	// Create block 2
	block2 := generateNextBlock(Blockchain[len(Blockchain)-1], []Transaction{tx3, tx4, tx5})
	Blockchain = append(Blockchain, block2)

	// Add transactions to AMF
	transactionsData = [][]byte{
		[]byte(fmt.Sprintf("%v", tx3)),
		[]byte(fmt.Sprintf("%v", tx4)),
		[]byte(fmt.Sprintf("%v", tx5)),
	}
	err = amf.AddShard(2, transactionsData)
	if err != nil {
		fmt.Printf("Error adding to shard: %v\n", err)
	}

	// Print blockchain with Merkle Tree info
	fmt.Println("\n=== BLOCKCHAIN WITH MERKLE TREES ===")
	for i, block := range Blockchain {
		fmt.Printf("\nBlock #%d:\n", i)
		fmt.Printf("  Index: %d\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp)
		fmt.Printf("  PrevHash: %s\n", block.PrevHash)
		fmt.Printf("  Hash: %s\n", block.Hash)
		fmt.Printf("  Merkle Root: %s\n", block.MerkleRoot)
		fmt.Printf("  Transactions: %d\n", len(block.Transactions))

		// Verify block integrity
		if VerifyBlockIntegrity(block) {
			fmt.Println("  Block Integrity: VALID")
		} else {
			fmt.Println("  Block Integrity: INVALID")
		}
	}

	// Demonstrate AMF operations
	fmt.Println("\n=== ADAPTIVE MERKLE FOREST OPERATIONS ===")

	// Get all shard root hashes
	rootHashes := amf.GetShardRootHashes()
	fmt.Println("Shard Root Hashes:")
	for shardID, rootHash := range rootHashes {
		fmt.Printf("  Shard #%d: %s\n", shardID, rootHash)
	}

	// Demonstrate cross-shard proof
	fmt.Println("\nGenerating Cross-Shard Proof:")
	txData := []byte(fmt.Sprintf("%v", tx1))
	crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Cross-shard proof generated: %d elements\n", len(crossProof))
	}

	// Demonstrate compressed proof
	fmt.Println("\nGenerating Compressed Proof:")
	compressedProof, err := amf.GetCompressedProof(1, txData)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Compressed proof: %s\n", hex.EncodeToString(compressedProof))

		// Verify compressed proof
		if amf.VerifyCompressedProof(1, txData, compressedProof) {
			fmt.Println("  Compressed proof verification: VALID")
		} else {
			fmt.Println("  Compressed proof verification: INVALID")
		}
	}

	// Demonstrate shard split
	fmt.Println("\nDemonstrating Shard Split:")
	shard1, shard2, err := amf.SplitShard(2)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Shard #2 split into Shards #%d and #%d\n", shard1, shard2)

		// Print updated shard info
		fmt.Println("  Updated Shard Info:")
		for shardID, info := range amf.ShardInfo {
			fmt.Printf("    Shard #%d: Size=%d, Load=%.2f\n",
				shardID, info.Size, info.ComputationalLoad)
		}
	}

	// ----- CAP OPTIMIZER TESTING -----
	fmt.Println("\n=== CAP OPTIMIZER TESTING ===")

	// Initialize CAP optimizer
	cap := NewCAPOptimizer()

	// Stage 1: Normal network conditions
	fmt.Println("\n[Stage 1] Normal network conditions:")
	// Simulate good network conditions (low latency, high reliability)
	for i := 1; i <= 5; i++ {
		nodeID := fmt.Sprintf("node%d", i)
		cap.UpdateNetworkTelemetry(nodeID, 50*time.Millisecond, true)
	}

	// Check the consistency level (should be strong in good conditions)
	fmt.Printf("Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("Optimal consistency level: %s\n", cap.GetConsistencyStatus())
	fmt.Printf("Dynamic timeout for 1s base: %v\n", cap.DynamicTimeout(time.Second))

	// Stage 2: Degraded network conditions
	fmt.Println("\n[Stage 2] Degraded network conditions:")
	// Simulate some network issues (higher latency, some failures)
	for i := 1; i <= 5; i++ {
		nodeID := fmt.Sprintf("node%d", i)
		success := i != 3 // Make node3 fail
		latency := 50 * time.Millisecond
		if i == 2 || i == 4 {
			latency = 350 * time.Millisecond // Higher latency for some nodes
		}
		cap.UpdateNetworkTelemetry(nodeID, latency, success)
	}

	// Check the consistency level (should adapt to deteriorating conditions)
	fmt.Printf("Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("Optimal consistency level: %s\n", cap.GetConsistencyStatus())
	fmt.Printf("Dynamic timeout for 1s base: %v\n", cap.DynamicTimeout(time.Second))

	// Stage 3: Network partition scenario
	fmt.Println("\n[Stage 3] Network partition scenario:")
	// Simulate severe network issues (high latency, many failures)
	for i := 1; i <= 5; i++ {
		nodeID := fmt.Sprintf("node%d", i)
		success := i <= 2                 // Only first two nodes responding
		latency := 800 * time.Millisecond // High latency
		cap.UpdateNetworkTelemetry(nodeID, latency, success)
	}

	// Check the consistency level (should prioritize availability now)
	fmt.Printf("Partition probability: %.2f\n", cap.networkState.PartitionProbability)
	fmt.Printf("Optimal consistency level: %s\n", cap.GetConsistencyStatus())
	fmt.Printf("Dynamic timeout for 1s base: %v\n", cap.DynamicTimeout(time.Second))

	// Stage 4: Test conflict detection and resolution
	fmt.Println("\n[Stage 4] Conflict detection and resolution:")

	// Create two conflicting state updates
	stateID := "block:1000:transaction:5"

	// First update from node1
	vClock1 := cap.RegisterStateUpdate(stateID, "node1", []byte("value1"))
	sv1 := StateVersion{
		Value:      []byte("value1"),
		UpdateTime: time.Now(),
		NodeID:     "node1",
		VClock:     vClock1,
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Concurrent update from node2
	vClock2 := cap.RegisterStateUpdate(stateID, "node2", []byte("value2"))
	sv2 := StateVersion{
		Value:      []byte("value2"),
		UpdateTime: time.Now(),
		NodeID:     "node2",
		VClock:     vClock2,
	}

	// Check if there's a conflict
	conflict := cap.DetectConflict(stateID, vClock1, vClock2)
	fmt.Printf("Conflict detected: %v\n", conflict)

	if conflict {
		// Register the conflict
		cap.RegisterConflict(stateID, []StateVersion{sv1, sv2})

		// Resolve the conflict
		resolved := cap.ResolveConflict(stateID)
		fmt.Printf("Conflict resolved, winner: Node %s (value: %s)\n",
			resolved.NodeID, string(resolved.Value))
	}

	// Stage 5: Integrate with blockchain operations
	fmt.Println("\n[Stage 5] Blockchain integration:")

	// Example: Get appropriate consistency level for next block commit
	consistencyLevel := cap.GetOptimalConsistencyLevel()
	fmt.Printf("Current network state suggests using: %s\n", cap.GetConsistencyStatus())

	// Example of using consistency level for determining consensus requirements
	consensusNodesRequired := 0
	switch consistencyLevel {
	case StrongConsistency:
		consensusNodesRequired = 5 // All nodes must agree
	case CausalConsistency:
		consensusNodesRequired = 3 // Majority must agree
	case EventualConsistency:
		consensusNodesRequired = 1 // Any node can proceed
	}
	fmt.Printf("Required consensus nodes: %d\n", consensusNodesRequired)

	// Example: Adjust transaction processing timeout based on network conditions
	baseTimeout := 2 * time.Second
	adjustedTimeout := cap.DynamicTimeout(baseTimeout)
	fmt.Printf("Adjusted transaction processing timeout: %v (base was %v)\n",
		adjustedTimeout, baseTimeout)

	// Complete the demo
	fmt.Println("\n=== CAP OPTIMIZER DEMO COMPLETE ===")
	fmt.Println("The CAP implementation is now ready to be integrated with blockchain consensus and transactions.")
}
