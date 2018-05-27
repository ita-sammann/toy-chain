package coin

import (
	"encoding/json"
	"regexp"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()

	if wallet.Balance != InitialBalance {
		t.Error("Incorrect initial balance")
	}
}

func TestWallet_MarshalJSON(t *testing.T) {
	jsonWallet, err := NewWallet().MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	var wallet map[string]interface{}
	err = json.Unmarshal(jsonWallet, &wallet)
	if err != nil {
		t.Error(err)
	}

	balance, ok := wallet["balance"].(float64)
	if ok {
		if int(balance) != InitialBalance {
			t.Errorf("invalid balance: got %d, expected %d", int(balance), InitialBalance)
		}
	} else {
		t.Error("balance is not uint64")
	}

	publicKey, ok := wallet["publicKey"].(string)
	if ok {
		if !regexp.MustCompile(`^[0-9a-zA-Z+/=]{88}$`).MatchString(publicKey) {
			t.Errorf("Bad public key format: %s", publicKey)
		}
	} else {
		t.Error("public key is not a string")
	}
}

//func TestWallet_SignTransaction(t *testing.T) {
//	wallet1 := NewWallet()
//	wallet2 := NewWallet()
//	tx, err := NewTransaction(wallet1, wallet2.PublicKey, 50)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	wallet1.SignTransaction(tx)
//}
