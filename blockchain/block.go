package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"bytes"
	"encoding/json"
)

// Block is one block from blockchain
type Block struct {
	Timestamp time.Time `json:"timestamp"`
	LastHash  BlockHash `json:"lastHash"`
	Hash      BlockHash `json:"hash"`
	Data      BlockData `json:"data"`
}

// BlockHash is hash used in blockchain
type BlockHash []byte

// BlockData is just byte slice for now
type BlockData []byte


// MarshalBinary Represents BlockHash as byte slice
func (hash BlockHash) MarshalBinary() []byte {
	return []byte(hash)
}

// MarshalJSON encodes BlockHash as hex digits
func (hash BlockHash) MarshalJSON() ([]byte, error) {
	hashString := hash.String()
	marshaledString, err := json.Marshal(hashString)
	if err != nil {
		return nil, err
	}
	return marshaledString, nil
}

// Eq checks equality of 2 hashes
func (hash BlockHash) Eq(other BlockHash) bool {
	if bytes.Equal(hash, other) {
		return true
	}
	return false
}

func (hash BlockHash) String() string {
	return hex.EncodeToString(hash)
}


// MarshalBinary represents BlockData as byte slice
func (data BlockData) MarshalBinary() []byte {
	return []byte(data)
}

// MarshalJSON returns BlockData as string
func (data BlockData) MarshalJSON() ([]byte, error) {
	dataString := string(data)
	marshaledString, err := json.Marshal(dataString)
	if err != nil {
		return nil, err
	}
	return marshaledString, nil
}

// UnmarshalJSON parses json string field as BlockData
func (data *BlockData) UnmarshalJSON(jsonData []byte) error {
	var dataString string
	if err := json.Unmarshal(jsonData, &dataString); err != nil {
		return err
	}
	*data = BlockData(dataString)
	return nil
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
