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
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ===== Helper Functions =====

// SimulateConsensusVote simulates validators voting on a block
func SimulateConsensusVote(validators []Validator, block Block) bool {
	votes := 0
	total := len(validators)

	for _, v := range validators {
		if simulateValidatorVote(v, block, total) {
			votes++
		}
	}
	requiredVotes := (2*total)/total + 1
	fmt.Printf("Consensus votes: %d/%d (required: %d)\n", votes, total, requiredVotes)
	return votes >= requiredVotes
}

// simulateValidatorVote simulates a validator approving a block
func simulateValidatorVote(v Validator, block Block, totalValidators int) bool {
	votes := 0
	requiredVotes := (2*totalValidators)/totalValidators + 1

	for i := 0; i < totalValidators; i++ {
		chance := rand.Intn(10)
		if chance >= 2 {
			votes++
		}
	}

	return votes >= requiredVotes
}

// createTransactions generates demo transactions
func createTransactions() []Transaction {
	return []Transaction{
		{Sender: "Alice", Recipient: "Bob", Amount: 10.0, Data: "Payment", Timestamp: time.Now().String(), Signature: "signature1"},
		{Sender: "Charlie", Recipient: "Dave", Amount: 5.5, Data: "Loan", Timestamp: time.Now().String(), Signature: "signature2"},
	}
}

// printBlockInfo prints detailed block information
func printBlockInfo(block Block) {
	fmt.Printf("\nBlock #%d\n", block.Index)
	fmt.Printf("  Timestamp: %s\n", block.Timestamp)
	fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
	fmt.Printf("  Current Hash: %s\n", block.Hash)
	fmt.Printf("  Merkle Root: %s\n", block.MerkleRoot)
	fmt.Printf("  Number of Transactions: %d\n", len(block.Transactions))
	fmt.Printf("  Nonce: %d\n", block.Nonce)

	if VerifyBlockIntegrity(block) {
		fmt.Println("  Block Integrity: VALID")
	} else {
		fmt.Println("  Block Integrity: INVALID")
	}
}

// ===== Main Function =====

func main() {
	// ========== PART 1: Validator Consensus Simulation ==========

	fmt.Println("\n=== VALIDATOR CONSENSUS DEMO ===")
	validators := []Validator{
		{ID: "Validator1", PublicKey: "pubKey1", Power: 1},
		{ID: "Validator2", PublicKey: "pubKey2", Power: 1},
		{ID: "Validator3", PublicKey: "pubKey3", Power: 1},
		{ID: "Validator4", PublicKey: "pubKey4", Power: 1},
	}

	var blockchain []Block

	// Step 1: Create genesis block
	genesisBlock := generateGenesisBlock()
	blockchain = append(blockchain, genesisBlock)

	// Step 2: Initialize Adaptive Merkle Forest
	amf := NewAdaptiveMerkleForest(100, 10)

	// Step 3: Add genesis block transactions to AMF
	genesisTxData := [][]byte{[]byte(fmt.Sprintf("%v", genesisBlock.Transactions))}
	if err := amf.AddShard(0, genesisTxData); err != nil {
		fmt.Printf("Error adding genesis shard: %v\n", err)
	}

	// Step 4: Create transactions for next block
	newTransactions := createTransactions()

	// Step 5: Generate new block
	difficulty := 1
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], newTransactions)
	newBlock.mineBlock(difficulty)

	// Step 6: Simulate consensus voting
	if SimulateConsensusVote(validators, newBlock) {
		fmt.Println("Consensus Reached: Block Approved ✅")
		blockchain = append(blockchain, newBlock)

		// Add transactions to AMF
		var txData [][]byte
		for _, tx := range newTransactions {
			txData = append(txData, []byte(fmt.Sprintf("%v", tx)))
		}
		if err := amf.AddShard(1, txData); err != nil {
			fmt.Printf("Error adding block 1 shard: %v\n", err)
		}
	} else {
		fmt.Println("Consensus Failed: Block Rejected ❌")
		// Block is discarded
	}

	// Step 7: Display blockchain state
	fmt.Println("\n=== BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}

	// Step 8: Display AMF state
	fmt.Println("\n=== ADAPTIVE MERKLE FOREST STATE ===")
	rootHashes := amf.GetShardRootHashes()
	for shardID, rootHash := range rootHashes {
		fmt.Printf("  Shard #%d Root Hash: %s\n", shardID, rootHash)
	}

	// ========== Adaptive Merkle Forest Operations ==========

	fmt.Println("\n\n===AMF OPERATIONS DEMO ===")

	// Create additional transactions
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

	// Create blocks
	block1 := generateNextBlock(blockchain[len(blockchain)-1], []Transaction{tx1, tx2})
	blockchain = append(blockchain, block1)

	block2 := generateNextBlock(blockchain[len(blockchain)-1], []Transaction{tx3, tx4, tx5})
	blockchain = append(blockchain, block2)

	// Add new transactions to AMF
	transactionsData := [][]byte{
		[]byte(fmt.Sprintf("%v", tx1)),
		[]byte(fmt.Sprintf("%v", tx2)),
	}
	if err := amf.AddShard(2, transactionsData); err != nil {
		fmt.Printf("Error adding to shard 2: %v\n", err)
	}

	transactionsData = [][]byte{
		[]byte(fmt.Sprintf("%v", tx3)),
		[]byte(fmt.Sprintf("%v", tx4)),
		[]byte(fmt.Sprintf("%v", tx5)),
	}
	if err := amf.AddShard(3, transactionsData); err != nil {
		fmt.Printf("Error adding to shard 3: %v\n", err)
	}

	// Demonstrate cross-shard proof
	fmt.Println("\nGenerating Cross-Shard Proof:")
	txData := []byte(fmt.Sprintf("%v", tx1))
	crossProof, err := amf.GenerateCrossShardProof(0, 2, txData)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Cross-shard proof generated: %d elements\n", len(crossProof))
	}

	// Demonstrate compressed proof
	fmt.Println("\nGenerating Compressed Proof:")
	compressedProof, err := amf.GetCompressedProof(2, txData)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Compressed proof: %s\n", hex.EncodeToString(compressedProof))

		// Verify compressed proof
		if amf.VerifyCompressedProof(2, txData, compressedProof) {
			fmt.Println("  Compressed proof verification: VALID")
		} else {
			fmt.Println("  Compressed proof verification: INVALID")
		}
	}
}
