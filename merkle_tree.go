package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
)

// MerkleNode represents a node in the Merkle Tree
type MerkleNode struct {
	Left    *MerkleNode
	Right   *MerkleNode
	Parent  *MerkleNode
	Hash    []byte
	Data    []byte
	IsLeaf  bool
	ShardID int // For AMF sharding
}

// MerkleTree represents a Merkle Tree data structure
type MerkleTree struct {
	Root       *MerkleNode
	LeafNodes  []*MerkleNode
	Depth      int
	TotalNodes int
}

// NewMerkleNode creates a new Merkle Tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		// Leaf node with data
		hash := sha256.Sum256(data)
		node.Hash = hash[:]
		node.Data = data
		node.IsLeaf = true
	} else {
		// Internal node with children
		prevHashes := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(prevHashes)
		node.Hash = hash[:]
		node.Left = left
		node.Right = right
		node.IsLeaf = false

		// Set parent references
		if left != nil {
			left.Parent = &node
		}
		if right != nil {
			right.Parent = &node
		}
	}

	return &node
}

// NewMerkleTree creates a new Merkle Tree from data blocks
func NewMerkleTree(dataBlocks [][]byte) (*MerkleTree, error) {
	if len(dataBlocks) == 0 {
		return nil, errors.New("cannot create Merkle tree with no data")
	}

	// Create leaf nodes
	var leafNodes []*MerkleNode
	for _, data := range dataBlocks {
		node := NewMerkleNode(nil, nil, data)
		leafNodes = append(leafNodes, node)
	}

	// Calculate tree depth
	depth := int(math.Ceil(math.Log2(float64(len(leafNodes)))))

	// Build the tree from bottom up
	nodes := leafNodes
	for len(nodes) > 1 {
		var level []*MerkleNode

		// Process nodes in pairs to create parent nodes
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				// Create parent with two children
				node := NewMerkleNode(nodes[i], nodes[i+1], nil)
				level = append(level, node)
			} else {
				// Odd number of nodes, promote the last one
				level = append(level, nodes[i])
			}
		}

		nodes = level
	}

	return &MerkleTree{
		Root:       nodes[0],
		LeafNodes:  leafNodes,
		Depth:      depth,
		TotalNodes: (1 << (depth + 1)) - 1, // 2^(depth+1) - 1
	}, nil
}

// GetRootHash returns the hex-encoded hash of the Merkle Tree root
func (mt *MerkleTree) GetRootHash() string {
	return hex.EncodeToString(mt.Root.Hash)
}

// VerifyData checks if the given data exists in the Merkle Tree
func (mt *MerkleTree) VerifyData(data []byte) bool {
	hash := sha256.Sum256(data)
	targetHash := hash[:]

	// Search through leaf nodes
	for _, leaf := range mt.LeafNodes {
		if hex.EncodeToString(leaf.Hash) == hex.EncodeToString(targetHash) {
			return true
		}
	}

	return false
}

// GenerateMerkleProof creates a Merkle proof for the specified data
func (mt *MerkleTree) GenerateMerkleProof(data []byte) ([][]byte, error) {
	hash := sha256.Sum256(data)
	targetHash := hash[:]

	// Find the corresponding leaf node
	var targetNode *MerkleNode
	for _, leaf := range mt.LeafNodes {
		if hex.EncodeToString(leaf.Hash) == hex.EncodeToString(targetHash) {
			targetNode = leaf
			break
		}
	}

	if targetNode == nil {
		return nil, errors.New("data not found in the Merkle Tree")
	}

	// Generate proof by traversing up to the root
	var proof [][]byte
	current := targetNode

	for current.Parent != nil {
		parent := current.Parent

		// Add sibling to proof
		if parent.Left == current {
			proof = append(proof, parent.Right.Hash)
		} else {
			proof = append(proof, parent.Left.Hash)
		}

		current = parent
	}

	return proof, nil
}

// VerifyMerkleProof verifies if data belongs to the tree using a Merkle proof
func VerifyMerkleProof(dataHash []byte, proof [][]byte, rootHash []byte) bool {
	currentHash := dataHash

	for _, proofElement := range proof {
		// Concatenate and hash
		var combined []byte
		combined = append(combined, currentHash...)
		combined = append(combined, proofElement...)

		hash := sha256.Sum256(combined)
		currentHash = hash[:]
	}

	return hex.EncodeToString(currentHash) == hex.EncodeToString(rootHash)
}
