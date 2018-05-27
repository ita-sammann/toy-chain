package blockchain

import (
	"github.com/ita-sammann/toy-chain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Blockchain", func() {
	var (
		bc Blockchain
	)

	BeforeEach(func() {
		bc = NewBlockchain()
		bc.AddBlock([]byte("test_data_1"))
		bc.AddBlock([]byte("test_data_2"))
		bc.AddBlock([]byte("test_data_3"))
	})

	It("starts with genesis block", func() {
		Expect(bc.ListBlocks()[0].String()).To(Equal(Genesis().String()))
	})

	It("has correct data in last block", func() {
		Expect(string(bc.ListBlocks()[3].Data)).To(Equal("test_data_3"))
	})

	It("is valid", func() {
		Expect(bc.IsValid()).To(BeTrue())
	})

	It("is invalid with corrupt genesis block", func() {
		bc.chain[0].Data = []byte("Corrupt data")
		Expect(bc.IsValid()).To(BeFalse())
	})

	It("is invalid with corrupt intermediate block", func() {
		bc.chain[1].Data = []byte("Corrupt data")
		Expect(bc.IsValid()).To(BeFalse())
	})

	It("is invalid with corrupt last block", func() {
		bc.chain[bc.Len()-1].Data = []byte("Corrupt data")
		Expect(bc.IsValid()).To(BeFalse())
	})

	Describe("Replacing with another chain", func() {
		BeforeEach(func() {
			bc = NewBlockchain()
			bc.AddBlock([]byte("data_1"))
			bc.AddBlock([]byte("data_2"))
		})

		It("should be replaced with valid and longer chain", func() {
			bc2 := NewBlockchain()
			bc2.AddBlock([]byte("foo 1"))
			bc2.AddBlock([]byte("foo 2"))
			bc2.AddBlock([]byte("bar 3"))

			err := bc.ReplaceChain(bc2)

			Expect(err).NotTo(HaveOccurred())
			Expect(bc.LastBlock().Hash).To(Equal(bc2.LastBlock().Hash))
		})

		It("should not be replaced with shorter one", func() {
			bc2 := NewBlockchain()
			bc2.AddBlock([]byte("baz 1"))
			bc2.AddBlock([]byte("baz 2"))

			err := bc.ReplaceChain(bc2)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(toy_chain.ErrChainReplaceTooShort))
			Expect(bc.LastBlock().Hash).NotTo(Equal(bc2.LastBlock().Hash))
		})

		It("should not be replaced with corrupt one", func() {
			bc2 := NewBlockchain()
			bc2.AddBlock([]byte("foobar 1"))
			bc2.AddBlock([]byte("foobar 2"))
			bc2.AddBlock([]byte("foobar 3"))
			bc2.AddBlock([]byte("foobar 4"))
			bc2.AddBlock([]byte("foobar 5"))
			bc2.AddBlock([]byte("foobar 6"))
			bc2.AddBlock([]byte("foobar 7"))
			bc2.AddBlock([]byte("foobar 8"))
			bc2.chain[5].Data = []byte("corrupt data")

			err := bc.ReplaceChain(bc2)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(toy_chain.ErrChainReplaceInvalid))
			Expect(bc.LastBlock().Hash).NotTo(Equal(bc2.LastBlock().Hash))
		})
	})
})
