// package main

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"time"
// )

// // Block represents a single block in the blockchain
// type Block struct {
// 	Index     int
// 	Timestamp string
// 	Data      string
// 	PrevHash  string
// 	Hash      string
// }

// // calculateHash generates a SHA-256 hash for the block
// func calculateHash(block Block) string {
// 	record := fmt.Sprintf("%d%s%s%s", block.Index, block.Timestamp, block.Data, block.PrevHash)
// 	hash := sha256.New()
// 	hash.Write([]byte(record))
// 	hashed := hash.Sum(nil)
// 	return hex.EncodeToString(hashed)
// }

// // generateGenesisBlock creates the first block in the blockchain
// func generateGenesisBlock() Block {
// 	genesis := Block{0, time.Now().String(), "Genesis Block", "", ""}
// 	genesis.Hash = calculateHash(genesis)
// 	return genesis
// }

// // generateNextBlock creates the next block using the previous block
// func generateNextBlock(prevBlock Block, data string) Block {
// 	newBlock := Block{
// 		Index:     prevBlock.Index + 1,
// 		Timestamp: time.Now().String(),
// 		Data:      data,
// 		PrevHash:  prevBlock.Hash,
// 	}
// 	newBlock.Hash = calculateHash(newBlock)
// 	return newBlock
// }

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Transaction represents a single transaction in the block
type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	Data      string
	Timestamp string
	Signature string
}

// Block represents a single block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	MerkleRoot   string      // Root hash of the Merkle tree
	MerkleTree   *MerkleTree // Merkle tree of transactions
	PrevHash     string
	Hash         string
	ShardID      int // For AMF integration
}

// calculateHash generates a SHA-256 hash for the block
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.Index, block.Timestamp, block.MerkleRoot, block.PrevHash)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// generateGenesisBlock creates the first block in the blockchain
func generateGenesisBlock() Block {
	// Create a genesis transaction
	genesisTx := Transaction{
		Sender:    "Genesis",
		Recipient: "System",
		Amount:    0,
		Data:      "Genesis Transaction",
		Timestamp: time.Now().String(),
		Signature: "0",
	}

	// Convert transaction to data for Merkle Tree
	txData := [][]byte{[]byte(fmt.Sprintf("%v", genesisTx))}

	// Create Merkle Tree
	merkleTree, _ := NewMerkleTree(txData)
	merkleRoot := merkleTree.GetRootHash()

	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{genesisTx},
		MerkleRoot:   merkleRoot,
		MerkleTree:   merkleTree,
		PrevHash:     "",
		ShardID:      0,
	}

	genesis.Hash = calculateHash(genesis)
	return genesis
}

// generateNextBlock creates the next block using the previous block
func generateNextBlock(prevBlock Block, transactions []Transaction) Block {
	// Convert transactions to data for Merkle Tree
	var txData [][]byte
	for _, tx := range transactions {
		txData = append(txData, []byte(fmt.Sprintf("%v", tx)))
	}

	// Create Merkle Tree
	merkleTree, _ := NewMerkleTree(txData)
	merkleRoot := merkleTree.GetRootHash()

	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		MerkleRoot:   merkleRoot,
		MerkleTree:   merkleTree,
		PrevHash:     prevBlock.Hash,
		ShardID:      prevBlock.ShardID, // Inherit shard ID by default
	}

	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

// VerifyBlockIntegrity checks if the block's data is valid and consistent
func VerifyBlockIntegrity(block Block) bool {
	// Verify Merkle Root
	var txData [][]byte
	for _, tx := range block.Transactions {
		txData = append(txData, []byte(fmt.Sprintf("%v", tx)))
	}

	// Recreate Merkle Tree
	merkleTree, err := NewMerkleTree(txData)
	if err != nil {
		return false
	}

	calculatedRoot := merkleTree.GetRootHash()
	if calculatedRoot != block.MerkleRoot {
		return false
	}

	// Verify block hash
	calculatedHash := calculateHash(block)
	return calculatedHash == block.Hash
}

// GetTransactionProof generates a Merkle proof for a specific transaction
func (block Block) GetTransactionProof(tx Transaction) ([][]byte, error) {
	txData := []byte(fmt.Sprintf("%v", tx))
	return block.MerkleTree.GenerateMerkleProof(txData)
}

// String provides a string representation of the block
func (block Block) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Block #%d\n", block.Index))
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n", block.Timestamp))
	sb.WriteString(fmt.Sprintf("MerkleRoot: %s\n", block.MerkleRoot))
	sb.WriteString(fmt.Sprintf("PrevHash: %s\n", block.Hash))
	sb.WriteString(fmt.Sprintf("Hash: %s\n", block.Hash))
	sb.WriteString(fmt.Sprintf("ShardID: %d\n", block.ShardID))
	sb.WriteString(fmt.Sprintf("Transactions: %d\n", len(block.Transactions)))

	return sb.String()
}
