package gochain

import (
	"github.com/google/go-cmp/cmp"
	"time"
)

// Reprecents a block chain.
type Chain struct {
	Blocks []*Block
}

// Return the chain length.
func (c Chain) Length() int {
	return len(c.Blocks)
}

// Gets the previous block relative to the provided index.
// If no available previous block is found, nil is returned
func (c Chain) PreviousBlock(i int) *Block {
	// Next index to target
	ni := i - 1

	if c.Length() > 0 && ni >= 0 {
		// Previous block index is available, get it.
		return c.Blocks[ni]
	}

	// No previous block, return empty.
	return new(Block)
}

// Get the last block in the chain.
func (c Chain) LastBlock() *Block {
	// Count the length of the chain.
	ct := c.Length()

	if ct > 0 {
		// We can get the last block.
		return c.Blocks[ct-1]
	}

	// No last block, return empty.
	return new(Block)
}

// Get the first block in the chain.
func (c Chain) FirstBlock() *Block {
	// Count the length of the chain.
	ct := c.Length()

	if ct > 0 {
		// We can get the last block.
		return c.Blocks[0]
	}

	// No first block, return empty.
	return new(Block)
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
	lb := c.LastBlock()

	// Build the new block based on the previous block.
	blk = &Block{
		Previous:     lb,
		PreviousHash: lb.Hash,
		Index:        lb.Index + 1,
		Timestamp:    time.Now(),
		Difficulty:   dif,
		Data:         dat,
	}

	// Add to the chain
	c.Blocks = append(c.Blocks, blk)

	return
}

// Walks the chain to ensure all blocks are valid.
func (c Chain) IsValid() bool {
	for _, blk := range c.Blocks {
		if ok := blk.IsValid(); ok != true {
			return false
		}
	}

	return true
}

func (c Chain) IsSameBlock(b1 *Block, b2 *Block) bool {
	return cmp.Equal(b1.Hash, b2.Hash)
}
