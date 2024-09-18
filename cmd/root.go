package cmd

import (
	"fmt"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	accessToken    string
	customerId     string
	baseUri        string
	previderClient *client.BaseClient
)

var rootCmd = &cobra.Command{
	Use:   "previder-cli",
	Short: "Previder CLI is the command line client for the Previder Portal",
	Long: `Previder CLI is the command line client for the Previder Portal.
More information can be found at https://portal.previder.com/api-docs.html or at https://previder.com`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if accessToken == "" {
			if os.Getenv("PREVIDER_TOKEN") != "" {
				accessToken = os.Getenv("PREVIDER_TOKEN")
			} else {
				log.Fatal("No token found")
			}
		}

		previderClient, err = client.New(&client.ClientOptions{Token: accessToken, CustomerId: customerId, BaseUrl: baseUri})
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&accessToken, "previder-access-token", "t", "", "The Previder access token")
	rootCmd.PersistentFlags().StringVarP(&customerId, "sub-customer", "c", "", "An optional subcustomer id")
	rootCmd.PersistentFlags().StringVarP(&baseUri, "uri", "u", "https://portal.previder.nl/api/", "Optional different URI")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
