package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of papyrus",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("papyrus %s\n", version)
		},
	}
}
