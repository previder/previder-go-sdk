package cmd

import (
	"fmt"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"os"
)

var accessToken string
var customerId string
var previderClient *client.BaseClient

var rootCmd = &cobra.Command{
	Use:   "previder-cli",
	Short: "Previder CLI is the command line client for the Previder Portal",
	Long: `Previder CLI is the command line client for the Previder Portal
                more information can be found at https://portal.previder.com/`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		previderClient, err = client.New(&client.ClientOptions{Token: accessToken, CustomerId: customerId})
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&accessToken, "previder-access-token", "a", "", "The Previder access token")
	rootCmd.PersistentFlags().StringVarP(&customerId, "sub-customer", "c", "", "An optional subcustomer id")
	err := rootCmd.MarkPersistentFlagRequired("previder-access-token")
	if err != nil {
		return
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
