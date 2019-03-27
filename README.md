# gochain

[![Build Status](https://secure.travis-ci.org/ohmybrew/gochain.png?branch=master)](http://travis-ci.org/ohmybrew/gochain)
[![Coverage Status](https://coveralls.io/repos/github/ohmybrew/gochain/badge.svg?branch=master)](https://coveralls.io/github/ohmybrew/gochain?branch=master)

This is a fast and simple implementation which uses a basic SHA256 problem to solve for mining based on a padding of zeros for the difficulty.

## Usage

```go
package main

import (
  "github.com/ohmybrew/gochain/chain"
  "github.com/ohmybrew/gochain/miner"

  "fmt"
)

// See tests for more examples...

// Create a new chain.
c := chain.New()

// Genesis block and another block joined to it.
dif := 3 // Difficulty for miner.
blk := miner.New(nil, dif, "Hello Data")
blk2 := miner.New(blk, dif, "Hi Data")

// Mine both blocks.
blk.Miner.Mine()
blk.Miner.GenerateHash(true)
blk2.Miner.Mine()
blk2.Miner.GenerateHash(true)

// Append the blocks.
c.Append(false, blk)
c.Append(false, blk2)

// See if block is valid.
fmt.Println("Block valid?", c.IsValid())

// Block to JSON
j := blk.Encode()
fmt.Println(string(j)) // {"parent_hash": ..., "hash": ..., "index": ..., "nonce": ..., "timestamp": ..., "difficulty": ..., "data": ...}

// Chain to JSON
cj := c.Encode()
fmt.Println(string(j)) // {"blocks":[{"previous_hash": ..., "hash": ..., "index": ..., "nonce": ..., "timestamp": ..., "difficulty": ..., "data": ...}, {...}]}

// Get first, last, previous blocks
fb, _ := c.First()     // equals blk, if no first block, error will be second return.
lb, _ := c.Last()      // equals blk2, if no last block, error will be second return.
pb, _ := c.Previous(1) // by index, 1 - 1 = 0, so this will equal blk1, if no previous block, error will be second return.
gb, _ := c.Get(1)      // get block by index.
```

## Custom Miner

`miner.New` in above example is a shortcut to create a block (`miner.Block`) with a miner which implements the `miner.Miner` interface.

The built-in miner is `miner.Chunk`.

```go
// Previous chunk, for example purposes is empty.
pck := new(miner.Chunk)

// Create the block with `miner.Chunk` since it satifies the interface.
blk := &miner.Block{
  Miner: &miner.Chunk{
    Parent:     pck,
    Index:      pck.Index + 1,
    Timestamp:  time.Now(),
    Difficulty: dif,
    Data:       data,
  },
}
```

You're free to supply any struct to `miner.Block.Miner` so long as it is compatible with the `miner.Miner` interface. This way, you're able to develop your own mining solutions and validity.

## Testing

`go test`, fully tested.

## Documentation

Available through [godoc.org](https://godoc.org/github.com/ohmybrew/gochain).

Important files:

+ `chain/block.go` contains the struct for a block and its methods.
+ `chain.go` contains the struct for the chain and its methods.
+ `gochain.go` is empty, simply the package index.

## LICENSE

This project is released under the MIT [license](https://github.com/ohmybrew/gochain/blob/master/LICENSE).
