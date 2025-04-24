package main

import (
	"math"
	"sync"
	"time"
)

// NetworkState represents the current state of the network
type NetworkState struct {
	PartitionProbability float64                  // Probability of network partition
	NodeLatencies        map[string]time.Duration // Latency measurements to each node
	NodeReliability      map[string]float64       // Reliability score for each node
	LastStateUpdate      time.Time
	mu                   sync.RWMutex
}

// ConsistencyLevel represents different consistency guarantees
type ConsistencyLevel int

const (
	StrongConsistency ConsistencyLevel = iota
	CausalConsistency
	EventualConsistency
)

// CAPOptimizer handles the dynamic adjustment of CAP theorem trade-offs
type CAPOptimizer struct {
	networkState       NetworkState
	currentConsistency ConsistencyLevel
	conflictRegistry   map[string]ConflictInfo
	vectorClocks       map[string]VectorClock
	mu                 sync.RWMutex
}

// ConflictInfo stores information about detected conflicts
type ConflictInfo struct {
	StateID      string
	ConflictTime time.Time
	Versions     []StateVersion
	EntropyScore float64
}

// StateVersion represents a specific version of a state item
type StateVersion struct {
	Value      []byte
	UpdateTime time.Time
	NodeID     string
	VClock     VectorClock
}

// VectorClock implements a vector clock for causal ordering
type VectorClock map[string]uint64

// NewCAPOptimizer creates a new CAP optimizer instance
func NewCAPOptimizer() *CAPOptimizer {
	return &CAPOptimizer{
		networkState: NetworkState{
			NodeLatencies:   make(map[string]time.Duration),
			NodeReliability: make(map[string]float64),
			LastStateUpdate: time.Now(),
		},
		currentConsistency: CausalConsistency, // Start with a balanced approach
		conflictRegistry:   make(map[string]ConflictInfo),
		vectorClocks:       make(map[string]VectorClock),
	}
}

// UpdateNetworkTelemetry updates the network state based on recent measurements
func (c *CAPOptimizer) UpdateNetworkTelemetry(nodeID string, latency time.Duration, success bool) {
	c.networkState.mu.Lock()
	defer c.networkState.mu.Unlock()

	// Update node latency
	c.networkState.NodeLatencies[nodeID] = latency

	// Update node reliability
	reliability, exists := c.networkState.NodeReliability[nodeID]
	if !exists {
		reliability = 0.5 // Default starting reliability
	}

	// Adjust reliability based on success/failure
	if success {
		reliability = reliability*0.9 + 0.1 // Slowly increase reliability on success
	} else {
		reliability = reliability * 0.7 // Quickly decrease reliability on failure
	}
	c.networkState.NodeReliability[nodeID] = reliability

	// Recalculate partition probability based on node latencies and reliability
	c.recalculatePartitionProbability()

	c.networkState.LastStateUpdate = time.Now()
}

// recalculatePartitionProbability updates the probability of network partition
func (c *CAPOptimizer) recalculatePartitionProbability() {
	if len(c.networkState.NodeLatencies) == 0 {
		c.networkState.PartitionProbability = 0.5
		return
	}

	// Calculate average latency
	var totalLatency time.Duration
	for _, latency := range c.networkState.NodeLatencies {
		totalLatency += latency
	}
	avgLatency := totalLatency / time.Duration(len(c.networkState.NodeLatencies))

	// Calculate average reliability
	var totalReliability float64
	for _, reliability := range c.networkState.NodeReliability {
		totalReliability += reliability
	}
	avgReliability := totalReliability / float64(len(c.networkState.NodeReliability))

	// Calculate partition probability based on latency and reliability
	// High latency and low reliability increase partition probability
	latencyFactor := math.Min(1.0, float64(avgLatency)/float64(time.Second))
	reliabilityFactor := 1.0 - avgReliability

	c.networkState.PartitionProbability = (latencyFactor*0.7 + reliabilityFactor*0.3)
}

// GetOptimalConsistencyLevel returns the best consistency level based on network conditions
func (c *CAPOptimizer) GetOptimalConsistencyLevel() ConsistencyLevel {
	c.networkState.mu.RLock()
	defer c.networkState.mu.RUnlock()

	// Choose consistency level based on partition probability
	if c.networkState.PartitionProbability < 0.2 {
		return StrongConsistency
	} else if c.networkState.PartitionProbability < 0.6 {
		return CausalConsistency
	} else {
		return EventualConsistency
	}
}

// DynamicTimeout calculates an optimal timeout for operations based on network state
func (c *CAPOptimizer) DynamicTimeout(baseTimeout time.Duration) time.Duration {
	c.networkState.mu.RLock()
	defer c.networkState.mu.RUnlock()

	// Adjust timeout based on network conditions
	adjustmentFactor := 1.0 + c.networkState.PartitionProbability

	// Find the 75th percentile of latencies to account for most nodes
	var latencies []time.Duration
	for _, latency := range c.networkState.NodeLatencies {
		latencies = append(latencies, latency)
	}

	// If we have enough data, use it for better timeout adjustment
	if len(latencies) > 0 {
		// Simple approximation of 75th percentile
		// In a real implementation, you'd want to sort and select properly
		adjustmentFactor *= 1.5 * float64(latencies[0]/time.Millisecond) / 100.0
	}

	// Ensure we don't go crazy with the timeout
	if adjustmentFactor < 1.0 {
		adjustmentFactor = 1.0
	}
	if adjustmentFactor > 10.0 {
		adjustmentFactor = 10.0
	}

	return time.Duration(float64(baseTimeout) * adjustmentFactor)
}

// RegisterStateUpdate updates vector clocks for a state change
func (c *CAPOptimizer) RegisterStateUpdate(stateID string, nodeID string, value []byte) VectorClock {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get or create vector clock for this state
	vClock, exists := c.vectorClocks[stateID]
	if !exists {
		vClock = make(VectorClock)
	}

	// Increment the counter for this node
	vClock[nodeID]++

	// Save the updated vector clock
	c.vectorClocks[stateID] = vClock

	// Return a copy of the vector clock
	result := make(VectorClock)
	for k, v := range vClock {
		result[k] = v
	}

	return result
}

// DetectConflict checks if two state versions are in conflict
func (c *CAPOptimizer) DetectConflict(stateID string, v1, v2 VectorClock) bool {
	// Check if either clock happened strictly before the other
	v1BeforeV2 := true
	v2BeforeV1 := true

	for node, count := range v1 {
		if v2[node] < count {
			v2BeforeV1 = false
		}
	}

	for node, count := range v2 {
		if v1[node] < count {
			v1BeforeV2 = false
		}
	}

	// If neither happened strictly before the other, we have a conflict
	return !v1BeforeV2 && !v2BeforeV1
}

// RegisterConflict records a conflict between multiple state versions
func (c *CAPOptimizer) RegisterConflict(stateID string, versions []StateVersion) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Calculate entropy score based on the divergence between versions
	entropyScore := c.calculateEntropyScore(versions)

	// Register the conflict
	c.conflictRegistry[stateID] = ConflictInfo{
		StateID:      stateID,
		ConflictTime: time.Now(),
		Versions:     versions,
		EntropyScore: entropyScore,
	}
}

// calculateEntropyScore determines how severe a conflict is
func (c *CAPOptimizer) calculateEntropyScore(versions []StateVersion) float64 {
	if len(versions) <= 1 {
		return 0.0
	}

	// A simple entropy calculation based on differences between versions
	// In a real implementation, this would use more sophisticated metrics
	return float64(len(versions)) / 10.0
}

// ResolveConflict implements probabilistic conflict resolution
func (c *CAPOptimizer) ResolveConflict(stateID string) StateVersion {
	c.mu.Lock()
	defer c.mu.Unlock()

	conflict, exists := c.conflictRegistry[stateID]
	if !exists || len(conflict.Versions) == 0 {
		return StateVersion{} // No conflict to resolve
	}

	// For serious conflicts (high entropy), use more complex resolution
	if conflict.EntropyScore > 0.5 {
		return c.resolveHighEntropyConflict(conflict)
	}

	// For simpler conflicts, use time-based resolution
	return c.resolveTimeBasedConflict(conflict)
}

// resolveHighEntropyConflict handles serious conflicts
func (c *CAPOptimizer) resolveHighEntropyConflict(conflict ConflictInfo) StateVersion {
	// Implement a weighted voting mechanism where more reliable nodes have more weight
	versionScores := make(map[int]float64)

	for i, version := range conflict.Versions {
		// Start with node reliability as the base score
		reliability := c.networkState.NodeReliability[version.NodeID]
		if reliability == 0 {
			reliability = 0.5 // Default if no data
		}

		// Favor more recent updates slightly
		timeScore := 1.0 - float64(time.Since(version.UpdateTime))/float64(time.Hour)
		if timeScore < 0 {
			timeScore = 0
		}

		// Calculate final score
		versionScores[i] = reliability*0.7 + timeScore*0.3
	}

	// Find version with highest score
	bestScore := -1.0
	bestIndex := 0

	for i, score := range versionScores {
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	// Remove the conflict from registry after resolution
	delete(c.conflictRegistry, conflict.StateID)

	return conflict.Versions[bestIndex]
}

// resolveTimeBasedConflict resolves simpler conflicts based primarily on timestamp
func (c *CAPOptimizer) resolveTimeBasedConflict(conflict ConflictInfo) StateVersion {
	var newestVersion StateVersion
	var newestTime time.Time

	for _, version := range conflict.Versions {
		if newestTime.IsZero() || version.UpdateTime.After(newestTime) {
			newestTime = version.UpdateTime
			newestVersion = version
		}
	}

	// Remove the conflict from registry after resolution
	delete(c.conflictRegistry, conflict.StateID)

	return newestVersion
}

// GetConsistencyStatus returns a human-readable status of the current CAP configuration
func (c *CAPOptimizer) GetConsistencyStatus() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	level := c.GetOptimalConsistencyLevel()

	switch level {
	case StrongConsistency:
		return "Strong Consistency (CP mode)"
	case CausalConsistency:
		return "Causal Consistency (balanced mode)"
	case EventualConsistency:
		return "Eventual Consistency (AP mode)"
	default:
		return "Unknown Consistency Mode"
	}
}
