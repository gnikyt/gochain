package gochain

import (
	"bytes"
	"errors"
	"time"

	"encoding/json"
)

// Reprecents a block chain.
type Chain struct {
	Blocks []*Block `json:"blocks"`
}

// Encodes the struct to JSON format.
func (c Chain) Encode() (j []byte) {
	j, _ = json.Marshal(c)

	return
}

// Return the chain length.
func (c Chain) Length() int {
	return len(c.Blocks)
}

// Gets the previous block relative to the provided index.
// If no available previous block is found, nil is returned
func (c Chain) PreviousBlock(i int) (*Block, error) {
	// Next index to target
	ni := i - 1

	if c.Length() > 0 && ni >= 0 {
		// Previous block index is available, get it.
		return c.Blocks[ni], nil
	}

	// No previous block.
	return nil, errors.New("no previous block")
}

// Get the last block in the chain.
func (c Chain) LastBlock() (*Block, error) {
	// Count the length of the chain.
	ct := c.Length()

	if ct > 0 {
		// We can get the last block.
		return c.Blocks[ct-1], nil
	}

	// No last block.
	return nil, errors.New("no last block")
}

// Get the first block in the chain.
func (c Chain) FirstBlock() (*Block, error) {
	// Count the length of the chain.
	ct := c.Length()

	if ct > 0 {
		// We can get the last block.
		return c.Blocks[0], nil
	}

	// No first block.
	return nil, errors.New("no first block")
}

// Adds a block to the chain and returns the chain.
func (c *Chain) AddBlock(blk *Block) *Chain {
	c.Blocks = append(c.Blocks, blk)

	return c
}

// Builds a block easily.
// Difficult and data are the only two items required.
// The newly created block is returned and added to the chain.
func (c *Chain) BuildBlock(dif int, dat string) (blk *Block) {
	// Get the last block
	lb, err := c.LastBlock()
	if err != nil {
		// Make a empty block to use as previous.
		lb = new(Block)
	}

	// Build the new block based on the previous block.
	blk = &Block{
		Previous:   lb,
		Index:      lb.Index + 1,
		Timestamp:  time.Now(),
		Difficulty: dif,
		Data:       dat,
	}

	// Add to the chain
	c.Blocks = append(c.Blocks, blk)

	return
}

// Checks if the blocks are the same (simple hash check).
func (c Chain) IsSameBlock(b1 *Block, b2 *Block) bool {
	return bytes.Equal(b1.Hash, b2.Hash)
}

// Walks the chain to ensure all blocks are valid.
func (c Chain) IsValid() bool {
	for _, blk := range c.Blocks {
		if ok := blk.IsValid(); !ok {
			return false
		}
	}

	return true
}
