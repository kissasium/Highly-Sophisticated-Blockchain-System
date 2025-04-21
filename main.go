package main

import "fmt"

func main() {
	var Blockchain []Block

	// Create genesis block
	genesis := generateGenesisBlock()
	Blockchain = append(Blockchain, genesis)

	// Add a couple more blocks
	block1 := generateNextBlock(Blockchain[len(Blockchain)-1], "First real block")
	Blockchain = append(Blockchain, block1)

	block2 := generateNextBlock(Blockchain[len(Blockchain)-1], "Second real block")
	Blockchain = append(Blockchain, block2)

	// Print the blockchain
	for _, block := range Blockchain {
		fmt.Printf("\nIndex: %d\nTimestamp: %s\nData: %s\nPrevHash: %s\nHash: %s\n",
			block.Index, block.Timestamp, block.Data, block.PrevHash, block.Hash)
	}
}
