package gochain

import (
	"reflect"
	"testing"
)

// Test chain has a length.
func TestChainLength(t *testing.T) {
	c := new(Chain)
	len := c.Length()

	if len != 0 {
		t.Errorf("Expected len to be 0 but got %d", len)
	}
}

// Test the block creation shortcut works.
func TestBuildBlock(t *testing.T) {
	// New chain.
	c := new(Chain)

	// Create an expected block to test against.
	dat := "Hello World"
	dif := 3
	epb := new(Block)
	eblk := &Block{
		Previous:   epb,
		Index:      1,
		Difficulty: dif,
		Data:       dat,
	}

	// Create a block
	blk := c.BuildBlock(dif, dat)

	a := reflect.TypeOf(blk).Elem().Name() // Actual
	e := "Block"                           // Expected

	// Test type.
	if a != e {
		t.Errorf("Expected type %s but got %s", e, a)
	}

	// Test index.
	if blk.Index != eblk.Index {
		t.Errorf("Expected first block index to be %d, got %d", eblk.Index, blk.Index)
	}

	// Test data.
	if blk.Data != eblk.Data {
		t.Errorf("Expected block data to be \"%s\", got \"%s\"", eblk.Data, blk.Data)
	}

	// Test difficulty.
	if blk.Difficulty != eblk.Difficulty {
		t.Errorf("Expected difficulty of %d, got %d", eblk.Difficulty, blk.Difficulty)
	}

	// Test chain length went up by 1.
	l := c.Length()
	if l != 1 {
		t.Errorf("Expected chain to be 1 in length, got %d", l)
	}
}

// Test getting previous block.
func TestPreviousBlock(t *testing.T) {
	// New chian.
	c := new(Chain)
	c.BuildBlock(1, "One")
	c.BuildBlock(2, "Two")

	if blk, _ := c.PreviousBlock(1); blk.Data != "One" {
		t.Errorf("Expected previous block behind index 1 to equal the first block's data")
	}

	// New chain. Test returns a an error for no previous block.
	c = new(Chain)
	c.BuildBlock(1, "One")

	if _, err := c.PreviousBlock(0); err == nil {
		t.Errorf("Expected previous block behind index 0 to return nil")
	}
}

// Test can get last block in chain.
func TestLastBlock(t *testing.T) {
	// New chain.
	c := new(Chain)
	c.BuildBlock(1, "One")
	c.BuildBlock(2, "Two")

	if blk, _ := c.LastBlock(); blk.Data != "Two" {
		t.Errorf("Expected last block to equal last built block")
	}

	// New chain. Test no blocks in chain, last block should not exist.
	c = new(Chain)

	if _, err := c.LastBlock(); err == nil {
		t.Errorf("Expected last block to be nil")
	}
}

// Test can get first block in chain.
func TestFirstBlock(t *testing.T) {
	// New chain.
	c := new(Chain)
	c.BuildBlock(1, "One")
	c.BuildBlock(2, "Two")

	if blk, _ := c.FirstBlock(); blk.Data != "One" {
		t.Errorf("Expected first block to equal first built block")
	}

	// New chain. Test no blocks in chain, first block should not exist.
	c = new(Chain)

	if _, err := c.FirstBlock(); err == nil {
		t.Errorf("Expected last block to be nil")
	}
}

// Test add block.
func TestAddBlockToChain(t *testing.T) {
	// New chain.
	c := new(Chain)

	// Add block to chain.
	blk := new(Block)
	c.AddBlock(blk)

	// Confirm its added.
	l := c.Length()
	if l != 1 {
		t.Errorf("Expected chain length to be 1, got %d", l)
	}
}

// Test blocks are same (just by hash)
func TestIsSameBlock(t *testing.T) {
	// New chain.
	c := new(Chain)
	blk := c.BuildBlock(1, "Hello World")

	// Get the first block (which should be blk)
	fb, _ := c.FirstBlock()

	if !c.IsSameBlock(blk, fb) {
		t.Errorf("Expected both block hashes to match")
	}
}

// Test chain validates.
func TestValidChain(t *testing.T) {
	// New chain.
	c := new(Chain)

	// Add two blocks, mine them.
	blk1 := c.BuildBlock(2, "One")
	blk2 := c.BuildBlock(2, "Two")

	blk1.Mine()
	blk1.GenerateHash(true)
	blk2.Mine()
	blk2.GenerateHash(true)

	// Check the chain.
	if !c.IsValid() {
		t.Errorf("Expected chain to validate but failed")
	}
}

// Test chain invalidates.
func TestInValidChain(t *testing.T) {
	// New chain.
	c := new(Chain)

	// Add two blocks, mine them.
	blk1 := c.BuildBlock(2, "One")
	blk2 := c.BuildBlock(2, "Two")

	blk1.Mine()
	blk1.GenerateHash(true)
	blk2.Mine()
	blk2.GenerateHash(true)

	// Change some data.
	blk2.Index = 3
	blk2.Hash = nil

	// Check the chain.
	if c.IsValid() {
		t.Errorf("Expected chain to invalid but result was valid")
	}
}
