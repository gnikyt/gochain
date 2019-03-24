# GoChain

[![Build Status](https://secure.travis-ci.org/ohmybrew/gochain.png?branch=master)](http://travis-ci.org/ohmybrew/gochain)
[![Coverage Status](https://coveralls.io/repos/github/ohmybrew/gochain/badge.svg?branch=master)](https://coveralls.io/github/ohmybrew/gochain?branch=master)

Port of my Blockchain-PHP library to a Golang. The speed is roughly 90% faster.

## Usage

```go
package main

import (
  gc "github.com/ohmybrew/gochain"
  "fmt"
)

// See tests for more examples...

// New chain.
c := new(gc.Chain)

// Add two blocks, mine them. Difficulty of "2".
dif := 2
blk1 := c.BuildBlock(dif, "One")
blk2 := c.BuildBlock(dif, "Two")

blk1.Mine()
blk1.GenerateHash(true)
blk2.Mine()
blk2.GenerateHash(true)

fmt.Println("Block valid?", blk.IsValid())
fmt.Println("Block valid?", blk2.IsValid())
fmt.Println("Same block?", c.IsSameBlock(blk, blk))
fmt.Println("Chain is valid?", c.IsValid())

// Block to JSON
j := blk1.Encode()
fmt.Println(string(j)) // example: {"previous_hash": ..., "hash": ..., "index": ..., "nonce": ..., "timestamp": ..., "difficulty": ..., "data": ...}

// Chain to JSON
cj := c.Encode()
fmt.Println(string(j)) // example: {"blocks":[{"previous_hash": ..., "hash": ..., "index": ..., "nonce": ..., "timestamp": ..., "difficulty": ..., "data": ...}, {...}]}

// Get first, last, previous blocks
fb, _ := c.FirstBlock() // equals blk1, if no first block, error will be second return
lb, _ := c.LastBlock() // equals blk2, if no last block, error will be second return
pb, _ := c.PreviousBlock(1) // by index, 1 - 1 = 0, so this will equal blk1, if no previous block, error will be second return
```

## Testing

`bin/test` for test suite.

`bin/cover` for test suite with coverage output.

## Documentation

Available through [godoc.org](https://godoc.org/github.com/ohmybrew/gochain).

Important files:

+ `block.go` contains the struct for a block and its methods.
+ `chain.go` contains the struct for the chain and its methods.
+ `gochain.go` is empty, simply the package index.

## LICENSE

This project is released under the MIT [license](https://github.com/ohmybrew/gochain/blob/master/LICENSE).
