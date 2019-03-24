package gochain

import (
	"encoding/hex"
	"testing"
	"time"
)

// Test block ability to encode its struct to JSON data.
func TestBlockEncode(t *testing.T) {
	blk := createBlock()

	// Actual and expected.
	a := string(blk.Encode())
	e := "{\"previous_hash\":null,\"hash\":null,\"index\":1,\"nonce\":0,\"difficulty\":1,\"data\":\"Hello World\",\"timestamp\":\"" + blk.Timestamp.Format(time.RFC3339Nano) + "\"}"

	if a != e {
		t.Errorf("Expected encode of %s but got %s", a, e)
	}
}

// Test block ability to encode its stuct to JSON data and create hash of it.
func TestBlockGenerateHash(t *testing.T) {
	blk := createBlock()

	// Actual and expected.
	a := hex.EncodeToString(blk.GenerateHash(false))
	e := "5a4a6d8ab0f989674913ba81208a324779717e1d82167759b9df36416016978b"

	if a != e {
		t.Errorf("Expected hash of %s but got %s", a, e)
	}
}

// Test if block was mined.
func TestBlockMined(t *testing.T) {
	blk := createBlock()

	if blk.IsMined() {
		t.Errorf("Expected block to be unmined")
	}

	blk.Nonce = 1
	if !blk.IsMined() {
		t.Errorf("Expected block to be mined")
	}
}

// Check block nonce can validate.
func TestBlockValidateNonce(t *testing.T) {
	// With a difficulty of "1".
	// And a previous block nonce of "0".
	// It should take "3" tries to solve the problem.
	// Because a SHA256 hash of "0" + "3" a a string,
	// Will equal a hash of "0b8efa5a3bf104413a725c6ff0459a6be12b1fd33314cbb138745baf39504ae5",
	// Which then "0"[:difficulty] == "0".

	blk := createBlock()

	n := 3                      // nonce of 3
	res := blk.ValidateNonce(n) // result

	if !res {
		t.Errorf("Expected nonce of %d to pass but failed", n)
	}

	// Set the nonce to the block and check self validation.
	blk.Nonce = n
	if !blk.IsValidNonce() {
		t.Errorf("Expected block nonce to pass but failed")
	}
}

// Test the block runs the solution to produce a valid nonce and become "mined".
func TestBlockMines(t *testing.T) {
	// Given our solution for validate nonce,
	// A difficulty of "1", should produce an nonce of "3".
	blk := createBlock()
	blk.Mine()

	if blk.Nonce != 3 {
		t.Errorf("Expected block to mine with an nonce result of 3 but failed")
	}

	if !blk.IsMined() {
		t.Errorf("Expected block to be mined but failed")
	}
}

// Test the entire mined block can self validate.
func TestBlockValidates(t *testing.T) {
	blk := createBlock()
	blk.Mine()
	blk.GenerateHash(true)

	if !blk.IsValid() {
		t.Errorf("Expected block to validate itself but failed")
	}
}

// Test a bad block fails to self validate.
func TestBlockIsNotValid(t *testing.T) {
	blk := createBlock()
	blk.Mine()
	blk.GenerateHash(true)

	// Fudge the data
	blk.Timestamp = time.Now()

	if blk.IsValid() {
		t.Errorf("Expected block to be invalid but it returned as valid")
	}
}

// Create a plain unmined block for test use.
func createBlock() (blk *Block) {
	// Create an empty previous block.
	pb := new(Block)

	// Set a static time we can test against.
	ts := time.Date(2019, 3, 24, 13, 42, 58, 0, time.UTC)

	// Create the block.
	blk = &Block{
		Previous:     pb,
		PreviousHash: pb.Hash,
		Index:        1,
		Difficulty:   1,
		Data:         "Hello World",
		Timestamp:    ts,
	}

	return
}
