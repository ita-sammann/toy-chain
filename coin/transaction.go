package coin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/ita-sammann/toy-chain"
	"github.com/satori/go.uuid"
)

type Transaction struct {
	ID      uuid.UUID
	Input   TransactionInput
	Outputs TransactionOutputs
}

func NewTransaction(senderWallet Wallet, recipientAddr ecdsa.PublicKey, amount uint64) (*Transaction, error) {
	if amount > senderWallet.Balance {
		return nil, toy_chain.ErrTransactionAmountExceedsBalance
	}

	tx := new(Transaction)
	tx.ID = NewTransactionID()

	tx.Outputs = append(
		tx.Outputs,
		TransactionOutput{
			senderWallet.Balance - amount,
			senderWallet.PublicKey,
		},
		TransactionOutput{
			amount,
			recipientAddr,
		},
	)
	err := senderWallet.SignTransaction(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tx *Transaction) Update(senderWallet Wallet, recipientAddr ecdsa.PublicKey, amount uint64) error {
	var senderIdx int
	for i, txo := range tx.Outputs {
		if txo.address == senderWallet.PublicKey {
			senderIdx = i
			break
		}
	}
	if amount > tx.Outputs[senderIdx].amount {
		return toy_chain.ErrTransactionAmountExceedsBalance
	}
	tx.Outputs[senderIdx].amount -= amount
	tx.Outputs = append(
		tx.Outputs,
		TransactionOutput{
			amount,
			recipientAddr,
		},
	)
	return senderWallet.SignTransaction(tx)
}

func (tx Transaction) VerifySignature(publicKey ecdsa.PublicKey) bool {
	txoHash := sha256.Sum256(tx.Outputs.MarshalBin())
	return ecdsa.Verify(&publicKey, txoHash[:], tx.Input.signature.r, tx.Input.signature.s)
}

type TransactionInput struct {
	timestamp time.Time
	amount    uint64
	address   ecdsa.PublicKey
	signature TXISignature
}

type TransactionOutput struct {
	amount  uint64
	address ecdsa.PublicKey
}

func (txo TransactionOutput) MarshalBin() []byte {
	txoBin := make([]byte, 8, 1024)
	binary.BigEndian.PutUint64(txoBin, txo.amount)
	txoBin = append(txoBin, elliptic.Marshal(txo.address.Curve, txo.address.X, txo.address.Y)...)
	return txoBin
}

type TransactionOutputs []TransactionOutput

func (txos TransactionOutputs) MarshalBin() []byte {
	txosBin := make([]byte, 0, 1024)
	for _, txo := range txos {
		txosBin = append(txosBin, txo.MarshalBin()...)
	}
	return txosBin
}

func NewTransactionID() uuid.UUID {
	id, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}
	return id
}

type TXISignature struct {
	r, s *big.Int
}

func NewTXISignature(r, s *big.Int) TXISignature {
	return TXISignature{r, s}
}
