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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"
)

// ===== Init =====

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// ===== Key Functions =====

func GenerateKeyPair() (*ecdsa.PrivateKey, string, *ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, "", nil, err
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, "", nil, err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	return priv, string(pubPEM), &priv.PublicKey, nil
}

func SignBlock(priv *ecdsa.PrivateKey, block Block) ([]byte, []byte, error) {
	hash := sha256.Sum256([]byte(block.Hash))
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return nil, nil, err
	}
	return r.Bytes(), s.Bytes(), nil
}

func VerifySignature(pub *ecdsa.PublicKey, block Block, rBytes, sBytes []byte) bool {
	hash := sha256.Sum256([]byte(block.Hash))
	var r, s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)
	return ecdsa.Verify(pub, hash[:], &r, &s)
}

func SimulateConsensusVote(validators []Validator, block Block) bool {
	yesVotes := 0
	fmt.Println("\n=== VALIDATOR VOTING ===")

	for _, v := range validators {
		// Simulate signature generation
		r, s, err := SignBlock(v.PrivateKey, block)
		if err != nil {
			fmt.Printf("❌ %s: Failed to sign block - %v\n", v.ID, err)
			continue
		}

		// Simulate signature verification
		isValid := VerifySignature(v.PublicKeyObj, block, r, s)
		if isValid {
			fmt.Printf("✅ %s: Voted YES (Signature Verified)\n", v.ID)
			yesVotes++
		} else {
			fmt.Printf("❌ %s: Voted NO (Signature Invalid)\n", v.ID)
		}
	}

	requiredVotes := len(validators)*2/3 + 1
	fmt.Printf("🧮 YES Votes: %d / %d (Required: %d)\n", yesVotes, len(validators), requiredVotes)

	return yesVotes >= requiredVotes
}

func simulateValidatorVote(v Validator, block Block) bool {
	r, s, err := SignBlock(v.PrivateKey, block)
	if err != nil {
		fmt.Printf("  Error signing by %s: %v\n", v.ID, err)
		return false
	}
	valid := VerifySignature(v.PublicKeyObj, block, r, s)
	if !valid {
		fmt.Printf("  ❌ Invalid signature by %s\n", v.ID)
	}
	return valid
}

func createTransactions() []Transaction {
	return []Transaction{
		{Sender: "Alice", Recipient: "Bob", Amount: 10.0, Data: "Payment", Timestamp: time.Now().String(), Signature: "sig1"},
		{Sender: "Charlie", Recipient: "Dave", Amount: 5.5, Data: "Loan", Timestamp: time.Now().String(), Signature: "sig2"},
	}
}

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

// ===== Main =====

func main() {
	fmt.Println("\n=== VALIDATOR CONSENSUS DEMO ===")

	var validators []Validator
	for i := 1; i <= 4; i++ {
		priv, pubStr, pubObj, err := GenerateKeyPair()
		if err != nil {
			fmt.Printf("Error generating key for validator %d: %v\n", i, err)
			continue
		}
		validators = append(validators, Validator{
			ID:           fmt.Sprintf("Validator%d", i),
			PublicKey:    pubStr,
			PrivateKey:   priv,
			PublicKeyObj: pubObj,
			Power:        1,
		})
	}

	var blockchain []Block
	genesisBlock := generateGenesisBlock()
	blockchain = append(blockchain, genesisBlock)

	amf := NewAdaptiveMerkleForest(100, 10)
	_ = amf.AddShard(0, [][]byte{[]byte(fmt.Sprintf("%v", genesisBlock.Transactions))})

	newTransactions := createTransactions()
	difficulty := 1
	newBlock := generateNextBlock(blockchain[len(blockchain)-1], newTransactions)
	newBlock.mineBlock(difficulty)

	if SimulateConsensusVote(validators, newBlock) {
		fmt.Println("✅ Consensus Reached: Block Approved")
		blockchain = append(blockchain, newBlock)

		var txData [][]byte
		for _, tx := range newTransactions {
			txData = append(txData, []byte(fmt.Sprintf("%v", tx)))
		}
		_ = amf.AddShard(1, txData)
	} else {
		fmt.Println("❌ Consensus Failed: Block Rejected")
	}

	fmt.Println("\n=== BLOCKCHAIN STATE ===")
	for _, block := range blockchain {
		printBlockInfo(block)
	}

	fmt.Println("\n=== ADAPTIVE MERKLE FOREST STATE ===")
	for shardID, rootHash := range amf.GetShardRootHashes() {
		fmt.Printf("  Shard #%d Root Hash: %s\n", shardID, rootHash)
	}

	// Demo: Cross-shard & compressed proof
	fmt.Println("\n=== PROOF DEMO ===")
	txData := []byte(fmt.Sprintf("%v", newTransactions[0]))

	crossProof, err := amf.GenerateCrossShardProof(0, 1, txData)
	if err != nil {
		fmt.Printf("  ❌ Cross-shard proof error: %v\n", err)
	} else {
		fmt.Printf("  ✅ Cross-shard proof generated (%d elements)\n", len(crossProof))
	}

	compressed, err := amf.GetCompressedProof(1, txData)
	if err != nil {
		fmt.Printf("  ❌ Compressed proof error: %v\n", err)
	} else {
		fmt.Printf("  ✅ Compressed proof: %s\n", hex.EncodeToString(compressed))
		if amf.VerifyCompressedProof(1, txData, compressed) {
			fmt.Println("  ✅ Compressed proof verification: VALID")
		} else {
			fmt.Println("  ❌ Compressed proof verification: INVALID")
		}
	}
}
