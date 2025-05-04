package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

// ShardInfo stores metadata about a shard
type ShardInfo struct {
	ID                int
	Size              int
	ComputationalLoad float64
	ActiveNodes       int
	LastRebalance     int64 // Unix timestamp
}

// AdaptiveMerkleForest represents the AMF structure
type AdaptiveMerkleForest struct {
	Shards            map[int]*MerkleTree
	ShardInfo         map[int]*ShardInfo
	MaxShardSize      int
	MinShardSize      int
	LoadThresholdHigh float64
	LoadThresholdLow  float64
	mutex             sync.RWMutex
}

// NewAdaptiveMerkleForest creates a new AMF instance
func NewAdaptiveMerkleForest(maxShardSize, minShardSize int) *AdaptiveMerkleForest {
	return &AdaptiveMerkleForest{
		Shards:            make(map[int]*MerkleTree),
		ShardInfo:         make(map[int]*ShardInfo),
		MaxShardSize:      maxShardSize,
		MinShardSize:      minShardSize,
		LoadThresholdHigh: 0.8, // 80% load triggers split
		LoadThresholdLow:  0.2, // 20% load triggers merge
		mutex:             sync.RWMutex{},
	}
}

// AddShard adds a new shard to the forest
func (amf *AdaptiveMerkleForest) AddShard(shardID int, dataBlocks [][]byte) error {
	amf.mutex.Lock()
	defer amf.mutex.Unlock()

	// Check if shard already exists
	if _, exists := amf.Shards[shardID]; exists {
		return errors.New("shard ID already exists")
	}

	// Create a new Merkle Tree for this shard
	tree, err := NewMerkleTree(dataBlocks)
	if err != nil {
		return err
	}

	// Add the tree to the forest
	amf.Shards[shardID] = tree

	// Create shard info
	amf.ShardInfo[shardID] = &ShardInfo{
		ID:                shardID,
		Size:              len(dataBlocks),
		ComputationalLoad: float64(len(dataBlocks)) / float64(amf.MaxShardSize),
		ActiveNodes:       len(dataBlocks),
		LastRebalance:     0, // Will be set during first rebalance
	}

	return nil
}

// SplitShard splits a shard into two when it exceeds the threshold
func (amf *AdaptiveMerkleForest) SplitShard(shardID int) (int, int, error) {
	amf.mutex.Lock()
	defer amf.mutex.Unlock()

	// Check if shard exists
	originalTree, exists := amf.Shards[shardID]
	if !exists {
		return 0, 0, errors.New("shard not found")
	}

	// Check if split is necessary
	shardInfo := amf.ShardInfo[shardID]
	if shardInfo.Size <= amf.MaxShardSize && shardInfo.ComputationalLoad < amf.LoadThresholdHigh {
		return shardID, shardID, nil // No split needed
	}

	// Generate new shard IDs
	newShardID1 := shardID
	newShardID2 := len(amf.Shards) + 1

	// Split leaf nodes into two groups
	midpoint := len(originalTree.LeafNodes) / 2

	// Create data blocks for the new shards
	dataBlocks1 := make([][]byte, midpoint)
	dataBlocks2 := make([][]byte, len(originalTree.LeafNodes)-midpoint)

	for i := 0; i < midpoint; i++ {
		dataBlocks1[i] = originalTree.LeafNodes[i].Data
	}

	for i := midpoint; i < len(originalTree.LeafNodes); i++ {
		dataBlocks2[i-midpoint] = originalTree.LeafNodes[i].Data
	}

	// Create new Merkle Trees
	tree1, err := NewMerkleTree(dataBlocks1)
	if err != nil {
		return 0, 0, err
	}

	tree2, err := NewMerkleTree(dataBlocks2)
	if err != nil {
		return 0, 0, err
	}

	// Update the forest
	delete(amf.Shards, shardID)
	delete(amf.ShardInfo, shardID)

	amf.Shards[newShardID1] = tree1
	amf.Shards[newShardID2] = tree2

	// Update shard info
	amf.ShardInfo[newShardID1] = &ShardInfo{
		ID:                newShardID1,
		Size:              len(dataBlocks1),
		ComputationalLoad: float64(len(dataBlocks1)) / float64(amf.MaxShardSize),
		ActiveNodes:       len(dataBlocks1),
		LastRebalance:     0,
	}

	amf.ShardInfo[newShardID2] = &ShardInfo{
		ID:                newShardID2,
		Size:              len(dataBlocks2),
		ComputationalLoad: float64(len(dataBlocks2)) / float64(amf.MaxShardSize),
		ActiveNodes:       len(dataBlocks2),
		LastRebalance:     0,
	}

	return newShardID1, newShardID2, nil
}

// MergeShards combines two underutilized shards
func (amf *AdaptiveMerkleForest) MergeShards(shardID1, shardID2 int) (int, error) {
	amf.mutex.Lock()
	defer amf.mutex.Unlock()

	// Check if shards exist
	tree1, exists1 := amf.Shards[shardID1]
	tree2, exists2 := amf.Shards[shardID2]

	if !exists1 || !exists2 {
		return 0, errors.New("one or both shards not found")
	}

	// Check if merge is beneficial
	info1 := amf.ShardInfo[shardID1]
	info2 := amf.ShardInfo[shardID2]

	combinedSize := info1.Size + info2.Size
	if combinedSize > amf.MaxShardSize {
		return 0, errors.New("merged shard would exceed maximum size")
	}

	// Combine data blocks
	combinedBlocks := make([][]byte, 0, combinedSize)

	for _, node := range tree1.LeafNodes {
		combinedBlocks = append(combinedBlocks, node.Data)
	}

	for _, node := range tree2.LeafNodes {
		combinedBlocks = append(combinedBlocks, node.Data)
	}

	// Create new merged tree
	newTree, err := NewMerkleTree(combinedBlocks)
	if err != nil {
		return 0, err
	}

	// Create new shard ID for the merged shard
	newShardID := len(amf.Shards) + 1

	// Update the forest
	delete(amf.Shards, shardID1)
	delete(amf.Shards, shardID2)
	delete(amf.ShardInfo, shardID1)
	delete(amf.ShardInfo, shardID2)

	amf.Shards[newShardID] = newTree
	amf.ShardInfo[newShardID] = &ShardInfo{
		ID:                newShardID,
		Size:              combinedSize,
		ComputationalLoad: float64(combinedSize) / float64(amf.MaxShardSize),
		ActiveNodes:       combinedSize,
		LastRebalance:     0,
	}

	return newShardID, nil
}

// GenerateCrossShardProof creates a proof that spans across two shards
func (amf *AdaptiveMerkleForest) GenerateCrossShardProof(shardID1, shardID2 int, data []byte) ([][]byte, error) {
	amf.mutex.RLock()
	defer amf.mutex.RUnlock()

	// Verify shards exist
	tree1, exists1 := amf.Shards[shardID1]
	tree2, exists2 := amf.Shards[shardID2]

	if !exists1 || !exists2 {
		return nil, errors.New("one or both shards not found")
	}

	// Try to generate proof from first shard
	proof1, err1 := tree1.GenerateMerkleProof(data)
	if err1 == nil {
		// Data found in shard1, include root hash of shard2 for cross-verification
		rootHash2 := tree2.Root.Hash
		proof1 = append(proof1, rootHash2)
		return proof1, nil
	}

	// Try to generate proof from second shard
	proof2, err2 := tree2.GenerateMerkleProof(data)
	if err2 == nil {
		// Data found in shard2, include root hash of shard1 for cross-verification
		rootHash1 := tree1.Root.Hash
		proof2 = append(proof2, rootHash1)
		return proof2, nil
	}

	return nil, errors.New("data not found in either shard")
}

// GetCompressedProof generates a probabilistic compressed proof
// This is a simplified implementation of AMQ (Approximate Membership Query) filters
func (amf *AdaptiveMerkleForest) GetCompressedProof(shardID int, data []byte) ([]byte, error) {
	amf.mutex.RLock()
	defer amf.mutex.RUnlock()

	tree, exists := amf.Shards[shardID]
	if !exists {
		return nil, errors.New("shard not found")
	}

	// First check if data exists
	if !tree.VerifyData(data) {
		return nil, errors.New("data not found in specified shard")
	}

	// Generate standard proof
	proof, err := tree.GenerateMerkleProof(data)
	if err != nil {
		return nil, err
	}

	// Compress the proof (simplified version)
	// In a real implementation, this would use Bloom filters or other AMQ techniques
	var compressed []byte

	// Create a hash of all proof elements
	for _, element := range proof {
		compressed = append(compressed, element...)
	}

	// Apply final compression hash
	hash := sha256.Sum256(compressed)
	return hash[:], nil
}

// VerifyCompressedProof verifies a compressed proof
func (amf *AdaptiveMerkleForest) VerifyCompressedProof(shardID int, data []byte, compressedProof []byte) bool {
	amf.mutex.RLock()
	defer amf.mutex.RUnlock()

	tree, exists := amf.Shards[shardID]
	if !exists {
		return false
	}

	// Generate standard proof for comparison
	proof, err := tree.GenerateMerkleProof(data)
	if err != nil {
		return false
	}

	// Apply the same compression algorithm
	var compressed []byte
	for _, element := range proof {
		compressed = append(compressed, element...)
	}

	hash := sha256.Sum256(compressed)
	return hex.EncodeToString(hash[:]) == hex.EncodeToString(compressedProof)
}

// GetShardRootHashes returns all shard root hashes for the forest
func (amf *AdaptiveMerkleForest) GetShardRootHashes() map[int]string {
	amf.mutex.RLock()
	defer amf.mutex.RUnlock()

	rootHashes := make(map[int]string)
	for id, tree := range amf.Shards {
		rootHashes[id] = hex.EncodeToString(tree.Root.Hash)
	}

	return rootHashes
}
