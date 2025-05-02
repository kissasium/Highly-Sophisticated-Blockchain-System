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
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// type of validator
type Validator struct {
	ID           string
	PublicKey    string
	PrivateKey   *ecdsa.PrivateKey
	PublicKeyObj *ecdsa.PublicKey
	Power        int
}

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
	ShardID      int   // For AMF integration
	Nonce        int64 // Nonce for Proof of Work
}

// calculateHash generates a SHA-256 hash for the block
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s%d", block.Index, block.Timestamp, block.MerkleRoot, block.PrevHash, block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
func generateBlockHash(block Block) string {
	// Convert block fields to a string (the order of fields is important)
	blockData := fmt.Sprintf("%d%s%s%d%v",
		block.Index, block.Timestamp, block.PrevHash, block.Nonce, block.Transactions)

	// Generate hash using SHA-256
	hash := sha256.Sum256([]byte(blockData))

	// Return the resulting hash as a string
	return fmt.Sprintf("%x", hash)
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
func isValidBlockHash(hash string) bool {
	// For example, check if the hash starts with "0000"
	return hash[:4] == "0000"
}

// generateNextBlock creates the next block using the previous block
// generateNextBlock creates the next block using the previous block, WITHOUT mining
func generateNextBlock(previousBlock Block, transactions []Transaction) Block {
	newIndex := previousBlock.Index + 1
	newTimestamp := time.Now().String()
	previousHash := previousBlock.Hash

	newBlock := Block{
		Index:        newIndex,
		Timestamp:    newTimestamp,
		PrevHash:     previousHash,
		Transactions: transactions,
		Nonce:        0, // start with nonce 0
	}

	// Create Merkle Tree and Root
	var txData [][]byte
	for _, tx := range transactions {
		txData = append(txData, []byte(fmt.Sprintf("%v", tx)))
	}
	merkleTree, _ := NewMerkleTree(txData)
	newBlock.MerkleTree = merkleTree
	newBlock.MerkleRoot = merkleTree.GetRootHash()

	// Initial hash (with nonce = 0, usually invalid)
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
	sb.WriteString(fmt.Sprintf("Nonce: %d\n", block.Nonce))
	sb.WriteString(fmt.Sprintf("Transactions: %d\n", len(block.Transactions)))

	return sb.String()
}

// mineBlock attempts to find the correct nonce for the block to satisfy the difficulty target
func (block *Block) mineBlock(difficulty int) {
	prefix := strings.Repeat("0", difficulty) // Difficulty level: Number of leading zeros
	for {
		block.Nonce++
		block.Hash = calculateHash(*block)
		if strings.HasPrefix(block.Hash, prefix) {
			fmt.Printf("Block mined! Nonce: %d\n", block.Nonce)
			break
		}
	}
}

// AuthenticateValidator checks if a validator is authorized (basic simulation)
func AuthenticateValidator(validator Validator, message, signature string) bool {
	// In a real blockchain, you would verify the signature against the public key.
	// Here, we'll simulate it by checking if the signature is just "VALID:"+validator.ID
	expectedSignature := "VALID:" + validator.ID
	return signature == expectedSignature
}

// ValidatorReputation stores reputation scores for validators
var ValidatorReputation = make(map[string]float64)

// UpdateValidatorReputation adjusts reputation based on behavior
func UpdateValidatorReputation(validatorID string, delta float64) {
	ValidatorReputation[validatorID] += delta
	if ValidatorReputation[validatorID] < 0 {
		ValidatorReputation[validatorID] = 0 // No negative reputation
	}
}

// GetValidatorReputation returns the current reputation score
func GetValidatorReputation(validatorID string) float64 {
	return ValidatorReputation[validatorID]
}

// HandleValidatorAction updates reputation based on whether the validator acted correctly
func HandleValidatorAction(validatorID string, success bool) {
	if success {
		UpdateValidatorReputation(validatorID, 1.0) // Reward good behavior
	} else {
		UpdateValidatorReputation(validatorID, -2.0) // Penalize bad behavior more harshly
	}
}
