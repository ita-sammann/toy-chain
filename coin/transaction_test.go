package coin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction", func() {
	var (
		recipientAddr ecdsa.PublicKey
		txAmount      uint64
		wallet        Wallet
		tx            *Transaction
		err           error
	)

	BeforeEach(func() {
		testKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		recipientAddr = testKey.PublicKey
		txAmount = 50
		wallet = NewWallet()
		tx, err = NewTransaction(wallet, recipientAddr, txAmount)
	})

	It("fails to create with excessive amount", func() {
		tx, err = NewTransaction(wallet, recipientAddr, txAmount+wallet.Balance)

		Expect(tx).To(BeNil())
		Expect(err).To(HaveOccurred())
	})

	It("is successfully created with correct data", func() {
		Expect(err).NotTo(HaveOccurred())
		Expect(tx.Input.amount).To(Equal(wallet.Balance), "bad TXI amount")

		for _, txo := range tx.Outputs {
			if txo.address == recipientAddr {
				Expect(txo.amount).To(Equal(txAmount), "bad `target` TXO amount")
			} else if txo.address == wallet.PublicKey {
				Expect(txo.amount).To(Equal(wallet.Balance-txAmount), "bad `change` TXO amount.")
			} else {
				Fail("unexpected TXO")
			}
		}

	})

	It("is correct and is verified successfully", func() {
		Expect(tx.VerifySignature(wallet.PublicKey)).To(BeTrue())
	})

	It("is corrupt and is not verified", func() {
		tx.Outputs[0].amount = wallet.Balance

		Expect(tx.VerifySignature(wallet.PublicKey)).To(BeFalse())
	})

	Describe("updating a transaction", func() {
		var (
			nextAmount    uint64
			nextRecipient ecdsa.PublicKey
		)

		BeforeEach(func() {
			nextAmount = 20
			testKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			nextRecipient = testKey.PublicKey
			err = tx.Update(wallet, nextRecipient, nextAmount)
		})

		It("subtracts next amount from sender's output", func() {
			var senderTxo TransactionOutput
			for _, txo := range tx.Outputs {
				if txo.address == wallet.PublicKey {
					senderTxo = txo
					break
				}
			}
			Expect(senderTxo.amount).To(Equal(wallet.Balance - txAmount - nextAmount))
		})

		It("outputs amount for next recipient", func() {
			var nextTxo TransactionOutput
			for _, txo := range tx.Outputs {
				if txo.address == nextRecipient {
					nextTxo = txo
					break
				}
			}
			Expect(nextTxo.amount).To(Equal(nextAmount))
		})
	})
})
