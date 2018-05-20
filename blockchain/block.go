package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"bytes"
)

// BlockHash is hash used in blockchain
type BlockHash []byte

// MarshalBinary Represents BlockHash as byte slice
func (hash BlockHash) MarshalBinary() []byte {
	return []byte(hash)
}

func (hash BlockHash) String() string {
	return hex.EncodeToString(hash)
}

// BlockData is just byte slice for now
type BlockData []byte

// MarshalBinary represents BlockData as byte slice
func (data BlockData) MarshalBinary() []byte {
	return []byte(data)
}

// Block is one block from blockchain
type Block struct {
	Timestamp time.Time
	LastHash  BlockHash
	Hash      BlockHash
	Data      BlockData
}

func (block Block) String() string {
	return fmt.Sprintf(
		"Block -\n\tTimestamp : %s\n\tLastHash  : %s\n\tHash      : %s\n\tData      : %s\n",
		block.Timestamp.Format("2006-01-02 15:04:05.999999"),
		block.LastHash,
		block.Hash,
		block.Data,
	)
}

func (block Block) checkHash() bool {
	return bytes.Equal(block.Hash, Hash(block.Timestamp, block.LastHash, block.Data))
}

// Genesis return genesis block
func Genesis() Block {
	timestamp := time.Unix(1526774400, 0)
	lastHash := []byte{
		0, 0, 0, 0, 255, 255, 255, 255,
		0, 0, 0, 0, 255, 255, 255, 255,
		0, 0, 0, 0, 255, 255, 255, 255,
		0, 0, 0, 0, 255, 255, 255, 255,
	}
	data := []byte("genesis")

	return Block{
		timestamp,
		lastHash,
		Hash(timestamp, lastHash, data),
		data,
	}
}

// MineBlock generates new block based on last block of chain
func MineBlock(lastBlock Block, data BlockData) Block {
	timestamp := time.Now()
	lastHash := lastBlock.Hash
	hash := Hash(timestamp, lastBlock.Hash, data)

	return Block{
		timestamp,
		lastHash,
		hash,
		data,
	}
}

// Hash generates hash for new block
func Hash(timestamp time.Time, lastHash BlockHash, data BlockData) BlockHash {
	binTimestamp, err := timestamp.MarshalBinary()
	if err != nil {
		panic(err)
	}
	dataToHash := make([]byte, 64)
	dataToHash = append(dataToHash, binTimestamp...)
	dataToHash = append(dataToHash, lastHash.MarshalBinary()...)
	dataToHash = append(dataToHash, data.MarshalBinary()...)

	hash := sha256.Sum256(dataToHash)
	return hash[:]
}
