package blockchain

import (
	"testing"
	"github.com/ita-sammann/toy-chain"
)

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()
	if bc.chain[0].String() != Genesis().String() {
		t.Error("Blockchain does not start with genesis block")
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock([]byte("TestBlockchain_AddBlock_some_test_data"))
	if string(bc.LastBlock().Data) != "TestBlockchain_AddBlock_some_test_data" {
		t.Error("Data mismatching in last block")
	}
}

func TestBlockchain_IsValid(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_1"))
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_2"))

	if !bc.IsValid() {
		t.Error("Failed to validate valid chain")
	}

	bc.chain[0].Data = []byte("Corrupt data")

	if bc.IsValid() {
		t.Error("Successfully validated chain with corrupt genesis block")
	}

	bc = NewBlockchain()
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_1"))
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_2"))
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_3"))
	bc.chain[2].Data = []byte("gimme all your money")

	if bc.IsValid() {
		t.Error("Successfully validated chain with corrupt intermediate block")
	}

	bc = NewBlockchain()
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_1"))
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_2"))
	bc.AddBlock([]byte("TestBlockchain_IsValid_data_3"))
	bc.chain[bc.Len()-1].Data = []byte("bad data")

	if bc.IsValid() {
		t.Error("Successfully validated chain with corrupt last block")
	}
}

func TestBlockchain_ReplaceChain(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock([]byte("foo 1"))
	bc.AddBlock([]byte("foo 2"))

	bc2 := NewBlockchain()
	bc2.AddBlock([]byte("foo 1"))
	bc2.AddBlock([]byte("foo 2"))
	bc2.AddBlock([]byte("bar 3"))

	err := bc.ReplaceChain(bc2)

	if err != nil || !bc.LastBlock().Hash.Eq(bc2.LastBlock().Hash) {
		t.Error("Failed to replace with valid chain")
	}

	bc3 := NewBlockchain()
	bc3.AddBlock([]byte("baz 1"))
	bc3.AddBlock([]byte("baz 2"))

	err = bc.ReplaceChain(bc3)

	if err != toy_chain.ErrChainReplaceTooShort || bc.LastBlock().Hash.Eq(bc3.LastBlock().Hash) {
		t.Error("Successfully replaced chain with shorter one")
	}

	bc4 := NewBlockchain()
	bc4.AddBlock([]byte("foobar 1"))
	bc4.AddBlock([]byte("foobar 2"))
	bc4.AddBlock([]byte("foobar 3"))
	bc4.AddBlock([]byte("foobar 4"))
	bc4.AddBlock([]byte("foobar 5"))
	bc4.AddBlock([]byte("foobar 6"))
	bc4.AddBlock([]byte("foobar 7"))
	bc4.AddBlock([]byte("foobar 8"))
	bc4.chain[5].Data = []byte("corrupt data")

	err = bc.ReplaceChain(bc4)

	if err != toy_chain.ErrChainReplaceInvalid || bc.LastBlock().Hash.Eq(bc4.LastBlock().Hash) {
		t.Error("Successfully replaced chain with corrupt one")
	}
}
