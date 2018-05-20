package blockchain

import (
	"testing"
	"time"
)

func TestBlock_String(t *testing.T) {
	block := Block{
		time.Date(2018, 5, 20, 15, 30, 45, 123456789, time.UTC),
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		[]byte("this is some data in block"),
	}
	expectedString := `Block -
	Timestamp : 2018-05-20 15:30:45.123456
	LastHash  : 6162636465666768696a6b6c6d6e6f707172737475767778797a
	Hash      : 4142434445464748494a4b4c4d4e4f505152535455565758595a
	Data      : this is some data in block
`
	if block.String() != expectedString {
		t.Errorf("Block.String() is incorrect.\nGot:\n%s\n\nExpected:\n%s", block.String(), expectedString)
	}
}

func TestGenesis(t *testing.T) {
	timestamp := time.Unix(1526774400, 0)
	lastHash := []byte{
		0, 0, 0, 0, 255, 255, 255, 255,
		0, 0, 0, 0, 255, 255, 255, 255,
		0, 0, 0, 0, 255, 255, 255, 255,
		0, 0, 0, 0, 255, 255, 255, 255,
	}
	data := []byte("genesis")

	expectedBlock := Block{
		timestamp,
		lastHash,
		Hash(timestamp, lastHash, data),
		data,
	}
	if Genesis().String() != expectedBlock.String() {
		t.Errorf("Bad genesis block.\nGot:\n%s\n\nExpected:\n%s", Genesis(), expectedBlock)
	}
}

func TestMineBlock(t *testing.T) {
	newBlock := MineBlock(Genesis(), []byte("test_data"))

	if newBlock.LastHash.String() != Genesis().Hash.String() {
		t.Errorf("Bad last hash: %s", newBlock.LastHash)
	}
	if newBlock.Hash.String() != Hash(newBlock.Timestamp, newBlock.LastHash, newBlock.Data).String() {
		t.Errorf("Bad hash: %s", newBlock.Hash)
	}
	if string(newBlock.Data) != "test_data" {
		t.Errorf("Bad data: %s", newBlock.Data)
	}
}
