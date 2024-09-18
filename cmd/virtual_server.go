package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/browser"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	var virtualServerCmd = &cobra.Command{
		Use:   "virtualserver",
		Short: "Virtual server commands",
	}
	rootCmd.AddCommand(virtualServerCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of virtual servers",
		Args:  cobra.NoArgs,
		RunE:  listVirtualServer,
	}
	virtualServerCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:   "get",
		Short: "Get a virtual server",
		Args:  cobra.ExactArgs(1),
		RunE:  getVirtualServer,
	}
	virtualServerCmd.AddCommand(cmdGet)

	var cmdConsole = &cobra.Command{
		Use:   "console",
		Short: "Open the console of a virtual server",
		Args:  cobra.ExactArgs(1),
		RunE:  openConsole,
	}
	virtualServerCmd.AddCommand(cmdConsole)

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a virtual server",
		Args:  cobra.NoArgs,
		RunE:  createVirtualServer,
		PreRun: func(cmd *cobra.Command, args []string) {
			template, err := cmd.Flags().GetString("template")
			if err != nil {
				return
			}
			sourceVirtualServer, err := cmd.Flags().GetString("sourceVirtualServer")
			if err != nil {
				return
			}
			guestId, err := cmd.Flags().GetString("guestId")
			if err != nil {
				return
			}
			if template == "" && sourceVirtualServer == "" && guestId == "" {
				log.Fatal("Either template, sourceVirtualServer or guestId is required")
			}
		},
	}

	cmdCreate.Flags().String("name", "", "Name of the virtual server")
	cmdCreate.MarkFlagRequired("name")
	cmdCreate.Flags().Int("cpu-cores", 1, "Number of CPU cores")
	cmdCreate.MarkFlagRequired("cpu-cores")
	cmdCreate.Flags().Uint64("memory", 1024, "Memory size in MB")
	cmdCreate.MarkFlagRequired("memory")
	cmdCreate.Flags().String("group", "", "Group")
	cmdCreate.Flags().StringArray("tag", []string{}, "Tag")
	cmdCreate.Flags().String("compute-cluster", "express", "Compute cluster")
	cmdCreate.Flags().StringArray("disk", []string{}, "Disk size in MB (multiple arguments allowed)")
	cmdCreate.MarkFlagRequired("disk")
	cmdCreate.Flags().StringArray("network-interface", []string{}, "Network interface Network:[Connected] (multiple arguments allowed)")
	cmdCreate.MarkFlagRequired("network-interface")
	cmdCreate.Flags().Bool("termination-protection", false, "Termination protection")
	cmdCreate.Flags().String("template", "", "Template")
	cmdCreate.Flags().String("source-virtual-server", "", "Source virtual server (clone)")
	cmdCreate.Flags().String("guest-id", "", "Guest ID")
	cmdCreate.Flags().String("user-data", "", "User data")
	cmdCreate.Flags().String("provisioning-type", "", "Provisioning type")
	virtualServerCmd.AddCommand(cmdCreate)

	var cmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete a virtual server",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteVirtualServer,
	}
	virtualServerCmd.AddCommand(cmdDelete)

}

func listVirtualServer(cmd *cobra.Command, args []string) error {
	var page client.PageRequest
	page.Size = 100
	page.Page = 0
	page.Sort = "+name"
	page.Query = ""

	_, content, err := previderClient.VirtualServer.Page(page)
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Group", "Type", "CPU cores", "Memory", "Storage", "State"})
	for _, virtualMachine := range *content {
		table.Append([]string{
			virtualMachine.Name,
			virtualMachine.Group,
			virtualMachine.ComputeCluster,
			strconv.Itoa(virtualMachine.CpuCores),
			ToHumanReadable(virtualMachine.Memory * 1048576),
			ToHumanReadable(uint64(virtualMachine.TotalDiskSize * 1048576)),
			virtualMachine.State,
		})
	}
	table.Render()
	return nil
}

func getVirtualServer(cmd *cobra.Command, args []string) error {
	content, err := previderClient.VirtualServer.Get(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", content)
	return nil
}

func createVirtualServer(cmd *cobra.Command, args []string) error {
	var err error
	var vm client.VirtualMachineCreate
	vm.Name, err = cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	vm.CpuCores, err = cmd.Flags().GetInt("cpuCores")
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

	vm.SourceVirtualMachine, err = cmd.Flags().GetString("source-virtual-server")
	if err != nil {
		return err
	}

	vm.GuestId, err = cmd.Flags().GetString("guest-id")
	if err != nil {
		return err
	}

	vm.Tags, err = cmd.Flags().GetStringArray("tag")
	vm.Tags = []string{}
	if err != nil {
		return err
	}

	// Parse disks arguments
	disks, err := cmd.Flags().GetStringArray("disk")
	if err != nil {
		return err
	}
	for _, disk := range disks {
		size, err := FromHumanReadable(disk)
		if err != nil {
			return nil
		}
		vm.Disks = append(vm.Disks, client.Disk{
			//	Id:   &id,
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
	}

	vm.UserData, err = cmd.Flags().GetString("user-data")
	if err != nil {
		return err
	}

	task, err := previderClient.VirtualServer.Create(&vm)
	if err != nil {
		return err
	}

	finishedTask, err := previderClient.Task.WaitFor(task.Id, client.DefaultTimeout)
	if err != nil {
		return err
	}

	fmt.Println(finishedTask)
	log.Println("Virtual Server create successful")

	return nil
}

func deleteVirtualServer(cmd *cobra.Command, args []string) error {
	_, err := previderClient.VirtualServer.Delete(args[0])
	if err != nil {
		return err
	}
	log.Println("Virtual Server delete successful")

	return nil
}

func openConsole(cmd *cobra.Command, args []string) error {
	res, err := previderClient.VirtualServer.OpenConsole(args[0])
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
