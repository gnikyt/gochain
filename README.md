# GoChain

[![Build Status](https://secure.travis-ci.org/ohmybrew/gochain.png?branch=master)](http://travis-ci.org/ohmybrew/gochain)
[![Coverage Status](https://coveralls.io/repos/github/ohmybrew/gochain/badge.svg?branch=master)](https://coveralls.io/github/ohmybrew/gochain?branch=master)

Port of my Blockchain-PHP library to a Golang. The speed is roughly 90% faster.

*Note: Not yet completed.*

## Usage

```go
package main

import (
  gc "github.com/ohmybrew/gochain"
  "fmt"
)

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

fmt.Println("Chain is valid?", c.IsValid()); // See tests for more examples
```

## Testing

`bin/test` for test suite.

`bin/cover` for test suite with coverage output.

## LICENSE

This project is released under the MIT [license](https://github.com/ohmybrew/gochain/blob/master/LICENSE).
