package zapredactor

import "go.uber.org/zap"

type CreditCard struct {
	PAN            string
	CardHolder     string
	CVV            string
	Expiry         string
	First4         string `redact:"allow"`
	Last4          string `redact:"allow"`
	IntenralStatus string `redact:"-"` // This field won't make to the log.
}

func ExampleRedact() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	card := &CreditCard{
		PAN:            "1234123412344321",
		CardHolder:     "Snake Eyes Joe",
		CVV:            "123",
		Expiry:         "02/29",
		First4:         "1234",
		Last4:          "4321",
		IntenralStatus: "pending",
	}

	logger.Info("card added to customer", Redact("card", card))
}