# zapredactor

Redact structs when logging them to zap.

## Getting Started

This library uses code generation over annotated structs to redact informations using the zap library. All fields are
redacted by default unless you explicitly mark it.

To log any information as redacted you need to use the `zapredactor.Redact` function. It returns a `zap.Field` and is
designed to work just like `zap.String` (for example).


### Tagging fields

Field tags are easy to use. You just need to define the tag `redact:"[field_name|],[allow|],[redactor]"`.

* **field_name** is the name of the field you want to redact. If not specified, the annotated `json` name is used. 
  Otherwise, the Go field name is used.
* **allow** is a boolean that mark this field not to be redacted. All fields are redacted by default unless marked with
  this option.
* **redactor** is the name of the redactor to use. If not specified, the default redactor is used. Check [here](#Available_redactors) the list 
  of redactors available.

```go
package main

import (
	"github.com/jamillosantos/zapredactor"
	"go.uber.org/zap
)

type CreditCard struct {
	PAN            string `redact:",,pan64"` // This will be redacted using the zaparray.PAN64 redactor (first 6, last 4).
	CardHolder     string
	CVV            string `redact:",,*"` // This will be masked with `*`.
	Expiry         string
	First4         string `redact:",allow"`
	Last4          string `redact:",allow"`
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

	logger.Info("card added to the customer", zap.String("customer_id", "123"), zapredactor.Redact("card", card))
}
```

This method has a limitation: You won't be able to log a struct that you cannot annotate.

> You will find the complete example on `/examples/tag/main.go`.

## Code generation

First, you will need to install the `zapredactor` CLI.

```bash
go install github.com/jamillosantos/zapredactor/cli/zapredactor@latest
```

Then, you can generate the code for your structs:

```bash
zapredactor <package> --destination=<file name>.go
```

The command above will find all the structs on the given package that contains the `redact` tag on it. Then, generate a
file named `<file name>.go` with the `zapredactor.Redactable` interface.

Alternatively, you can use `go:generate` directive:

```go
//go:generate zapredactor . --destination=redactor_gen.go
package domain

type Card struct {
	ID         string `redact:",allow"`
	PAN        string `redact:",,pan64"`
	CVV        string
	CardHolder string
}
```

## Available redactors

| Name            | Description                                                                   | Types compatible                                 |
|-----------------|-------------------------------------------------------------------------------|--------------------------------------------------|
| (empty string)  | Returns `[redacted]` regardless the input.                                    | `any`                                            |
| `pan`           | Redacts a card number outputting only its first 6 and last 4 digits of a PAN. | `string`, `*string`                              |
| `bin`           | Redacts a card number outputting only its first 6.                            | `string`, `*string`                              |
| `star`, `*`     | Redacts the a given string (or *string).                                      | `string`, `*string`                              |
| `len`           | Redacts the first digit of a PAN.                                             | `string`, `*string`, arrays (through reflection) |

### Anatomy of a redactor

A redactor is a simple function that receives an value and returns a string:

```go
package myredactors

// ...

func MyRedactor(_ interface{}) string {
    return "[redacted]"
}
```

Now, you need to inform zapredactor your redactor is available to be used. Use the `zapredactor.RegisterRedactor` function
to bind its key to the given  redactor. Then, you will be able to use it in the `redact` tag.

Usually you will use the `init` function to register your redactors.

```go
package myredactors

// ...

func init() {
    zapredactor.RegisterRedactor("my_redactor", MyRedactor)
}
```

The you will be able to use it in the `redact` tag.

```go
// ...

type MyStruct struct {
    Field string `redact:",,my_redactor"`
}
```

# License

MIT license
