//go:generate go run ../../cli/zapredactor/main.go --destination redactor_gen.go
package main

import (
	"go.uber.org/zap"

	"github.com/jamillosantos/zapredactor"
)

type Expiry struct {
	Month int
	Year  int `redact:",allow"`
}

type CreditCard struct {
	PAN            string `json:"pan"`
	CardHolder     string `redact:"card_holder"`
	CVV            string
	Expiry         Expiry
	First4         string `redact:",allow"`
	Last4          string `redact:",allow"`
	IntenralStatus string `redact:"-"` // This field won't make to the log.
}

type Address struct {
	Street string `redact:","`
	City   string
}

type Customer struct {
	Name    string      `redact:""`
	Card    *CreditCard `redact:""`
	Address Address     `redact:""`
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	customerWithoutCard := Customer{
		Name: "Customer without Card",
		Address: Address{
			Street: "Rua do Professor",
			City:   "Natal",
		},
	}

	customerWithCard := Customer{
		Name: "Customer without Card",
		Card: &CreditCard{
			PAN:        "1234123412344321",
			CardHolder: "Snake Eyes Joe",
			CVV:        "123",
			Expiry: Expiry{
				Month: 2,
				Year:  29,
			},
			First4:         "1234",
			Last4:          "4321",
			IntenralStatus: "pending",
		},
		Address: Address{
			Street: "Rua do Professor",
			City:   "Natal",
		},
	}

	logger.Info("customer without card", zapredactor.Redact("customer", &customerWithoutCard))
	logger.Info("card added to customer", zapredactor.Redact("customer", &customerWithCard))
}
