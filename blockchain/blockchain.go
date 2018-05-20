package blockchain

import "bytes"

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

// AddBlock adds new block to blockchain
func (blockchain *Blockchain) AddBlock(data BlockData) Block {
	block := MineBlock(blockchain.chain[len(blockchain.chain)-1], data)
	blockchain.chain = append(blockchain.chain, block)
	return block
}

func (blockchain Blockchain) isValid() bool {
	if blockchain.chain[0].String() != Genesis().String() {
		return false
	}

	for i, block := range blockchain.chain[1:] {
		lastBlock := blockchain.chain[i-1]
		if !bytes.Equal(block.LastHash, lastBlock.Hash) || !block.checkHash() {
			return false
		}
	}

	return true
}
