package blockchain

import (
	"bytes"

	"github.com/ita-sammann/toy-chain"
)

// Blockchain is blockchain
type Blockchain struct {
	chain []Block
}

// NewBlockchain is constructor
func NewBlockchain() Blockchain {
	return Blockchain{
		[]Block{Genesis()},
	}
}

//NewBlockchainBlocks creates new blockchain from blocks slice
func NewBlockchainBlocks(chain []Block) Blockchain {
	return Blockchain{
		chain,
	}
}

// Len returns length of blockchain
func (bc Blockchain) Len() int {
	return len(bc.chain)
}

// LastBlock returns last block of chain
func (bc Blockchain) LastBlock() Block {
	return bc.chain[bc.Len()-1]
}

// AddBlock adds new block to blockchain
func (bc *Blockchain) AddBlock(data BlockData) Block {
	block := MineBlock(bc.LastBlock(), data)
	bc.chain = append(bc.chain, block)
	return block
}

// IsValid checks validity of current blockchain
func (bc Blockchain) IsValid() bool {
	if bc.chain[0].String() != Genesis().String() {
		return false
	}

	for i := 1; i < bc.Len(); i++ {
		block := bc.chain[i]
		lastBlock := bc.chain[i-1]
		if !bytes.Equal(block.LastHash, lastBlock.Hash) || !block.checkHash() {
			return false
		}
	}

	return true
}

// ReplaceChain replaces current blockchain with new one if it's valid
func (bc *Blockchain) ReplaceChain(newChain Blockchain) error {
	if newChain.Len() <= bc.Len() {
		return toy_chain.ErrChainReplaceTooShort
	}
	if !newChain.IsValid() {
		return toy_chain.ErrChainReplaceInvalid
	}

	bc.chain = newChain.chain
	return nil
}

// ListBlocks returns slice of blocks in chain
func (bc Blockchain) ListBlocks() []Block {
	return bc.chain
}
