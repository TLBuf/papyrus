// Package command defines the papyrus command line utility.
package main

import (
	"github.com/TLBuf/papyrus/cmd/papyrus/cmd"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "papyrus",
		Short: "A CLI for working with Papyrus",
		Long:  `papyrus is a command line utility for working with Bethesda's Papyrus scripting language.`,
	}

	root.AddCommand(cmd.Version())
	root.AddCommand(cmd.Format())

	root.Execute()
}
