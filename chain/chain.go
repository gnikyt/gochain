package chain

import (
	"errors"

	"encoding/json"

	"github.com/ohmybrew/gochain/miner"
)

// Reprecents a blockchain.
type Chain struct {
	Blocks []*miner.Block `json:"blocks"`
}

// Creates a new chain.
func New() *Chain {
	return new(Chain)
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

// Gets a block by index.
// If no available block is found, error is returned.
func (c Chain) Get(i int) (*miner.Block, error) {
	ct := c.Length()

	if ct == 0 || i > (ct-1) || i < 0 {
		// Out of range.
		return nil, errors.New("no block found")
	}

	return c.Blocks[i], nil
}

// Gets the previous block relative to the provided index.
// If no available previous block is found, error is returned.
func (c Chain) Previous(i int) (*miner.Block, error) {
	return c.Get(i - 1)
}

// Gets the next block relative to the provided index.
// If no available next block is found, error is returned.
func (c Chain) Next(i int) (*miner.Block, error) {
	return c.Get(i + 1)
}

// Get the last block in the chain.
// If no last block is found, error is returned.
func (c Chain) Last() (*miner.Block, error) {
	return c.Get(c.Length() - 1)
}

// Get the first block in the chain.
// If no first block is found, error is returned.
func (c Chain) First() (*miner.Block, error) {
	return c.Get(0)
}

// Appends block to the chain directly.
// Will return error if block is invalid and validation was asked for.
func (c *Chain) Append(ver bool, blk *miner.Block) error {
	// Verify the block if asked to verify by argument one.
	if blk.Miner == nil || (ver && !blk.Miner.IsValid()) {
		return errors.New("can not store block to chain, miner is not valid")
	}

	// All good, append.
	c.Blocks = append(c.Blocks, blk)

	return nil
}

// Walks the chain to ensure all blocks are valid.
func (c Chain) IsValid() bool {
	for _, blk := range c.Blocks {
		if !blk.Miner.IsValid() {
			return false
		}
	}

	return true
}
