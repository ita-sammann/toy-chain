package coin

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wallet", func() {
	var (
		wallet Wallet
	)

	BeforeEach(func() {
		wallet = NewWallet()
	})

	It("has correct initial balance", func() {
		Expect(wallet.Balance).To(Equal(uint64(InitialBalance)))
	})

	Context("When marshalled to json", func() {
		var (
			jsonWallet         []byte
			unmarshalledWallet map[string]interface{}
			err, err2          error
		)

		BeforeEach(func() {
			jsonWallet, err = wallet.MarshalJSON()
			err2 = json.Unmarshal(jsonWallet, &unmarshalledWallet)
		})

		It("produces no errors", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("has correct balance", func() {
			balance, ok := unmarshalledWallet["balance"].(float64)

			Expect(ok).To(BeTrue())
			Expect(uint64(balance)).To(Equal(uint64(InitialBalance)))
		})

		It("has correct public key", func() {
			publicKey, ok := unmarshalledWallet["publicKey"].(string)

			Expect(ok).To(BeTrue())
			Expect(publicKey).To(MatchRegexp(`^[0-9a-zA-Z+/=]{88}$`))
		})
	})
})
