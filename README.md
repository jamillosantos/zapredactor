# zapredactor

Redact structs when logging them to zap.

## Getting Started

This library redacts all fields by default, unless you actively tell the redactor not to do so.

To log any information as redacted you need to use the `zapredactor.Redact` function. It returns a `zap.Field` and is
designed to work just like `zap.String` (for example).

### Redacting with field tags

Field tags are easy to use. You just need to define the tag `redact:"allow"`.

```go
package main

import (
	"github.com/jamillosantos/zapredactor"
	"go.uber.org/zap
)

type CreditCard struct {
	PAN            string
	CardHolder     string
	CVV            string
	Expiry         string
	First4         string `redact:"allow"`
	Last4          string `redact:"allow"`
	IntenralStatus string `redact:"-"` // This field won't make to the log.
}

func main() {
	card := &CreditCard{
		PAN:            "1234123412344321",
		CardHolder:     "Snake Eyes Joe",
		CVV:            "123",
		Expiry:         "02/29",
		First4:         "1234",
		Last4:          "4321",
		IntenralStatus: "",
	}

	var logger *zap.Logger

	// You need to initialize logger before next line.

	logger.Info("card added to the customer", zap.String("customer_id"), zapredactor.Redact("card", card))
}
```

This method has a limitation: You won't be able to log a struct that you cannot annotate.

> You will find the complete example on `/examples/tag/main.go`.

## Using non-intrusive methods

> **⚠ WARNING: This features is not implemented yet.**

# License

MIT license
