package blockchain

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
