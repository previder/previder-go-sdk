package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Previder CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}

	rootCmd.AddCommand(versionCmd)
}

func getVersion() string {
	return fmt.Sprintf("Previder CLI %s-%s (%s)", version, commit, date)
}
