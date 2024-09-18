package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func init() {
	var virtualFirewallCmd = &cobra.Command{
		Use:   "virtualfirewall",
		Short: "Virtual firewall commands",
	}
	rootCmd.AddCommand(virtualFirewallCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of virtual firewalls",
		Args:  cobra.NoArgs,
		RunE:  listVirtualFirewall,
	}
	virtualFirewallCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:   "get",
		Short: "Get a virtual firewall",
		Args:  cobra.ExactArgs(1),
		RunE:  getVirtualFirewall,
	}
	virtualFirewallCmd.AddCommand(cmdGet)

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a virtual firewall",
		Args:  cobra.NoArgs,
		RunE:  createVirtualFirewall,
		PreRun: func(cmd *cobra.Command, args []string) {
			dhcpEnabled, err := cmd.Flags().GetBool("dhcpEnabled")
			if err != nil {
				log.Fatal("Error parsing DHCP config")
			}
			if dhcpEnabled {
				cmd.MarkFlagRequired("dhcpRangeStart")
				cmd.MarkFlagRequired("dhcpRangeEnd")
			}
		},
	}

	cmdCreate.Flags().String("name", "", "Name of the firewall")
	cmdCreate.MarkFlagRequired("name")
	cmdCreate.Flags().String("type", "previder", "Type of the firewall")
	cmdCreate.Flags().String("group", "", "Group of the firewall")
	cmdCreate.Flags().String("network", "", "ID or name of the network")
	cmdCreate.MarkFlagRequired("network")
	cmdCreate.Flags().String("lan", "192.168.1.1/24", "LAN CIDR in the local network")
	cmdCreate.MarkFlagRequired("lan")
	cmdCreate.Flags().Bool("dhcp-enabled", true, "Enable DHCP")
	cmdCreate.Flags().String("local-domain-name", "int", "Internal network name")
	cmdCreate.Flags().IP("dhcp-range-start", net.IP("192.168.1.10"), "Start of the DHCP range (must be in the LAN range)")
	cmdCreate.Flags().IP("dhcp-range-end", net.IP("192.168.1.100"), "End of the DHCP range (must be in the LAN range and after the range start)")
	cmdCreate.Flags().Bool("dns-enabled", true, "Enable DNS")
	cmdCreate.Flags().String("nameservers", "80.65.96.50,62.165.127.222", "List of nameservers")
	cmdCreate.Flags().Bool("termination-protected", false, "Enabled termination protection")
	cmdCreate.Flags().Bool("icmp-wan-enabled", true, "Enable ICMP on WAN interface")
	cmdCreate.Flags().Bool("icmp-lan-enabled", true, "Enable ICMP on LAN interface")

	virtualFirewallCmd.AddCommand(cmdCreate)

	var cmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete a virtual firewall",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteVirtualFirewall,
	}
	virtualFirewallCmd.AddCommand(cmdDelete)

	var virtualFirewallNatRulesCmd = &cobra.Command{
		Use:   "nat",
		Args:  cobra.NoArgs,
		Short: "virtual firewall nat rules commands",
	}
	virtualFirewallCmd.AddCommand(virtualFirewallNatRulesCmd)

	var cmdNatList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of a virtual firewall NAT rules",
		Args:  cobra.ExactArgs(1),
		RunE:  listVirtualFirewallNatRules,
	}
	virtualFirewallNatRulesCmd.AddCommand(cmdNatList)

	var cmdNatCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a virtual firewall NAT rule",
		Args:  cobra.ExactArgs(1),
		RunE:  createVirtualFirewallNatRule,
	}

	cmdNatCreate.Flags().String("description", "", "Description of the NAT rule")
	cmdNatCreate.MarkFlagRequired("description")
	cmdNatCreate.Flags().String("source", "", "Source CIDR")
	cmdNatCreate.Flags().String("protocol", "TCP", "Protocol of the port, tcp or udp")
	cmdNatCreate.Flags().Int("port", 0, "External port")
	cmdNatCreate.MarkFlagRequired("port")
	cmdNatCreate.Flags().Int("nat-port", 0, "NAT port")
	cmdNatCreate.MarkFlagRequired("nat-port")
	cmdNatCreate.Flags().IP("nat-destination", nil, "NAT Destination")
	cmdNatCreate.MarkFlagRequired("nat-destination")
	cmdNatCreate.Flags().Bool("active", true, "Rule is active")

	virtualFirewallNatRulesCmd.AddCommand(cmdNatCreate)

	var cmdNatDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete a virtual firewall NAT rule",
		Args:  cobra.ExactArgs(2),
		RunE:  deleteVirtualFirewallNatRule,
	}
	virtualFirewallNatRulesCmd.AddCommand(cmdNatDelete)
}

func listVirtualFirewall(cmd *cobra.Command, args []string) error {
	var page client.PageRequest
	page.Size = 100
	page.Page = 0
	page.Sort = "+name"
	page.Query = ""

	_, content, err := previderClient.VirtualFirewall.Page(page)
	if err != nil {
		log.Fatal(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Type", "Group", "Network", "LAN address", "WAN address", "State"})
	for _, virtualFirewall := range *content {
		table.Append([]string{
			virtualFirewall.Id,
			virtualFirewall.Name,
			virtualFirewall.TypeName,
			virtualFirewall.GroupName,
			virtualFirewall.NetworkName,
			virtualFirewall.LanAddress,
			strings.Join(virtualFirewall.WanAddress, ","),
			virtualFirewall.State,
		})
	}
	table.Render()
	return nil
}

func getVirtualFirewall(cmd *cobra.Command, args []string) error {
	content, err := previderClient.VirtualFirewall.Get(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", content)
	return nil
}

func createVirtualFirewall(cmd *cobra.Command, args []string) error {
	var err error

	var virtualfirewall client.VirtualFirewallCreate
	virtualfirewall.Name, err = cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	virtualfirewall.Type, err = cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	virtualfirewall.Group, err = cmd.Flags().GetString("group")
	if err != nil {
		return err
	}
	virtualfirewall.Network, err = cmd.Flags().GetString("network")
	if err != nil {
		return err
	}
	lanAddress, err := cmd.Flags().GetString("lan")
	if err != nil {
		return err
	}
	lanAddressIp, lanAddressNet, err := net.ParseCIDR(lanAddress)
	if err != nil {
		log.Fatal("Invalid LAN CIDR")
	}
	mask, _ := lanAddressNet.Mask.Size()
	virtualfirewall.LanAddress = fmt.Sprintf("%v/%v", lanAddressIp.String(), mask)
	if lanAddressNet.IP.To4().Equal(lanAddressIp) {
		log.Fatal("LAN Address cannot be a network address")
	}

	virtualfirewall.DhcpEnabled, err = cmd.Flags().GetBool("dhcp-enabled")
	if err != nil {
		return err
	}
	if virtualfirewall.DhcpEnabled {
		virtualfirewall.DhcpRangeStart, err = cmd.Flags().GetIP("dhcp-range-start")
		if err != nil {
			return err
		}
		if !lanAddressNet.Contains(virtualfirewall.DhcpRangeStart) {
			log.Fatal("dhcpRangeStart is not in lan CIDR")
		}

		virtualfirewall.DhcpRangeEnd, err = cmd.Flags().GetIP("dhcp-range-end")
		if err != nil {
			return err
		}
		if !lanAddressNet.Contains(virtualfirewall.DhcpRangeEnd) {
			log.Fatal("dhcpRangeEnd is not in lan CIDR")
		}

		virtualfirewall.LocalDomainName, err = cmd.Flags().GetString("local-domain-name")
		if err != nil {
			return err
		}
	}

	virtualfirewall.DnsEnabled, err = cmd.Flags().GetBool("dns-enabled")
	if err != nil {
		return err
	}
	if virtualfirewall.DnsEnabled {
		nameservers, err := cmd.Flags().GetString("nameservers")
		if err != nil {
			return err
		}

		for _, nameserver := range strings.Split(nameservers, ",") {
			parsedNameserver := net.ParseIP(nameserver)
			if parsedNameserver == nil {
				log.Fatal(fmt.Sprintf("Invalid nameserver %v", nameserver))
			}
			virtualfirewall.Nameservers = append(virtualfirewall.Nameservers, parsedNameserver)
		}
	}
	virtualfirewall.IcmpLanEnabled, err = cmd.Flags().GetBool("icmp-lan-enabled")
	if err != nil {
		return err
	}
	virtualfirewall.IcmpWanEnabled, err = cmd.Flags().GetBool("icmp-wan-enabled")
	if err != nil {
		return err
	}

	err = previderClient.VirtualFirewall.Create(virtualfirewall)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Virtual Firewall create successful")
	return nil
}

func deleteVirtualFirewall(cmd *cobra.Command, args []string) error {
	err := previderClient.VirtualFirewall.Delete(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Virtual Firewall delete successful")
	return nil
}

// NAT rules
func listVirtualFirewallNatRules(cmd *cobra.Command, args []string) error {
	var page client.PageRequest
	page.Size = 100
	page.Page = 0
	page.Sort = "+description"
	page.Query = ""

	firewall, err := previderClient.VirtualFirewall.Get(args[0])
	if err != nil {
		return err
	}

	_, content, err := previderClient.VirtualFirewall.PageNatRules(firewall.Id, page)
	if err != nil {
		log.Fatal(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Source", "Port", "NAT Destination", "NAT Port", "Active"})
	for _, virtualFirewallNatRule := range *content {
		table.Append([]string{
			virtualFirewallNatRule.Id,
			virtualFirewallNatRule.Description,
			virtualFirewallNatRule.Source,
			strconv.Itoa(virtualFirewallNatRule.Port) + " / " + virtualFirewallNatRule.Protocol,
			virtualFirewallNatRule.NatDestination,
			strconv.Itoa(virtualFirewallNatRule.NatPort),
			strconv.FormatBool(virtualFirewallNatRule.Active),
		})
	}
	table.Render()
	return nil
}

func createVirtualFirewallNatRule(cmd *cobra.Command, args []string) error {
	firewall, err := previderClient.VirtualFirewall.Get(args[0])
	if err != nil {
		return err
	}

	var natrule client.VirtualFirewallNatRuleCreate

	natrule.Description, err = cmd.Flags().GetString("description")
	if err != nil {
		return err
	}

	natrule.Active, err = cmd.Flags().GetBool("active")
	if err != nil {
		return err
	}

	natrule.Source, err = cmd.Flags().GetString("source")
	if err != nil {
		return err
	}
	if len(natrule.Source) > 0 {
		_, _, err = net.ParseCIDR(natrule.Source)
		if err != nil {
			log.Fatal("Invalid source CIDR: ", err.Error())
		}
	}

	natrule.Protocol, err = cmd.Flags().GetString("protocol")
	if err != nil {
		return err
	}
	if natrule.Protocol != "TCP" && natrule.Protocol != "UDP" {
		log.Fatal("Invalid protocol, only values TCP or UDP allowed")
	}

	natrule.Port, err = cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}
	if natrule.Port < 1 || natrule.Port > 65535 {
		log.Fatal("Invalid external port")
	}
	natrule.NatPort, err = cmd.Flags().GetInt("nat-port")
	if err != nil {
		return err
	}
	if natrule.NatPort < 1 || natrule.NatPort > 65535 {
		log.Fatal("Invalid NAT port")
	}

	natDestination, err := cmd.Flags().GetIP("nat-destination")
	if err != nil {
		return err
	}
	natrule.NatDestination = natDestination.String()

	err = previderClient.VirtualFirewall.CreateNatRule(firewall.Id, natrule)
	if err != nil {
		return err
	}
	log.Println("NAT rule create successful")

	return nil
}

func deleteVirtualFirewallNatRule(cmd *cobra.Command, args []string) error {
	err := previderClient.VirtualFirewall.DeleteNatRule(args[0], args[1])

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("NAT rule delete successful")
	return nil
}
