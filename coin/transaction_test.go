package coin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	testKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	recipientAddr := testKey.PublicKey
	txAmount := uint64(50)
	wallet := NewWallet()

	tx, err := NewTransaction(wallet, recipientAddr, txAmount+wallet.Balance)
	if err == nil || tx != nil {
		t.Error("Managed to create TX with excessive amount")
	}

	tx, err = NewTransaction(wallet, recipientAddr, txAmount)
	if err != nil {
		t.Error("Failed to create correct TX")
		return
	}

	for _, txo := range tx.Outputs {
		if txo.address == recipientAddr {
			if txo.amount != txAmount {
				t.Errorf("Bad `change` TXO amount. Got: %d, expected: %d", txo.amount, wallet.Balance-txAmount)
			}
		} else if txo.address == wallet.PublicKey {
			if txo.amount != wallet.Balance-txAmount {
				t.Errorf("Bad `change` TXO amount. Got: %d, expected: %d", txo.amount, wallet.Balance-txAmount)
			}
		} else {
			t.Error("Unexpected TXO")
		}
	}

	if tx.Input.amount != wallet.Balance {
		t.Errorf("Bad TXI amount. Got: %d, expected: %d", tx.Input.amount, wallet.Balance)
	}
}

func TestTransaction_VerifySignature(t *testing.T) {
	testKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	recipientAddr := testKey.PublicKey
	txAmount := uint64(50)
	wallet := NewWallet()

	tx, err := NewTransaction(wallet, recipientAddr, txAmount)
	if err != nil {
		t.Error("Failed to create TX")
		return
	}

	if !tx.VerifySignature(wallet.PublicKey) {
		t.Error("Failed to validate correct TX")
	}

	tx.Outputs[0].amount = wallet.Balance
	if tx.VerifySignature(wallet.PublicKey) {
		t.Error("Validated corrupt TX")
	}
}
