package miner

import (
	"encoding/hex"
	"testing"
	"time"
)

// Test new miner is created properly.
func TestNewMiner(t *testing.T) {
	// Create a new block
	var pblk *Block
	dif := 1
	data := "Hello World!"
	blk := New(pblk, dif, data)

	// Get the chunk created
	ck := getChunk(blk)

	if ck.Index != 0 {
		t.Errorf("expected index of 0 where no parent block is present, but got %d", ck.Index)
	}

	if ck.Difficulty != dif {
		t.Errorf("expected difficulty to be %d, but got %d", dif, ck.Difficulty)
	}

	if ck.Data != data {
		t.Errorf("expected data to be \"%s\" but got \"%s\"", data, ck.Data)
	}

	// Previous block should be allowed.
	blk2 := New(blk, dif, data)
	if (blk2.Miner).(*Chunk).Parent != (blk.Miner).(*Chunk) {
		t.Errorf("expected parent to match")
	}
}

// Test miner ability to encode its struct to JSON data.
func TestMinerEncode(t *testing.T) {
	blk := createBlock()
	ck := getChunk(blk)

	// Actual and expected.
	a := string(ck.Encode())
	e := "{\"parent_hash\":null,\"hash\":null,\"index\":0,\"pow\":0,\"difficulty\":1,\"data\":\"Hello World!\",\"timestamp\":\"" + ck.Timestamp.Format(time.RFC3339Nano) + "\"}"

	if a != e {
		t.Errorf("expected encode of %s but got %s", a, e)
	}
}

// Test miner ability to encode its stuct to JSON data and create hash of it.
func TestMinerGenerateHash(t *testing.T) {
	blk := createBlock()
	ck := getChunk(blk)

	// Actual and expected.
	a := hex.EncodeToString(ck.GenerateHash(false))
	e := "866da7defdc2df09616dea46e80ffa8c0ea66517a6bd1f70146c86b1b6a54efd"

	if a != e {
		t.Errorf("expected hash of %s but got %s", a, e)
	}
}

// Test if miner mined.
func TestMinerMined(t *testing.T) {
	blk := createBlock()
	ck := getChunk(blk)

	if ck.IsMined() {
		t.Errorf("expected miner to have mined")
	}

	ck.PoW = 1
	if !ck.IsMined() {
		t.Errorf("expected miner to not have mined")
	}
}

// Check miner PoW can validate.
func TestMineValidateNonce(t *testing.T) {
	// With a difficulty of "1".
	// And a parent chunk PoW of "0".
	// It should take "3" tries to solve the problem.
	// Because a SHA256 hash of "0" + "3" a a string,
	// Will equal a hash of "0b8efa5a3bf104413a725c6ff0459a6be12b1fd33314cbb138745baf39504ae5",
	// Which then "0"[:difficulty] == "0".

	blk := createBlock()
	ck := getChunk(blk)

	n := 3                   // PoW of 3
	res := ck.ValidatePoW(n) // result

	if !res {
		t.Errorf("expected PoW of %d to pass but failed", n)
	}

	// Set the PoW to the block and check self validation.
	ck.PoW = n
	if !ck.IsValidPoW() {
		t.Errorf("expected miner PoW to pass but failed")
	}
}

// Test the miner runs the solution to produce a valid PoW and become "mined".
func TestMinerMines(t *testing.T) {
	// Given our solution for validate PoW,
	// A difficulty of "1", should produce an PoW of "3".
	blk := createBlock()
	ck := getChunk(blk)
	ck.Mine()

	if ck.PoW != 3 {
		t.Errorf("expected miner to have mined with a PoW result of 3 but failed")
	}

	if !ck.IsMined() {
		t.Errorf("expected miner to have mined but failed")
	}
}

// Test the entire mined miner can self validate.
func TestMinerValidates(t *testing.T) {
	blk := createBlock()
	ck := getChunk(blk)
	ck.Mine()
	ck.GenerateHash(true)

	if !ck.IsValid() {
		t.Errorf("expected miner to validate itself but failed")
	}
}

// Test a bad miner fails to self validate.
func TestMinerIsNotValid(t *testing.T) {
	blk := createBlock()
	ck := getChunk(blk)
	ck.Mine()
	ck.GenerateHash(true)

	// Fudge the data.
	ck.Timestamp = time.Now()

	if ck.IsValid() {
		t.Errorf("expected miner to be invalid but it returned as valid")
	}
}

// Test the entire miner can self validate with a full parent chunk.
func TestMinerValidatesWithParentChunk(t *testing.T) {
	// Chunk 1
	ck := &Chunk{
		Parent:     nil,
		Index:      0,
		Difficulty: 1,
		Data:       "Hello World",
		Timestamp:  time.Now(),
	}

	ck.Mine()
	ck.GenerateHash(true)

	// Chunk 2
	ck2 := &Chunk{
		Parent:     ck,
		Index:      1,
		Difficulty: 1,
		Data:       "Hello World, Again",
		Timestamp:  time.Now(),
	}

	ck2.Mine()
	ck2.GenerateHash(true)

	if !ck.IsValid() || !ck2.IsValid() {
		t.Errorf("expected miners to validate but failed")
	}
}

// Test a bad parent chunk fails to validate.
func TestMinerIsNotValidWithParentMiner(t *testing.T) {
	// Chunk 1
	pck := new(Chunk)
	ck := &Chunk{
		Parent:     pck,
		Index:      1,
		Difficulty: 1,
		Data:       "Hello World",
		Timestamp:  time.Now(),
	}

	ck.Mine()
	ck.GenerateHash(true)

	// Chunk 2
	ck2 := &Chunk{
		Parent:     ck,
		Index:      1, // Change the index
		Difficulty: 1,
		Data:       "Hellow World, Again",
		Timestamp:  time.Now(),
	}

	ck2.Mine()
	ck2.GenerateHash(true)

	if ck2.IsValid() {
		t.Errorf("expected miner to be invalid but result was valid")
	}
}

// Create a plain miner implementation for the test to use.
func createBlock() (blk *Block) {
	// Create the block with an empty previous.
	blk = New(nil, 1, "Hello World!")

	// Get the chunk created
	ck := getChunk(blk)
	ck.Timestamp = time.Date(2019, 3, 24, 13, 42, 58, 0, time.UTC) // Change the timestamp to something testable.

	return
}

// Quick way to get the chunk created in a block.
func getChunk(blk *Block) *Chunk {
	return (blk.Miner).(*Chunk)
}
