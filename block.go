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
	Metadata     string
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

// package main

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"errors"
// 	"sync"
// )

// // ShardInfo stores metadata about a shard
// type ShardInfo struct {
// 	ID                int
// 	Size              int
// 	ComputationalLoad float64
// 	ActiveNodes       int
// 	LastRebalance     int64 // Unix timestamp
// }

// // AdaptiveMerkleForest represents the AMF structure
// type AdaptiveMerkleForest struct {
// 	Shards            map[int]*MerkleTree
// 	ShardInfo         map[int]*ShardInfo
// 	MaxShardSize      int
// 	MinShardSize      int
// 	LoadThresholdHigh float64
// 	LoadThresholdLow  float64
// 	mutex             sync.RWMutex
// }

// // NewAdaptiveMerkleForest creates a new AMF instance
// func NewAdaptiveMerkleForest(maxShardSize, minShardSize int) *AdaptiveMerkleForest {
// 	return &AdaptiveMerkleForest{
// 		Shards:            make(map[int]*MerkleTree),
// 		ShardInfo:         make(map[int]*ShardInfo),
// 		MaxShardSize:      maxShardSize,
// 		MinShardSize:      minShardSize,
// 		LoadThresholdHigh: 0.8, // 80% load triggers split
// 		LoadThresholdLow:  0.2, // 20% load triggers merge
// 		mutex:             sync.RWMutex{},
// 	}
// }

// // AddShard adds a new shard to the forest
// func (amf *AdaptiveMerkleForest) AddShard(shardID int, dataBlocks [][]byte) error {
// 	amf.mutex.Lock()
// 	defer amf.mutex.Unlock()

// 	// Check if shard already exists
// 	if _, exists := amf.Shards[shardID]; exists {
// 		return errors.New("shard ID already exists")
// 	}

// 	// Create a new Merkle Tree for this shard
// 	tree, err := NewMerkleTree(dataBlocks)
// 	if err != nil {
// 		return err
// 	}

// 	// Add the tree to the forest
// 	amf.Shards[shardID] = tree

// 	// Create shard info
// 	amf.ShardInfo[shardID] = &ShardInfo{
// 		ID:                shardID,
// 		Size:              len(dataBlocks),
// 		ComputationalLoad: float64(len(dataBlocks)) / float64(amf.MaxShardSize),
// 		ActiveNodes:       len(dataBlocks),
// 		LastRebalance:     0, // Will be set during first rebalance
// 	}

// 	return nil
// }

// // SplitShard splits a shard into two when it exceeds the threshold
// func (amf *AdaptiveMerkleForest) SplitShard(shardID int) (int, int, error) {
// 	amf.mutex.Lock()
// 	defer amf.mutex.Unlock()

// 	// Check if shard exists
// 	originalTree, exists := amf.Shards[shardID]
// 	if !exists {
// 		return 0, 0, errors.New("shard not found")
// 	}

// 	// Check if split is necessary
// 	shardInfo := amf.ShardInfo[shardID]
// 	if shardInfo.Size <= amf.MaxShardSize && shardInfo.ComputationalLoad < amf.LoadThresholdHigh {
// 		return shardID, shardID, nil // No split needed
// 	}

// 	// Generate new shard IDs
// 	newShardID1 := shardID
// 	newShardID2 := len(amf.Shards) + 1

// 	// Split leaf nodes into two groups
// 	midpoint := len(originalTree.LeafNodes) / 2

// 	// Create data blocks for the new shards
// 	dataBlocks1 := make([][]byte, midpoint)
// 	dataBlocks2 := make([][]byte, len(originalTree.LeafNodes)-midpoint)

// 	for i := 0; i < midpoint; i++ {
// 		dataBlocks1[i] = originalTree.LeafNodes[i].Data
// 	}

// 	for i := midpoint; i < len(originalTree.LeafNodes); i++ {
// 		dataBlocks2[i-midpoint] = originalTree.LeafNodes[i].Data
// 	}

// 	// Create new Merkle Trees
// 	tree1, err := NewMerkleTree(dataBlocks1)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	tree2, err := NewMerkleTree(dataBlocks2)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	// Update the forest
// 	delete(amf.Shards, shardID)
// 	delete(amf.ShardInfo, shardID)

// 	amf.Shards[newShardID1] = tree1
// 	amf.Shards[newShardID2] = tree2

// 	// Update shard info
// 	amf.ShardInfo[newShardID1] = &ShardInfo{
// 		ID:                newShardID1,
// 		Size:              len(dataBlocks1),
// 		ComputationalLoad: float64(len(dataBlocks1)) / float64(amf.MaxShardSize),
// 		ActiveNodes:       len(dataBlocks1),
// 		LastRebalance:     0,
// 	}

// 	amf.ShardInfo[newShardID2] = &ShardInfo{
// 		ID:                newShardID2,
// 		Size:              len(dataBlocks2),
// 		ComputationalLoad: float64(len(dataBlocks2)) / float64(amf.MaxShardSize),
// 		ActiveNodes:       len(dataBlocks2),
// 		LastRebalance:     0,
// 	}

// 	return newShardID1, newShardID2, nil
// }

// // MergeShards combines two underutilized shards
// func (amf *AdaptiveMerkleForest) MergeShards(shardID1, shardID2 int) (int, error) {
// 	amf.mutex.Lock()
// 	defer amf.mutex.Unlock()

// 	// Check if shards exist
// 	tree1, exists1 := amf.Shards[shardID1]
// 	tree2, exists2 := amf.Shards[shardID2]

// 	if !exists1 || !exists2 {
// 		return 0, errors.New("one or both shards not found")
// 	}

// 	// Check if merge is beneficial
// 	info1 := amf.ShardInfo[shardID1]
// 	info2 := amf.ShardInfo[shardID2]

// 	combinedSize := info1.Size + info2.Size
// 	if combinedSize > amf.MaxShardSize {
// 		return 0, errors.New("merged shard would exceed maximum size")
// 	}

// 	// Combine data blocks
// 	combinedBlocks := make([][]byte, 0, combinedSize)

// 	for _, node := range tree1.LeafNodes {
// 		combinedBlocks = append(combinedBlocks, node.Data)
// 	}

// 	for _, node := range tree2.LeafNodes {
// 		combinedBlocks = append(combinedBlocks, node.Data)
// 	}

// 	// Create new merged tree
// 	newTree, err := NewMerkleTree(combinedBlocks)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Create new shard ID for the merged shard
// 	newShardID := len(amf.Shards) + 1

// 	// Update the forest
// 	delete(amf.Shards, shardID1)
// 	delete(amf.Shards, shardID2)
// 	delete(amf.ShardInfo, shardID1)
// 	delete(amf.ShardInfo, shardID2)

// 	amf.Shards[newShardID] = newTree
// 	amf.ShardInfo[newShardID] = &ShardInfo{
// 		ID:                newShardID,
// 		Size:              combinedSize,
// 		ComputationalLoad: float64(combinedSize) / float64(amf.MaxShardSize),
// 		ActiveNodes:       combinedSize,
// 		LastRebalance:     0,
// 	}

// 	return newShardID, nil
// }

// // GenerateCrossShardProof creates a proof that spans across two shards
// func (amf *AdaptiveMerkleForest) GenerateCrossShardProof(shardID1, shardID2 int, data []byte) ([][]byte, error) {
// 	amf.mutex.RLock()
// 	defer amf.mutex.RUnlock()

// 	// Verify shards exist
// 	tree1, exists1 := amf.Shards[shardID1]
// 	tree2, exists2 := amf.Shards[shardID2]

// 	if !exists1 || !exists2 {
// 		return nil, errors.New("one or both shards not found")
// 	}

// 	// Try to generate proof from first shard
// 	proof1, err1 := tree1.GenerateMerkleProof(data)
// 	if err1 == nil {
// 		// Data found in shard1, include root hash of shard2 for cross-verification
// 		rootHash2 := tree2.Root.Hash
// 		proof1 = append(proof1, rootHash2)
// 		return proof1, nil
// 	}

// 	// Try to generate proof from second shard
// 	proof2, err2 := tree2.GenerateMerkleProof(data)
// 	if err2 == nil {
// 		// Data found in shard2, include root hash of shard1 for cross-verification
// 		rootHash1 := tree1.Root.Hash
// 		proof2 = append(proof2, rootHash1)
// 		return proof2, nil
// 	}

// 	return nil, errors.New("data not found in either shard")
// }

// // GetCompressedProof generates a probabilistic compressed proof
// // This is a simplified implementation of AMQ (Approximate Membership Query) filters
// func (amf *AdaptiveMerkleForest) GetCompressedProof(shardID int, data []byte) ([]byte, error) {
// 	amf.mutex.RLock()
// 	defer amf.mutex.RUnlock()

// 	tree, exists := amf.Shards[shardID]
// 	if !exists {
// 		return nil, errors.New("shard not found")
// 	}

// 	// First check if data exists
// 	if !tree.VerifyData(data) {
// 		return nil, errors.New("data not found in specified shard")
// 	}

// 	// Generate standard proof
// 	proof, err := tree.GenerateMerkleProof(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Compress the proof (simplified version)
// 	// In a real implementation, this would use Bloom filters or other AMQ techniques
// 	var compressed []byte

// 	// Create a hash of all proof elements
// 	for _, element := range proof {
// 		compressed = append(compressed, element...)
// 	}

// 	// Apply final compression hash
// 	hash := sha256.Sum256(compressed)
// 	return hash[:], nil
// }

// // VerifyCompressedProof verifies a compressed proof
// func (amf *AdaptiveMerkleForest) VerifyCompressedProof(shardID int, data []byte, compressedProof []byte) bool {
// 	amf.mutex.RLock()
// 	defer amf.mutex.RUnlock()

// 	tree, exists := amf.Shards[shardID]
// 	if !exists {
// 		return false
// 	}

// 	// Generate standard proof for comparison
// 	proof, err := tree.GenerateMerkleProof(data)
// 	if err != nil {
// 		return false
// 	}

// 	// Apply the same compression algorithm
// 	var compressed []byte
// 	for _, element := range proof {
// 		compressed = append(compressed, element...)
// 	}

// 	hash := sha256.Sum256(compressed)
// 	return hex.EncodeToString(hash[:]) == hex.EncodeToString(compressedProof)
// }

// // GetShardRootHashes returns all shard root hashes for the forest
// func (amf *AdaptiveMerkleForest) GetShardRootHashes() map[int]string {
// 	amf.mutex.RLock()
// 	defer amf.mutex.RUnlock()

// 	rootHashes := make(map[int]string)
// 	for id, tree := range amf.Shards {
// 		rootHashes[id] = hex.EncodeToString(tree.Root.Hash)
// 	}

// 	return rootHashes
// }
