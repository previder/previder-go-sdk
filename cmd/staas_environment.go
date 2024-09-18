package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	var staasEnvironmentCmd = &cobra.Command{
		Use:   "staasenvironment",
		Short: "STaaS environment commands",
	}
	rootCmd.AddCommand(staasEnvironmentCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of staas environments",
		Args:  cobra.NoArgs,
		RunE:  listSTaaSEnvironments,
	}
	staasEnvironmentCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:   "get",
		Short: "Get a staas environment",
		Args:  cobra.ExactArgs(1),
		RunE:  getSTaaSEnvironment,
	}
	staasEnvironmentCmd.AddCommand(cmdGet)

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a staas environment",
		RunE:  createSTaaSEnvironment,
	}

	cmdCreate.Flags().String("name", "name", "Name of the cluster")
	cmdCreate.MarkFlagRequired("name")
	cmdCreate.Flags().String("type", "nfs", "Type (NFS or ISCSI)")
	cmdCreate.Flags().String("cluster", "", "ID of the STaaS cluster")
	cmdCreate.MarkFlagRequired("cluster")

	/*cmdCreate.Flags().String("volume-name", "my-volume", "Name of the initial volume")
	cmdCreate.Flags().Int("volume-sizemb", 10240, "Size of the initial volume")
	cmdCreate.Flags().String("volume-type", "express", "Type of the initial volume. Refer to STaaS cluster information for available storage types")
	cmdCreate.Flags().String("volume-allowedipro", "192.168.0.0/24", "Comma-seperated list of cidrs that are allowed to read the volume")
	cmdCreate.Flags().String("volume-allowediprw", "192.168.0.0/24", "Comma-seperated list of cidrs that are allowed to write to the volume")

	cmdCreate.Flags().String("network-id", "my-volume", "ID of the virtual network. Target network must be type Cloud VLAN")
	cmdCreate.Flags().String("network-cidr", "192.168.0.50/24", "IP address in CIDR format.")*/

	staasEnvironmentCmd.AddCommand(cmdCreate)

}

func listSTaaSEnvironments(cmd *cobra.Command, args []string) error {
	var page client.PageRequest
	page.Size = 100
	page.Page = 0
	page.Sort = "+name"
	page.Query = ""
	_, content, err := previderClient.STaaSEnvironment.Page(page)
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "State", "Cluster", "ClusterId", "Type"})
	for _, sTaaSEnvironment := range *content {
		table.Append([]string{
			sTaaSEnvironment.Id,
			sTaaSEnvironment.Name,
			sTaaSEnvironment.State,
			sTaaSEnvironment.Cluster,
			sTaaSEnvironment.ClusterId,
			sTaaSEnvironment.Type,
		})
	}
	table.Render()
	return nil
}

func getSTaaSEnvironment(cmd *cobra.Command, args []string) error {
	content, err := previderClient.STaaSEnvironment.Get(args[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", content)
	return nil
}

func createSTaaSEnvironment(cmd *cobra.Command, args []string) error {
	var create client.STaaSEnvironmentCreate
	var err error
	create.Name, err = cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	create.Type, err = cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	create.Cluster, err = cmd.Flags().GetString("cluster")
	if err != nil {
		return err
	}

	previderClient.STaaSEnvironment.Create(create)
	return nil
}
