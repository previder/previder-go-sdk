package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	var kubernetesClusterCmd = &cobra.Command{
		Use:   "kubernetescluster",
		Short: "Kubernetes cluster commands",
	}
	rootCmd.AddCommand(kubernetesClusterCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of kubernetes clusters",
		Args:  cobra.NoArgs,
		RunE:  listKubernetesCluster,
	}
	kubernetesClusterCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:   "get",
		Short: "Get a kubernetes cluster",
		Args:  cobra.ExactArgs(1),
		RunE:  getKubernetesCluster,
	}
	kubernetesClusterCmd.AddCommand(cmdGet)

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a kubernetes cluster",
		RunE:  createKubernetesCluster,
	}

	cmdCreate.Flags().StringP("name", "", "cluster", "Name of the cluster")
	cmdCreate.Flags().StringP("version", "", "", "Version of the cluster (only when auto update is disabled)")
	cmdCreate.Flags().IntP("minimal-nodes", "", 1, "Minimal nodes")
	cmdCreate.Flags().IntP("maximal-nodes", "", 1, "Maximal Nodes")
	cmdCreate.Flags().BoolP("auto-update", "", true, "Automatically get the newest Kubernetes version")
	cmdCreate.Flags().BoolP("auto-scale-enabled", "", false, "Install an cluster autoscaler")
	cmdCreate.Flags().IntP("control-plane-cpu-cores", "", 2, "Number of cpu cores per node")
	cmdCreate.Flags().IntP("control-plane-memory-gb", "", 4, "Number of memory GB per node")
	cmdCreate.Flags().IntP("control-plane-storage-gb", "", 25, "Storage capacity per node (minimum 25)")
	cmdCreate.Flags().IntP("node-cpu-cores", "", 4, "Number of cpu cores per node")
	cmdCreate.Flags().IntP("node-memory-gb", "", 8, "Number of memory GB per node")
	cmdCreate.Flags().IntP("node-storage-gb", "", 30, "Storage capacity per node (minimum 25)")
	cmdCreate.Flags().StringP("compute-cluster", "", "", "Compute cluster")
	cmdCreate.Flags().StringP("cni", "", "", "CNI to install")
	cmdCreate.Flags().BoolP("high-available-control-plane", "", false, "Install 1 or 3 control planes")
	cmdCreate.Flags().StringArrayP("vips", "", []string{}, "VIPS as comma seperated list")
	cmdCreate.Flags().StringArrayP("endpoints", "", []string{}, "Endpoints as comma seperated list")
	cmdCreate.Flags().StringP("network", "", "", "Network to deploy in (name or objectID)")

	kubernetesClusterCmd.AddCommand(cmdCreate)

	var cmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete a kubernetes cluster",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteKubernetesCluster,
	}
	kubernetesClusterCmd.AddCommand(cmdDelete)

}

func listKubernetesCluster(cmd *cobra.Command, args []string) error {
	var page client.PageRequest
	page.Size = 100
	page.Page = 0
	page.Sort = "+name"
	page.Query = ""
	_, content, err := previderClient.KubernetesCluster.Page(page)
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "State", "Version"})
	for _, kubernetesCluster := range *content {
		table.Append([]string{
			kubernetesCluster.Id,
			kubernetesCluster.Name,
			kubernetesCluster.State,
			kubernetesCluster.Version,
		})
	}
	table.Render()
	return nil
}

func getKubernetesCluster(cmd *cobra.Command, args []string) error {
	content, err := previderClient.KubernetesCluster.Get(args[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", content)
	return nil
}

func createKubernetesCluster(cmd *cobra.Command, args []string) error {
	var create client.KubernetesClusterCreate
	var err error
	create.Name, err = cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	create.AutoUpdate, err = cmd.Flags().GetBool("auto-update")
	if err != nil {
		return err
	}
	create.Version, err = cmd.Flags().GetString("version")
	if err != nil {
		if !create.AutoUpdate {
			return err
		}
	}
	create.AutoScaleEnabled, err = cmd.Flags().GetBool("auto-scale-enabled")
	if err != nil {
		create.AutoScaleEnabled = false
	}
	create.MinimalNodes, err = cmd.Flags().GetInt("minimal-nodes")
	if err != nil {
		return err
	}
	create.MaximalNodes, err = cmd.Flags().GetInt("maximal-nodes")
	if err != nil {
		if create.AutoScaleEnabled {
			return err
		}
	}
	create.ControlPlaneCpuCores, err = cmd.Flags().GetInt("control-plane-cpu-cores")
	if err != nil {
		return err
	}
	create.ControlPlaneMemoryGb, err = cmd.Flags().GetInt("control-plane-memory-gb")
	if err != nil {
		return err
	}
	create.ControlPlaneStorageGb, err = cmd.Flags().GetInt("control-plane-storage-gb")
	if err != nil {
		return err
	}
	create.NodeCpuCores, err = cmd.Flags().GetInt("node-cpu-cores")
	if err != nil {
		return err
	}
	create.NodeMemoryGb, err = cmd.Flags().GetInt("node-memory-gb")
	if err != nil {
		return err
	}
	create.NodeStorageGb, err = cmd.Flags().GetInt("node-storage-gb")
	if err != nil {
		return err
	}
	create.ComputeCluster, err = cmd.Flags().GetString("compute-cluster")
	if err != nil {
		create.ComputeCluster = "express"
	}
	create.HighAvailableControlPlane, err = cmd.Flags().GetBool("high-available-control-plane")
	if err != nil {
		create.HighAvailableControlPlane = false
	}
	create.Vips, err = cmd.Flags().GetStringArray("vips")
	if err != nil {
		return err
	}
	create.Endpoints, err = cmd.Flags().GetStringArray("endpoints")
	if err != nil {
		create.Endpoints = []string{}
	}
	create.CNI, err = cmd.Flags().GetString("cni")
	if err != nil {
		return err
	}
	create.Network, err = cmd.Flags().GetString("network")
	if err != nil {
		return err
	}
	previderClient.KubernetesCluster.Create(create)
	return nil
}

func deleteKubernetesCluster(cmd *cobra.Command, args []string) error {
	err := previderClient.KubernetesCluster.Delete(args[0])
	if err != nil {
		return err
	}
	return nil
}
