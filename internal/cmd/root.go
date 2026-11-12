package cmd

import (
	"bytes"
	"go/format"
	"os"

	mah "github.com/jamillosantos/go-my-ast-hurts"
	"github.com/spf13/cobra"

	"github.com/jamillosantos/zapredactor/internal/domain"
	"github.com/jamillosantos/zapredactor/internal/templates"
)

var (
	destination = "-"
)

// rootCmd is the generator: it reads a package and writes the Redact methods for its tagged structs.
var rootCmd = &cobra.Command{
	Use:   "zapredactor [package]",
	Short: "Generates the Redact method for the structs that have redact tags",
	Run: func(cmd *cobra.Command, args []string) {
		env, err := mah.NewEnvironment()
		if err != nil {
			panic(err)
		}

		packageName := "."
		if len(args) > 0 {
			packageName = args[0]
		}

		pkg, err := env.Parse(packageName)
		if err != nil {
			panic(err)
		}

		p := domain.Package{
			Name:    pkg.Name,
			Structs: make([]domain.RedactedStruct, 0),
		}
		for _, s := range pkg.Structs {
			if !hasRedactorTags(s) {
				continue
			}
			s := domain.ToRedactedStruct(s)
			p.Structs = append(p.Structs, s)

			for _, f := range s.Fields {
				if templates.UsesZap(f) {
					p.IncludeZap = true
				}
			}
		}

		buff := bytes.NewBuffer(nil)
		templates.WriteRedactor(buff, p)
		formatted, err := format.Source(buff.Bytes())
		if err != nil {
			_, _ = os.Stdout.Write(buff.Bytes())
			panic(err)
		}
		dst := os.Stdout
		if destination != "-" {
			dst, err = os.Create(destination)
			if err != nil {
				panic(err)
			}
		}
		_, err = dst.Write(formatted)
		if err != nil {
			panic(err)
		}
	},
}

func hasRedactorTags(s *mah.Struct) bool {
	for _, f := range s.Fields {
		t := f.Tag.TagParamByName("redact")
		if t != nil {
			return true
		}
	}
	return false
}

// Execute runs the generator. It is called by main.main().
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&destination, "destination", "d", destination, `Destination file ("-" for stdout)`)
}
