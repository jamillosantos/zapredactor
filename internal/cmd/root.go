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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zapredactor",
	Short: "A brief description of your application",
	Long:  `.`,
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
				if f.IsArray {
					p.IncludeZapArray = true
				}
			}
		}

		buff := bytes.NewBuffer(nil)
		templates.WriteRedactor(buff, p)
		formatted, err := format.Source(buff.Bytes())
		if err != nil {
			os.Stdout.Write(buff.Bytes())
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().StringVarP(&destination, "destination", "d", destination, `Destination file ("-" for stdout)`)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zapredactor.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
