package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

// Version returns a command that prints the version number of the papyrus CLI.
func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of papyrus",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("papyrus %s\n", version)
		},
	}
}
