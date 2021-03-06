// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// goprompt is a simple tool for prompting the user for input.
//
// Usage
//
// Use `goprompt help` to see full help documentation.
//
//	goprompt [--secret] [--json] [--question QUESTION]
//
// Examples
//
//	# show the
//	goprompt --secret --question "MFA SERIAL"
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/spatialcurrent/goprompt/pkg/prompt"
)

const (
	flagQuestion = "question"
	flagSecret   = "secret"
	flagJSON     = "json"
	flagLoop     = "loop"
)

func initFlags(flag *pflag.FlagSet) {
	flag.StringP(flagQuestion, "q", "", "the question for the prompt")
	flag.BoolP(flagSecret, "s", false, "use secret prompt")
	flag.BoolP(flagJSON, "j", false, "validate input as JSON")
	flag.BoolP(flagLoop, "l", false, "loop until non-blank input")
}

func main() {

	rootCommand := &cobra.Command{
		Use:                   "goprompt [--question QUESTION] [--secret] [--json] [--loop]",
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    false,
		Short:                 `goprompt is a simple tool for prompting the user for input.`,
		SilenceUsage:          true,
		SilenceErrors:         true,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()

			if errorBind := v.BindPFlags(cmd.Flags()); errorBind != nil {
				return errorBind
			}

			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			if len(args) != 0 {
				return errors.New("no positional arguments expected")
			}

			question := v.GetString(flagQuestion)
			loop := v.GetBool(flagLoop)

			if v.GetBool(flagSecret) {
				if v.GetBool(flagJSON) {
					value, err := prompt.SecretJSON(question, false, loop)
					if err != nil {
						return fmt.Errorf("error prompting for secret JSON: %w", err)
					}
					// print value to stdout
					fmt.Println(value)
					return nil
				}
				value, err := prompt.SecretString(question, false, loop)
				if err != nil {
					return fmt.Errorf("error prompting for secret string: %w", err)
				}
				// print value to stdout
				fmt.Println(value)
				return nil
			}

			if v.GetBool(flagJSON) {
				value, err := prompt.JSON(question, false, loop)
				if err != nil {
					return fmt.Errorf("error prompting for JSON: %w", err)
				}
				// print value to stdout
				fmt.Println(value)
				return nil
			}

			value, err := prompt.String(question, false, loop)
			if err != nil {
				return fmt.Errorf("error prompting for string: %w", err)
			}
			// print value to stdout
			fmt.Println(value)
			return nil
		},
	}
	initFlags(rootCommand.Flags())

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "goprompt: "+err.Error())
		fmt.Fprintln(os.Stderr, "Try goprompt --help for more information.")
		os.Exit(1)
	}
}
