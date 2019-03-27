package gochain

import (
	"testing"

	"github.com/ohmybrew/gochain/chain"
	"github.com/ohmybrew/gochain/miner"
)

func TestWhole(t *testing.T) {
	// Difficulty level, data for blocks, and new chain.
	dif := 2
	data := "Hello World"
	c := chain.New()

	// Genesis block and another block joined to it.
	blk := miner.New(nil, dif, data)
	blk2 := miner.New(blk, dif, data)

	// Mine both blocks.
	blk.Miner.Mine()
	blk.Miner.GenerateHash(true)
	blk2.Miner.Mine()
	blk2.Miner.GenerateHash(true)

	// Append the blocks.
	c.Append(false, blk)
	c.Append(false, blk2)

	// Validate the chain
	if !c.IsValid() {
		t.Errorf("chain should be valid but returned invalid")
	}
}
