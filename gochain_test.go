package gochain

import (
	"fmt"
	"testing"
)

// Testing not yet started
func TestStuff(t *testing.T) {
	c := new(Chain)

	a := c.BuildBlock(5, "Hello")
	a.Mine()
	a.GenerateHash()

	b := c.BuildBlock(5, "Hello 2")
	b.Mine()
	b.GenerateHash()

	fmt.Println("Chain valid?", c.IsValid())
}
