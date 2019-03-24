package gochain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Reprecents a block.
type Block struct {
	Previous     *Block    `json:"-"`
	PreviousHash []byte    `json:"previous_hash"`
	Hash         []byte    `json:"hash"`
	Index        int       `json:"index"`
	Nonce        int       `json:"nonce"`
	Difficulty   int       `json:"difficulty"`
	Data         string    `json:"data"`
	Timestamp    time.Time `json:"timestamp"`
}

// Encodes the struct to JSON format.
func (blk Block) Encode() (j []byte) {
	j, _ = json.Marshal(blk)

	return
}

// Validates the nonce by combining previous block's nonce with input n.
// Adding both together and hashing, should equal the padding of the difficulty.
func (blk Block) ValidateNonce(n int) bool {
	// Convert the nonce to strings and combine.
	cn := strconv.Itoa(blk.Previous.Nonce) + strconv.Itoa(n)

	// Hash the combined nonce, convert the hash to string.
	h := sha256.New()
	h.Write([]byte(cn))
	hs := hex.EncodeToString(h.Sum(nil))

	// Pad a "0" string by the difficulty level.
	dif := blk.Difficulty
	pad := fmt.Sprintf("%0*d", dif, 0)

	return hs[:dif] == pad
}

// Checks if this block's nonce is valid.
func (blk Block) IsValidNonce() bool {
	return blk.ValidateNonce(blk.Nonce)
}

// Mines the block.
// Will keep running until the nonce is valid and solved for the difficulty.
func (blk *Block) Mine() (n int) {
	for {
		if ok := blk.ValidateNonce(n); ok {
			// Solved
			break
		}

		// Not solved, increase.
		n = n + 1
	}

	// Save the nonce to the block.
	blk.Nonce = n

	return
}

// Check if the block is minded. Simply checks it has an nonce value.
func (blk Block) IsMined() bool {
	return blk.Nonce > 0
}

// Generate a hash for the block based on the block's struct data in JSON format.
// Option to save or simply generate.
func (blk *Block) GenerateHash(save bool) (hs []byte) {
	// We can't generate a hash of a block with a hash.
	// So this hash must be temp removed so we can recaulate.
	oh := blk.Hash
	blk.Hash = nil

	// Encode the JSON of the struct to a hash.
	h := sha256.New()
	h.Write(blk.Encode())
	hs = h.Sum(nil)

	if save {
		// Save the new hash to the block.
		blk.Hash = hs
	} else {
		// Put back the old hash to the block.
		blk.Hash = oh
	}

	return
}

// Confirms the block is valid.
func (blk Block) IsValid() bool {
	pok := true
	bok := true

	// Check if we have a previous block to use
	if pb := blk.Previous; pb.IsMined() {
		// Test previous block's index plus one, will equal this block's index.
		// Test the hash of previous block's hash is what is set for this block's previous hash.
		// Test this block's nonce is valid
		if ((pb.Index + 1) == blk.Index) && bytes.Equal(pb.Hash, blk.PreviousHash) && blk.IsValidNonce() {
			pok = true
		} else {
			pok = false
		}
	}

	// Test this blocks hash is equal to a regeneration of the hash.
	// Test this block is also mined.
	if bytes.Equal(blk.GenerateHash(false), blk.Hash) && blk.IsMined() {
		bok = true
	} else {
		bok = false
	}

	return pok && bok
}
