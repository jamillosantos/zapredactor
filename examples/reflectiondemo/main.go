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
