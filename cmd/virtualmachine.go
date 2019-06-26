package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/previder/previder-go-sdk/client"
	"os"
	"strconv"
	"strings"
)

func init() {
	var virtualMachineCmd = &cobra.Command{
		Use:   "virtualmachine",
		Short: "Virtual machine commands",
	}
	rootCmd.AddCommand(virtualMachineCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of virtual machines",
		Args:  cobra.NoArgs,
		RunE:  list,
	}
	virtualMachineCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:   "get",
		Short: "Get a virtual machine",
		Args:  cobra.ExactArgs(1),
		RunE:  get,
	}
	virtualMachineCmd.AddCommand(cmdGet)

	var cmdConsole = &cobra.Command{
		Use:   "console",
		Short: "Open the console of a virtual machine",
		Args:  cobra.ExactArgs(1),
		RunE:  openConsole,
	}
	virtualMachineCmd.AddCommand(cmdConsole)

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a virtual machine",
		Args:  cobra.ExactArgs(1),
		RunE:  create,
	}

	cmdCreate.Flags().IntP("cpu-cores", "c", 1, "Number of CPU cores")
	cmdCreate.Flags().Uint64P("memory", "m", 1024, "Memory size in MB")
	cmdCreate.Flags().StringP("group", "g", "", "Group")
	cmdCreate.Flags().StringArrayP("tag", "t", []string{}, "Tag")
	cmdCreate.Flags().StringP("compute-cluster", "", "express", "Compute cluster")
	cmdCreate.Flags().StringArrayP("disk", "d", []string{}, "Disk size in MB (multiple arguments allowed)")
	cmdCreate.Flags().StringArrayP("network-interface", "n", []string{}, "Network interface Network:[Connected] (multiple arguments allowed)")
	cmdCreate.Flags().BoolP("termination-protection", "", false, "Termination protection")
	cmdCreate.Flags().StringP("template", "", "ubuntu1604lts", "Template")
	cmdCreate.Flags().StringP("source-virtualmachine", "", "", "Source virtual machine (clone)")
	cmdCreate.Flags().StringP("user-data", "u", "", "User data")
	cmdCreate.Flags().StringP("guest-id", "", "", "Guest ID")
	cmdCreate.Flags().StringP("provisioning-type", "p", "", "Provisioning type")
	virtualMachineCmd.AddCommand(cmdCreate)

	var cmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete a virtual machine",
		Args:  cobra.ExactArgs(1),
		RunE:  delete,
	}
	virtualMachineCmd.AddCommand(cmdDelete)

}

func list(cmd *cobra.Command, args []string) error {
	_, content, err := previderClient.VirtualMachine.Page()
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "CPU cores", "Memory"})
	for _, virtualMachine := range *content {
		table.Append([]string{
			virtualMachine.Name,
			strconv.Itoa(virtualMachine.CpuCores),
			ToHumanReadable(uint64(virtualMachine.Memory * 1048576)),
		})
	}
	table.Render()
	return nil
}

func get(cmd *cobra.Command, args []string) error {
	content, err := previderClient.VirtualMachine.Get(args[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)
	return nil
}

func create(cmd *cobra.Command, args []string) error {
	var err error
	var vm client.VirtualMachineCreate
	vm.Name = args[0]

	vm.CpuCores, err = cmd.Flags().GetInt("cpu-cores")
	if err != nil {
		return err
	}

	vm.Memory, err = cmd.Flags().GetUint64("memory")
	if err != nil {
		return err
	}

	vm.ComputeCluster, err = cmd.Flags().GetString("compute-cluster")
	if err != nil {
		return err
	}

	vm.Template, err = cmd.Flags().GetString("template")
	if err != nil {
		return err
	}

	vm.Tags, err = cmd.Flags().GetStringArray("tag")
	vm.Tags = []string{}
	fmt.Println(vm.Tags)
	if err != nil {
		return err
	}

	// Parse disks arguments
	disks, err := cmd.Flags().GetStringArray("disk")
	if err != nil {
		return err
	}
	for id, disk := range disks {
		size, err := FromHumanReadable(disk)
		if err != nil {
			return nil
		}
		vm.Disks = append(vm.Disks, client.Disk{
			Id:   &id,
			Size: size / 1048576,
		})
	}

	// Parse network interface arguments
	nics, err := cmd.Flags().GetStringArray("network-interface")
	if err != nil {
		return err
	}
	for _, nic := range nics {
		connected := true
		var network string
		p := strings.Split(nic, ":")
		if len(p) > 2 || len(p) == 0 {
			return fmt.Errorf("invalid nic %s", nic)
		}
		if len(p) > 0 {
			network = p[0]
		}
		if len(p) > 1 {
			connected = strings.ToLower(p[1]) == "connected"
		}

		vm.NetworkInterfaces = append(vm.NetworkInterfaces, client.NetworkInterface{
			Network:   network,
			Connected: connected,
		})
		fmt.Println(nic)
	}

	task, err := previderClient.VirtualMachine.Create(&vm)
	if err != nil {
		return err
	}

	finishedTask, err := previderClient.Task.WaitFor(task.Id, client.DefaultTimeout)
	if err != nil {
		return err
	}

	fmt.Println(finishedTask)
	return nil
}

func delete(cmd *cobra.Command, args []string) error {
	_, err := previderClient.VirtualMachine.Delete(args[0])
	if err != nil {
		return err
	}
	return nil
}

func openConsole(cmd *cobra.Command, args []string) error {
	res, err := previderClient.VirtualMachine.OpenConsole(args[0])
	if err != nil {
		return err
	}

	err = browser.OpenURL(res.ConsoleUrl)
	if err != nil {
		fmt.Print("Unable to open a browser. Use the following URL to open the console for this virtual machine: ")
		fmt.Println(res.ConsoleUrl)
	}
	return nil
}
