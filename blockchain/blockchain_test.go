package blockchain

import "testing"

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()
	if bc.chain[0].String() != Genesis().String() {
		t.Error("Blockchain does not start with genesis block")
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock([]byte("TestBlockchain_AddBlock_some_test_data"))
	if string(bc.chain[len(bc.chain)-1].Data) != "TestBlockchain_AddBlock_some_test_data" {
		t.Error("Data mismatching in last block")
	}
}
