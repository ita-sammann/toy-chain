package blockchain

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Block", func() {

	It("is is correctly converted to string", func() {
		block := Block{
			time.Date(2018, 5, 20, 15, 30, 45, 123456789, time.UTC),
			[]byte("abcdefghijklmnopqrstuvwxyz"),
			[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			[]byte("this is some data in block"),
			uint64(0),
			2,
		}
		expectedString := `Block -
	Timestamp : 2018-05-20 15:30:45.123456
	LastHash  : 6162636465666768696a6b6c6d6e6f707172737475767778797a
	Hash      : 4142434445464748494a4b4c4d4e4f505152535455565758595a
	Data      : this is some data in block
	Nonce     : 0
	Difficulty: 2
`

		Expect(block.String()).To(Equal(expectedString))
	})

	It("is correct genesis block", func() {
		timestamp := time.Unix(1526774400, 0)
		lastHash := []byte{
			0, 0, 0, 0, 255, 255, 255, 255,
			0, 0, 0, 0, 255, 255, 255, 255,
			0, 0, 0, 0, 255, 255, 255, 255,
			0, 0, 0, 0, 255, 255, 255, 255,
		}
		data := []byte("genesis")

		nonce := uint64(0)

		expectedBlock := Block{
			timestamp,
			lastHash,
			Hash(timestamp, lastHash, data, nonce, DefaultDifficulty),
			data,
			nonce,
			DefaultDifficulty,
		}

		Expect(Genesis().String()).To(Equal(expectedBlock.String()))
	})

	It("is mined correctly", func() {
		newBlock := MineBlock(Genesis(), []byte("test_data"))

		Expect(newBlock.LastHash).To(Equal(Genesis().Hash))
		Expect(newBlock.Hash).To(Equal(
			Hash(newBlock.Timestamp, newBlock.LastHash, newBlock.Data, newBlock.Nonce, newBlock.Difficulty),
		))
		Expect(string(newBlock.Data)).To(Equal("test_data"))
	})
})
