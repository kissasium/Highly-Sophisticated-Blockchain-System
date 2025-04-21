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
}
