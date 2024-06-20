package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Previder CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Previder CLI v1.0")
		},
	}

	rootCmd.AddCommand(versionCmd)
}
