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
