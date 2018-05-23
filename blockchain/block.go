package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	DefaultDifficulty = 2
	MineRate          = 3 * time.Second
	MineRateError     = 500 * time.Millisecond
)

// Block is one block from blockchain
type Block struct {
	Timestamp  time.Time `json:"timestamp"`
	LastHash   BlockHash `json:"lastHash"`
	Hash       BlockHash `json:"hash"`
	Data       BlockData `json:"data"`
	Nonce      uint64    `json:"nonce"`
	Difficulty uint8     `json:"difficulty"`
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

// UnmarshalJSON decodes BlockHash from hex digits
func (hash *BlockHash) UnmarshalJSON(jsonHash []byte) error {
	var hashString string
	if err := json.Unmarshal(jsonHash, &hashString); err != nil {
		return err
	}

	hashBytes, err := NewBlockHashFromString(hashString)
	if err != nil {
		return err
	}
	*hash = hashBytes
	return nil
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

func NewBlockHashFromString(str string) (BlockHash, error) {
	hash, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return hash, nil
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
		`Block -
	Timestamp : %s
	LastHash  : %s
	Hash      : %s
	Data      : %s
	Nonce     : %d
	Difficulty: %d
`,
		block.Timestamp.Format("2006-01-02 15:04:05.999999"),
		block.LastHash,
		block.Hash,
		block.Data,
		block.Nonce,
		block.Difficulty,
	)
}

func (block Block) checkHash() bool {
	return bytes.Equal(block.Hash, Hash(block.Timestamp, block.LastHash, block.Data, block.Nonce, block.Difficulty))
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

	nonce := uint64(0)

	return Block{
		timestamp,
		lastHash,
		Hash(timestamp, lastHash, data, nonce, DefaultDifficulty),
		data,
		nonce,
		DefaultDifficulty,
	}
}

// MineBlock generates new block based on last block of chain
func MineBlock(lastBlock Block, data BlockData) Block {
	timestamp := time.Now()
	lastHash := lastBlock.Hash
	difficulty := lastBlock.Difficulty

	//hash := Hash(timestamp, lastBlock.Hash, data, nonce)
	timestamp, nonce, blockhash, difficulty := MineHash(lastBlock, data)

	return Block{
		timestamp,
		lastHash,
		blockhash,
		data,
		nonce,
		difficulty,
	}
}

func AdjustDifficulty(lastBlock Block, timestamp time.Time) uint8 {
	difficulty := lastBlock.Difficulty
	timeDiff := timestamp.Sub(lastBlock.Timestamp)
	if timeDiff > (MineRate+MineRateError) && difficulty > 0 {
		difficulty--
	} else if timeDiff < (MineRate - MineRateError) {
		difficulty++
	}
	return difficulty
}

// MineHash finds nonce that produces correct hash
func MineHash(lastBlock Block, data BlockData) (time.Time, uint64, BlockHash, uint8) {
	timestamp := time.Now()
	nonce := uint64(0)
	var blockhash BlockHash
ToMine:
	for {
		timestamp = time.Now()
		difficulty := AdjustDifficulty(lastBlock, timestamp)
		blockhash = Hash(timestamp, lastBlock.Hash, data, nonce, difficulty)

		for i := uint8(0); i < difficulty; i++ {
			if blockhash[i] != 0 {
				nonce++
				continue ToMine
			}
		}
		log.Println("cur diff:", difficulty)
		return timestamp, nonce, blockhash, difficulty
	}
}

// Hash generates hash for new block
func Hash(timestamp time.Time, lastHash BlockHash, data BlockData, nonce uint64, difficulty uint8) BlockHash {
	binTimestamp, err := timestamp.MarshalBinary()
	if err != nil {
		panic(err)
	}
	dataToHash := make([]byte, 0, 1024)
	dataToHash = append(dataToHash, binTimestamp...)
	dataToHash = append(dataToHash, lastHash.MarshalBinary()...)
	dataToHash = append(dataToHash, data.MarshalBinary()...)
	dataToHash = append(dataToHash, difficulty)

	// Adding 8 bytes of placeholder for nonce
	dataToHash = append(dataToHash, 0, 1, 2, 3, 4, 5, 6, 7)

	// This is to speed up hash calculation in future
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	for i, nonceByte := range nonceBytes {
		j := len(dataToHash) - len(nonceBytes) + i
		dataToHash[j] = nonceByte
	}

	blockhash := sha256.Sum256(dataToHash)
	return blockhash[:]
}
