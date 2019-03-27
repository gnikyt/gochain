package chain

import (
	"testing"
	"time"

	"github.com/ohmybrew/gochain/miner"
)

// Test chain has a length.
func TestChainLength(t *testing.T) {
	c := New()
	len := c.Length()

	if len != 0 {
		t.Errorf("expected len to be 0 but got %d", len)
	}
}

// Test getting previous block.
func TestPreviousBlock(t *testing.T) {
	// New chian.
	c := createFakeChain()
	if _, err := c.Previous(1); err != nil {
		t.Errorf("expected previous block behind index 1 to equal the first block's data")
	}

	if _, err := c.Previous(0); err == nil {
		t.Errorf("expected previous block behind index 0 to return nil")
	}
}

// Test getting next block.
func TestNextBlock(t *testing.T) {
	// New chian.
	c := createFakeChain()
	if _, err := c.Next(0); err != nil {
		t.Errorf("expected next block ahead of index 0 to equal the second block's data")
	}

	if _, err := c.Next(1); err == nil {
		t.Errorf("expected next block ahead of index 1 to return nil")
	}
}

// Test can get last block in chain.
func TestLastBlock(t *testing.T) {
	// New chain.
	c := createFakeChain()
	if _, err := c.Last(); err != nil {
		t.Errorf("expected last block to equal last built block")
	}

	// New chain. Test no blocks in chain, last block should not exist.
	c = new(Chain)
	if _, err := c.Last(); err == nil {
		t.Errorf("expected last block to be nil")
	}
}

// Test can get first block in chain.
func TestFirstBlock(t *testing.T) {
	// New chain.
	c := createFakeChain()
	if _, err := c.First(); err != nil {
		t.Errorf("expected first block to equal first built block")
	}

	// New chain. Test no blocks in chain, first block should not exist.
	c = new(Chain)
	if _, err := c.First(); err == nil {
		t.Errorf("expected last block to be nil")
	}
}

// Test append to chain.
func TestAppendToChain(t *testing.T) {
	// New chain.
	c := New()

	// Append block to chain.
	blk := miner.New(nil, 1, "One")
	c.Append(false, blk)

	// Confirm its added.
	l := c.Length()
	if l != 1 {
		t.Errorf("expected chain length to be 1, got %d", l)
	}
}

// Test append to chain with invalid.
func TestAppendToChainWithInvalid(t *testing.T) {
	// No parent.
	c := New()
	if err := c.Append(true, new(miner.Block)); err == nil {
		t.Errorf("expected append to be invalid but resulted in success")
	}

	// Invalid block.
	c = New()
	blk := miner.New(nil, 1, "One")
	blk2 := miner.New(blk, 1, "Two")
	(blk.Miner).(*miner.Chunk).Data = "Oops" // Create the invalid issue.

	if err := c.Append(true, blk, blk2); err == nil {
		t.Errorf("expected append to be invalid but resulted in success")
	}
}

// Test chain validates.
func TestValidChain(t *testing.T) {
	// New chain.
	c := createFakeChain()
	blks := c.Blocks

	blks[0].Mine()
	blks[0].GenerateHash(true)
	blks[1].Mine()
	blks[1].GenerateHash(true)

	// Check the chain.
	if !c.IsValid() {
		t.Errorf("expected chain to validate but failed")
	}
}

// Test chain invalidates.
func TestInvalidChain(t *testing.T) {
	// New chain.
	c := createFakeChain()
	blk1, _ := c.Get(0)
	blk2, _ := c.Get(1)

	blk1.Mine()
	blk1.GenerateHash(true)
	blk1.Mine()
	blk1.GenerateHash(true)

	// Change some data.
	ck := (blk2.Miner).(*miner.Chunk)
	ck.Index = 3
	ck.Hash = nil

	// Check the chain.
	if c.IsValid() {
		t.Errorf("expected chain to invalid but result was valid")
	}
}

// Test chain ability to encode its struct to JSON data.
func TestChainEncode(t *testing.T) {
	// New block.
	blk := &miner.Block{
		Miner: &miner.Chunk{
			Parent:     nil,
			Timestamp:  time.Now(),
			Index:      1,
			Difficulty: 1,
			Data:       "Hello World",
		},
	}
	ck := (blk.Miner).(*miner.Chunk)

	// New chain.
	c := New()
	c.Append(false, blk)

	// Actual and expected.
	a := string(c.Encode())
	e := "{\"blocks\":[{\"parent_hash\":null,\"hash\":null,\"index\":1,\"pow\":0,\"difficulty\":1,\"data\":\"Hello World\",\"timestamp\":\"" + ck.Timestamp.Format(time.RFC3339Nano) + "\"}]}"

	if a != e {
		t.Errorf("expected encode of %s but got %s", a, e)
	}
}

// Create a fake chain for testing
func createFakeChain() (c *Chain) {
	c = New()

	blk := miner.New(nil, 1, "One")
	blk2 := miner.New(blk, 1, "Two")

	c.Append(false, blk, blk2)

	return
}
