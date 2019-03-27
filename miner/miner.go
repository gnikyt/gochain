package miner

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type (
	// Miner implementation which much be adheard to for Block struct.
	Miner interface {
		Mine() (pow int)
		IsMined() bool
		MarshalJSON() ([]byte, error)
		Encode() (j []byte)
		ValidatePoW(pow int) bool
		IsValidPoW() bool
		GenerateHash(save bool) (sum []byte)
		IsValid() bool
	}

	// Reprecents a block in the chain which contains the miner.
	// The miner will contain a struct like "chunk" which implements the miner interface.
	Block struct {
		Miner
	}

	// Reprecents a chunk and it's data used for mining.
	Chunk struct {
		Parent     *Chunk    `json:"-"`
		Hash       []byte    `json:"hash"`
		Index      int       `json:"index"`
		PoW        int       `json:"pow"`
		Difficulty int       `json:"difficulty"`
		Data       string    `json:"data"`
		Timestamp  time.Time `json:"timestamp"`
	}
)

// Helper to create a new block based on a previous block.
func New(blk *Block, dif int, data string) *Block {
	var pck *Chunk
	if blk == nil {
		// Genesis block.
		pck = new(Chunk)
	} else {
		// Chained.
		pck = (blk.Miner).(*Chunk)
	}

	return &Block{
		Miner: &Chunk{
			Parent:     pck,
			Index:      pck.Index + 1,
			Timestamp:  time.Now(),
			Difficulty: dif,
			Data:       data,
		},
	}
}

// Mines a chunk.
// Will keep running until the PoW is valid and solved for the difficulty.
func (ck *Chunk) Mine() (pow int) {
	for {
		if ck.ValidatePoW(pow) {
			// Solved
			break
		}

		// Not solved, increase.
		pow++
	}

	// Save the PoW to the block.
	ck.PoW = pow

	return
}

// Check if the chunk is mined. Simply checks it has a PoW value.
func (ck Chunk) IsMined() bool {
	return ck.PoW > 0
}

// Marshal for JSON encode.
// Used to add "parent_hash" to the JSON output.
func (ck Chunk) MarshalJSON() ([]byte, error) {
	// Create an alias to the chunck struct to prevent recursion.
	type Alias Chunk

	// Get the parent chunk hash.
	var h []byte
	pck := ck.Parent
	if pck != nil {
		h = pck.Hash
	}

	// Add our parent hash to the alias struct.
	return json.Marshal(
		struct {
			ParentHash []byte `json:"parent_hash"`
			Alias
		}{
			ParentHash: h,
			Alias:      Alias(ck),
		},
	)
}

// Encodes the struct to JSON format.
func (ck Chunk) Encode() (j []byte) {
	j, _ = json.Marshal(ck)

	return
}

// Validates the PoW by combining parent chunk's PoW with input pow.
// Adding both together and hashing, should equal the padding of the difficulty.
func (ck Chunk) ValidatePoW(pow int) bool {
	// Convert the PoW to strings and combine.
	c := strconv.Itoa(ck.Parent.PoW) + strconv.Itoa(pow)

	// Hash the combined PoW, convert the hash to string.
	h := sha256.New()
	h.Write([]byte(c))
	sum := hex.EncodeToString(h.Sum(nil))

	// Pad a "0" string by the difficulty level.
	dif := ck.Difficulty
	pad := fmt.Sprintf("%0*d", dif, 0)

	return sum[:dif] == pad
}

// Checks if this chunk's PoW is valid.
func (ck Chunk) IsValidPoW() bool {
	return ck.ValidatePoW(ck.PoW)
}

// Generate a hash for the chunk based on the chunk's struct data in JSON format.
// Option to save or simply generate.
func (ck *Chunk) GenerateHash(save bool) (sum []byte) {
	// We can't generate a hash of a chunk with a hash.
	// So this hash must be temp removed so we can recaulate.
	osum := ck.Hash
	ck.Hash = nil

	// Encode the JSON of the struct to a hash.
	h := sha256.New()
	h.Write(ck.Encode())
	sum = h.Sum(nil)

	if save {
		// Save the new hash to the block.
		ck.Hash = sum
	} else {
		// Put back the old hash to the block.
		ck.Hash = osum
	}

	return
}

// Confirms the block validity.
func (ck Chunk) IsValid() bool {
	pok, bok := true, true

	// Check if we have a parent chunk to check.
	if pck := ck.Parent; pck != nil && pck.IsMined() {
		// Test parent chunk's index plus one, will equal this chunk's index.
		// Test the hash of parent chunk's hash is what is set for this chunk's parent hash.
		// Test this chunk's PoW is valid.
		if ((pck.Index + 1) != ck.Index) || !bytes.Equal(pck.Hash, ck.Parent.Hash) || !ck.IsValidPoW() {
			pok = false
		}
	}

	// Test this blocks hash is equal to a regeneration of the hash.
	// Test this block is also mined.
	if !bytes.Equal(ck.GenerateHash(false), ck.Hash) || !ck.IsMined() {
		bok = false
	}

	return pok && bok
}
