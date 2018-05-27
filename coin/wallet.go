package coin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

const (
	InitialBalance = 500
)

func NewKeyPair() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to generate private key"))
	}
	return privateKey
}

type Wallet struct {
	Balance    uint64 `json:"balance"`
	privateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey `json:"publicKey"`
}

func NewWallet() Wallet {
	privateKey := NewKeyPair()
	return Wallet{
		InitialBalance,
		privateKey,
		privateKey.PublicKey,
	}
}

func (wallet Wallet) String() string {
	return fmt.Sprintf(
		`Wallet -
	PublicKey: %s
	Balance  : %d
`,
		base64.StdEncoding.EncodeToString(elliptic.Marshal(wallet.PublicKey.Curve, wallet.PublicKey.X, wallet.PublicKey.Y)),
		wallet.Balance,
	)
}

func (wallet Wallet) MarshalJSON() ([]byte, error) {
	jsonWallet := make(map[string]interface{})
	jsonWallet["balance"] = wallet.Balance
	jsonWallet["publicKey"] = base64.StdEncoding.EncodeToString(elliptic.Marshal(wallet.PublicKey.Curve, wallet.PublicKey.X, wallet.PublicKey.Y))
	return json.Marshal(jsonWallet)
}

func (wallet Wallet) SignTransaction(tx *Transaction) error {
	txoHash := sha256.Sum256(tx.Outputs.MarshalBin())
	r, s, err := ecdsa.Sign(rand.Reader, wallet.privateKey, txoHash[:])
	if err != nil {
		return errors.Wrap(err, "Failed to sign transaction")
	}

	txi := TransactionInput{
		time.Now().UTC(),
		wallet.Balance,
		wallet.PublicKey,
		NewTXISignature(r, s),
	}
	tx.Input = txi
	return nil
}
