# zapredactor

Redact structs when logging them to zap.

## Getting Started

This library uses code generation over annotated structs to redact informations using the zap library. All fields are
redacted by default unless you explicitly mark it.

To log any information as redacted you need to use the `zapredactor.Redact` function. It returns a `zap.Field` and is
designed to work just like `zap.String` (for example).

## Redacting fields with Code Generation

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

## Redacting fields with Reflection

This method uses reflection to iterate over all fields, redacting whenever needd. It is indicated when you need to log a
struct that you cannot annotate (eg. an struct from a third-party library).

As the other method, all fields are redacted by default. For a field to be not redacted, you need to explicitly inform
the redactor.

__Example 1__: The following example will redact with the definitions predefined (examples/reflectiondemo/main.go).

```go
package main

import (
  "go.uber.org/zap"

  "github.com/jamillosantos/zapredactor/redactors"
  "github.com/jamillosantos/zapredactor/redactreflection"
)

type Demo struct {
  String string
  Int    int
  PAN    string
}

var (
  // redactor can be reused multiple times.
  redactor = redactreflection.Redactor(
    redactreflection.WithAllowFields("String"),
    redactreflection.WithRedactor("PAN", redactors.PAN64),
  )
)

func main() {
  // Initialize logger
  logger, err := zap.NewDevelopment()
  if err != nil {
    panic(err)
  }

  logger.Info("demo entry reusing the redactor configuration 1", redactor("demo", &Demo{
    String: "string",
    Int:    1,
    PAN:    "1234567890124321",
  }))
  // Output: {"demo": {"String": "string", "Int": "[redacted]", "PAN": "123456******4321"}}
  logger.Info("demo entry reusing the redactor configuration 2", redactor("demo", &Demo{
    String: "string 2",
    Int:    2,
    PAN:    "1111112222223333",
  }))
  // Output: {"demo": {"String": "string 2", "Int": "[redacted]", "PAN": "111111******3333"}}
}

```

__Example 2__: The following example will redact with the definitions inline. So, for each call you need to specify the 
redactor configuration (examples/reflectiondemo2/main.go).

```go
package main

import (
  "go.uber.org/zap"

  "github.com/jamillosantos/zapredactor/redactors"
  "github.com/jamillosantos/zapredactor/redactreflection"
)

type Demo struct {
  String string
  Int    int
  PAN    string
}

func main() {
  // Initialize logger
  logger, err := zap.NewDevelopment()
  if err != nil {
    panic(err)
  }

  demo := Demo{
    String: "string",
    Int:    1,
    PAN:    "1234567890124321",
  }

  logger.Info("demo entry with all fields redacted", redactreflection.Redact("demo", &demo))
  // Output: {"demo": {"String": "[redacted]", "Int": "[redacted]", "PAN": "[redacted]"}}
  logger.Info("demo entry with one field redacted", redactreflection.Redact("demo", &demo, redactreflection.WithAllowFields("String")))
  // Output: {"demo": {"String": "string", "Int": "[redacted]", "PAN": "[redacted]"}}
  logger.Info("demo entry with a custom redactor", redactreflection.Redact("demo", &demo, redactreflection.WithRedactor("PAN", redactors.PAN64)))
  // Output: {"demo": {"String": "[redacted]", "Int": "[redacted]", "PAN": "123456******4321"}}
  logger.Info("demo entry hiding a field", redactreflection.Redact("demo", &demo, redactreflection.WithHiddenField("PAN")))
  // Output: {"demo": {"String": "[redacted]", "Int": "[redacted]"}}
}
```

## Redactors

| Name            | Code              | Description                                                                   | Types compatible                                 |
|-----------------|-------------------|-------------------------------------------------------------------------------|--------------------------------------------------|
| (empty string)  | `redactors.Default` | Returns `[redacted]` regardless the input.                                    | `any`                                            |
| `pan`           | `redactors.PAN64`   | Redacts a card number outputting only its first 6 and last 4 digits of a PAN. | `string`, `*string`                              |
| `bin`           | `redactors.BIN`     | Redacts a card number outputting only its first 6.                            | `string`, `*string`                              |
| `star`, `*`     | `redactors.Star`    | Redacts the a given string (or *string).                                      | `string`, `*string`                              |
| `len`           | `redactors.Len`     | Redacts the first digit of a PAN.                                             | `string`, `*string`, arrays (through reflection) |

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
