# zapredactor

Log structs and maps with [zap](https://github.com/uber-go/zap) without leaking sensitive data.

Every field is redacted by default. A value only reaches the log when you explicitly allow it.

```go
logger.Info("card added", zapredactor.Redact("card", &card))
// {"card": {"pan": "123412******4321", "card_holder": "[redacted]", "CVV": "***", "Expiry": {"Month": "[redacted]", "Year": 29}, "Last4": "4321"}}
```

## Contents

- [Installation](#installation)
- [Choosing an approach](#choosing-an-approach)
- [Code generation](#code-generation)
- [Reflection](#reflection)
- [Redactors](#redactors)
- [Examples](#examples)
- [Development](#development)
- [License](#license)

## Installation

```bash
go get github.com/jamillosantos/zapredactor
```

Requires Go 1.26 or newer.

## Choosing an approach

There are two ways to redact a value. Both return a `zap.Field`, so they are used like `zap.String` or `zap.Object`.

|                                    | [Code generation](#code-generation)   | [Reflection](#reflection)                 |
|------------------------------------|---------------------------------------|-------------------------------------------|
| Performance                        | faster, no reflection when logging    | slower, reflects on every log call        |
| Configured with                    | `redact` struct tags                  | options passed when logging               |
| Works with types you cannot edit   | no                                    | yes                                       |
| Redacts maps key by key            | no                                    | yes                                       |
| Extra build step                   | yes, run the generator                | no                                        |
| Entry point                        | `zapredactor.Redact`                  | `redactreflection.Redact`                 |

**Code generation is the more performant option.** The generated method writes each field straight to the zap
encoder, while reflection inspects the value again on every log call. Logging a small struct was about three times
faster with code generation, using roughly a third of the memory and allocations.

Use code generation for your own types, and whenever logging sits on a hot path. Keep reflection for types from other
packages and for maps.

## Code generation

You annotate a struct with `redact` tags, run the `zapredactor` generator, and it writes a `Redact` method for that
struct. `zapredactor.Redact` then logs any value that has this method.

### 1. Tag the struct

```go
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
	PAN            string `json:"pan" redact:",,pan64"` // first 6 and last 4 digits
	CardHolder     string `redact:"card_holder"`        // redacted, logged as "card_holder"
	CVV            string `redact:",,*"`                // one "*" per character
	Expiry         Expiry                               // nested struct, follows its own tags
	Last4          string `redact:",allow"`             // logged as is
	InternalStatus string `redact:"-"`                  // never logged
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	card := CreditCard{
		PAN:            "1234123412344321",
		CardHolder:     "Snake Eyes Joe",
		CVV:            "123",
		Expiry:         Expiry{Month: 2, Year: 29},
		Last4:          "4321",
		InternalStatus: "pending",
	}

	logger.Info("card added", zap.String("customer_id", "123"), zapredactor.Redact("card", &card))
}
```

Output:

```
{"customer_id": "123", "card": {"pan": "123412******4321", "card_holder": "[redacted]", "CVV": "***", "Expiry": {"Month": "[redacted]", "Year": 29}, "Last4": "4321"}}
```

#### The `redact` tag

```
redact:"<name>,<allow>,<redactor>"
```

All three parts are optional.

| Part       | Meaning                                                                                                              |
|------------|----------------------------------------------------------------------------------------------------------------------|
| `name`     | Key used in the log. Defaults to the name in the `json` tag, then to the Go field name. `-` removes the field from the log. |
| `allow`    | The literal word `allow`. The field is logged with its real value.                                                   |
| `redactor` | Name of the [redactor](#redactors) applied to the field. Ignored when the field is allowed.                          |

| Tag                     | Result                                                        |
|-------------------------|---------------------------------------------------------------|
| no tag                  | Redacted as `[redacted]`.                                     |
| `redact:"card_holder"`  | Redacted as `[redacted]`, under the key `card_holder`.        |
| `redact:",allow"`       | Logged with its real value.                                   |
| `redact:",,pan64"`      | Redacted with the `pan64` redactor.                           |
| `redact:"pan,,pan64"`   | Redacted with the `pan64` redactor, under the key `pan`.      |
| `redact:"-"`            | Not logged.                                                   |

The generator only picks up structs that have **at least one field with a `redact` tag**. To get a struct where every
field is redacted, tag one field with an empty tag: `redact:""`.

### 2. Generate the code

Add the generator to your `go.mod` as a tool. Its version then always matches the library version, and nobody has to
install anything by hand:

```bash
go get -tool github.com/jamillosantos/zapredactor/cli/zapredactor
```

Then run it:

```bash
go tool zapredactor [package] --destination redactor_gen.go
```

| Argument              | Meaning                                                                                                   |
|-----------------------|-----------------------------------------------------------------------------------------------------------|
| `package`             | Package to read. Defaults to `.`. Accepts a relative path starting with `./` or a full import path.       |
| `--destination`, `-d` | File to write. Defaults to `-`, which prints the code to the standard output.                             |

The usual setup is a `go:generate` directive next to the structs:

```go
//go:generate go tool zapredactor --destination redactor_gen.go
package domain
```

Run `go generate ./...` again whenever a tagged struct changes.

#### Installing the binary instead

If you prefer a standalone binary, install it and call `zapredactor` directly. The arguments are the same.

```bash
go install github.com/jamillosantos/zapredactor/cli/zapredactor@latest
```

```go
//go:generate zapredactor --destination redactor_gen.go
package domain
```

### 3. Log it

The generated file gives each struct two methods:

- `Redact(zapcore.ObjectEncoder) error`, which satisfies the `zapredactor.Redactable` interface.
- `MarshalLogObject(zapcore.ObjectEncoder) error`, which calls `Redact` and satisfies zap's own
  `zapcore.ObjectMarshaler` interface.

So the value, or a pointer to it, can be logged with `zapredactor.Redact` or with `zap.Object`:

```go
logger.Info("card added", zapredactor.Redact("card", &card))
logger.Info("card added", zap.Object("card", card))
```

`zapredactor.Redact` logs a nil value as `null`. `zap.Object` panics on a nil pointer, so give it a value or a pointer
you know is set.

`zap.Any` and the sugared logger (`Infow`) redact as well, but **only when the value itself is the generated type**. A
slice, a map or a struct that contains it goes through reflection and is logged with every field in clear, and so is
anything passed to `zap.Reflect`. Log a slice with `zap.Objects("cards", cards)`.

You can also write the `Redact` method by hand for types the generator cannot handle. Add a `MarshalLogObject` method
that calls it if the type should work with `zap.Object` too.

### How fields are handled

- **Redacted fields** (the default) can have any type. The value is handed to the redactor and logged as a string.
- **Allowed fields** must have one of these types: `string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`,
  `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64`, `complex128`, `time.Time` or
  `time.Duration`, as a value (`T`), a pointer (`*T`) or a slice (`[]T`), plus `[]byte`. A nil pointer is logged as
  `null`. Any other type, including named types such as `type Status string`, stops the generator with
  `unsupported field type`.
- **Struct fields** whose type is declared in the same package, as a value or as a pointer, are logged as a nested
  object through their own `Redact` method. That type needs a generated (or handwritten) `Redact` method too. The
  `allow` and `redactor` parts of the tag are ignored for these fields.
- **Struct types from other packages**, such as `url.URL`, are treated like any other value: redacted as a whole.
- **Unexported fields** are included and follow the same rules.

### Known limitations

- A field that is a slice of structs from the same package (`[]Inner`, `[]*Inner`) or an embedded struct is not
  supported. The generated code does not compile, or the generator fails. Tag these fields with `redact:"-"`.
- When every field of every generated struct in a package is allowed or skipped, the generated file has an unused
  import and does not compile. At least one field must be redacted or be a nested struct.
- An allowed field that is a fixed-size array (`[4]string`) is read as a slice, and the generated code does not
  compile. Use a slice, or tag the field with `redact:"-"`.
- The generator reports errors by panicking.

## Reflection

`redactreflection` walks the value with reflection every time it is logged. It needs no tags and no generated code,
so it works with types from other packages and with maps.

```go
package main

import (
	"go.uber.org/zap"

	"github.com/jamillosantos/zapredactor/redactors"
	"github.com/jamillosantos/zapredactor/redactreflection"
)

type Address struct {
	Street string
	City   string `json:"city"`
}

type Customer struct {
	Name    string
	PAN     string
	Token   string
	Address Address
}

// Build the configuration once and reuse it.
var redactCustomer = redactreflection.Redactor(
	redactreflection.WithAllowFields("Name", "Address.city"),
	redactreflection.WithRedactor("PAN", redactors.PAN64),
	redactreflection.WithHiddenField("Token"),
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	customer := Customer{
		Name:    "Joe",
		PAN:     "1234567890124321",
		Token:   "secret",
		Address: Address{Street: "Rua do Professor", City: "Natal"},
	}

	logger.Info("reused configuration", redactCustomer("customer", &customer))
	// {"customer": {"Name": "Joe", "PAN": "123456******4321", "Address": {"Street": "[redacted]", "city": "Natal"}}}

	// Or pass the options on each call.
	logger.Info("options per call", redactreflection.Redact("customer", &customer,
		redactreflection.WithAllowFields("Name"),
	))
	// {"customer": {"Name": "Joe", "PAN": "[redacted]", "Token": "[redacted]", "Address": {"Street": "[redacted]", "city": "[redacted]"}}}

	// Maps work the same way.
	logger.Info("map", redactreflection.Redact("payload", map[string]any{"id": 10, "password": "secret"},
		redactreflection.WithAllowFields("id"),
	))
	// {"payload": {"id": 10, "password": "[redacted]"}} (map keys are logged in no particular order)
}
```

### Functions

| Function                                       | Use                                                                  |
|------------------------------------------------|----------------------------------------------------------------------|
| `redactreflection.Redact(key, value, opts...)` | Returns a `zap.Field` for one value.                                 |
| `redactreflection.Redactor(opts...)`           | Returns a `func(key, value) zap.Field` that reuses the same options. |

`value` must be a struct, a map with string keys, or a pointer to one of them. Anything else, including a nil pointer,
is logged as an empty object.

### Options

| Option                             | Effect                                                                                  |
|------------------------------------|-----------------------------------------------------------------------------------------|
| `WithAllowFields(paths...)`        | Logs the fields with their real values.                                                 |
| `WithRedactor(path, redactor)`     | Redacts the field with the given [redactor](#redactors) instead of the default one.     |
| `WithHiddenField(paths...)`        | Removes the fields from the log.                                                        |
| `FieldNameExtractor(fn)`           | Replaces the function that names struct fields. The default is `DefaultNameExtractor`.  |

Each field has a single setting. When several options name the same path, the last one wins.

### Field paths

Options identify a field by its path: the names of the fields leading to it, joined with dots. The names are the same
keys you see in the log. `DefaultNameExtractor` picks, in this order: the name in the `redact` tag, the name in the
`json` tag, the Go field name. Map keys are used as they are.

```go
type Customer struct {
	Name    string                    // "Name"
	Address Address `json:"address"`  // "address.Street", "address.city"
	Extra   map[string]any            // "Extra.<key>", "Extra.<key>.<nested key>"
}
```

### How values are walked

- **Nested structs and maps** are walked field by field, at any depth. Options apply to the fields inside them, not to
  the struct or map itself: allowing or hiding `"Address"` does nothing, use `"Address.Street"`.
- **Everything else** is a single value: scalars, slices, pointers (including pointers to structs) and `time.Time`. It
  is redacted as a whole, or logged in full when allowed. Allowing a `*Address` or a `[]Address` field therefore logs
  all of its content unredacted.
- **Unexported fields** and fields tagged `redact:"-"` are skipped.
- Only the name part of the `redact` tag is used. The `allow` and `redactor` parts are ignored here: use the options.
- Map keys must be strings.

## Redactors

A redactor turns a value into the string that goes to the log.

| Tag name    | Constant                              | Output                                                               | Accepts                        |
|-------------|---------------------------------------|----------------------------------------------------------------------|--------------------------------|
| (empty)     | `redactors.Default`                   | `[redacted]`                                                         | anything                       |
| `pan64`     | `redactors.PAN64`                     | First 6 and last 4 characters, `*` in between: `123456******4321`    | `string`, `*string`            |
| `bin`       | `redactors.BIN`                       | First 6 characters: `123456`                                         | `string`, `*string`            |
| `star`, `*` | `redactors.Star`, `redactors.Asterisk` | One `*` per byte: `******`                                          | `string`, `*string`            |
| `len`       | `redactors.Len`                       | The length: `[len:16]`                                               | `string`, `*string`, slices    |

Edge cases:

- A value of a type the redactor does not accept is logged as `[redactor not compatible]`.
- `pan64` needs at least 12 characters and `bin` at least 7. Shorter values are logged as `[redacted]`.
- A nil value is logged as an empty string by `pan64`, `bin` and `star`, and as `[len:0]` by `len`.
- An unknown redactor name falls back to the default redactor.

A redactor can also be applied to a single value:

```go
logger.Info("payment", zap.String("pan", zapredactor.RedactValue(pan, redactors.PAN64)))
```

### Changing the default text

```go
redactors.DefaultRedactedString = "***"
```

Set it once at startup, before anything is logged.

### Custom redactors

A redactor is a function that receives the value and returns a string:

```go
func MyRedactor(_ interface{}) string {
	return "[mine]"
}
```

Register it under a name. Registration is not safe for concurrent use, so do it at startup, usually in `init`.
Registering an existing name replaces that redactor.

```go
func init() {
	zapredactor.RegisterRedactor("my_redactor", MyRedactor)
}
```

Then use the name in a tag or in an option:

```go
type Card struct {
	CVV string `redact:",,my_redactor"`
}

redactreflection.Redact("card", &card, redactreflection.WithRedactor("CVV", "my_redactor"))
```

## Examples

Each directory under [`examples`](examples) is a runnable program:

| Example                      | Shows                                                         |
|------------------------------|---------------------------------------------------------------|
| `examples/card`              | Code generation with nested structs and redactors.            |
| `examples/demo`              | Code generation with every supported field type.              |
| `examples/reflectiondemo`    | Reflection with a reused configuration.                       |
| `examples/reflectiondemo2`   | Reflection with options passed on each call.                  |
| `examples/reflectiondemo3`   | Reflection over a map.                                        |

```bash
go run ./examples/card
```

## Development

```bash
go test ./...   # run the tests
make lint       # run golangci-lint
make generate   # regenerate the template code and the examples
```

## License

MIT license
