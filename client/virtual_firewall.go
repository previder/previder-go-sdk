package client

import (
	"encoding/json"
	"net"
)

type VirtualFirewallService interface {
	Page(request PageRequest) (*Page, *[]VirtualFirewall, error)
	Get(id string) (*VirtualFirewallExt, error)
	Create(create VirtualFirewallCreate) (*Reference, error)
	Delete(id string) error
	Update(id string, update VirtualFirewallUpdate) error
	PageNatRules(firewallId string, request PageRequest) (*Page, *[]VirtualFirewallNatRule, error)
	CreateNatRule(firewallId string, create VirtualFirewallNatRuleCreate) (*Reference, error)
	UpdateNatRule(firewallId string, id string, create VirtualFirewallNatRuleCreate) error
	DeleteNatRule(firewallId string, id string) error
}

type VirtualFirewallServiceImpl struct {
	client *PreviderClient
}

type VirtualFirewall struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Group       string   `json:"group,omitempty"`
	GroupName   string   `json:"groupName,omitempty"`
	TypeLabel   string   `json:"typeLabel"`
	TypeName    string   `json:"typeName"`
	Network     string   `json:"network"`
	NetworkName string   `json:"networkName"`
	WanAddress  []string `json:"wanAddress,omitempty"`
	LanAddress  string   `json:"lanAddress"`
	State       string   `json:"state"`
}

type VirtualFirewallExt struct {
	VirtualFirewall
	TerminationProtected bool     `json:"terminationProtected"`
	Audit                Audit    `json:"audit"`
	DhcpEnabled          bool     `json:"dhcpEnabled"`
	DhcpRangeStart       string   `json:"dhcpRangeStart"`
	DhcpRangeEnd         string   `json:"dhcpRangeEnd"`
	LocalDomainName      string   `json:"localDomainName"`
	DnsEnabled           bool     `json:"dnsEnabled"`
	Nameservers          []string `json:"nameservers"`
	IcmpWanEnabled       bool     `json:"icmpWanEnabled"`
	IcmpLanEnabled       bool     `json:"icmpLanEnabled"`
}

type VirtualFirewallUpdate struct {
	Name                 string   `json:"name"`
	Group                string   `json:"group,omitempty"`
	Network              string   `json:"network"`
	LanAddress           string   `json:"lanAddress"`
	DhcpEnabled          bool     `json:"dhcpEnabled"`
	DhcpRangeStart       net.IP   `json:"dhcpRangeStart"`
	DhcpRangeEnd         net.IP   `json:"dhcpRangeEnd"`
	LocalDomainName      string   `json:"localDomainName"`
	DnsEnabled           bool     `json:"dnsEnabled"`
	Nameservers          []net.IP `json:"nameservers"`
	TerminationProtected bool     `json:"terminationProtected"`
	IcmpWanEnabled       bool     `json:"icmpWanEnabled"`
	IcmpLanEnabled       bool     `json:"icmpLanEnabled"`
}

type VirtualFirewallCreate struct {
	VirtualFirewallUpdate
	Type string `json:"type"`
}

type VirtualFirewallNatRule struct {
	Id             string `json:"id"`
	Description    string `json:"description"`
	Active         bool   `json:"active"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	Source         string `json:"source"`
	Destination    string `json:"destination"`
	NatDestination string `json:"natDestination"`
	NatPort        int    `json:"natPort"`
	WanInterface   string `json:"wanInterface"`
}

type VirtualFirewallNatRuleCreate struct {
	Description    string `json:"description"`
	Active         bool   `json:"active"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	Source         string `json:"source"`
	NatDestination string `json:"natDestination"`
	NatPort        int    `json:"natPort"`
}

func (c *VirtualFirewallServiceImpl) Page(request PageRequest) (*Page, *[]VirtualFirewall, error) {
	page := new(Page)
	err := c.client.Get(iaasBasePath+"/virtualfirewall", page, &request)
	if err != nil {
		return nil, nil, err
	}

	response := new([]VirtualFirewall)
	if err := json.Unmarshal(page.Content, &response); err != nil {
		return nil, nil, err
	}

	return page, response, err
}

func (c *VirtualFirewallServiceImpl) Get(id string) (*VirtualFirewallExt, error) {
	response := new(VirtualFirewallExt)
	err := c.client.Get(iaasBasePath+"/virtualfirewall/"+id, response, nil)
	return response, err
}

func (c *VirtualFirewallServiceImpl) Create(create VirtualFirewallCreate) (*Reference, error) {
	response := new(Reference)
	err := c.client.Post(iaasBasePath+"/virtualfirewall", create, response)
	return response, err
}

func (c *VirtualFirewallServiceImpl) Update(id string, update VirtualFirewallUpdate) error {
	err := c.client.Put(iaasBasePath+"/virtualfirewall/"+id, update, nil)
	return err
}

func (c *VirtualFirewallServiceImpl) Delete(id string) error {
	err := c.client.Delete(iaasBasePath+"/virtualfirewall/"+id, nil)
	return err
}

// NAT Rules
func (c *VirtualFirewallServiceImpl) PageNatRules(firewallId string, request PageRequest) (*Page, *[]VirtualFirewallNatRule, error) {
	page := new(Page)
	err := c.client.Get(iaasBasePath+"/virtualfirewall/"+firewallId+"/natrules", page, &request)
	if err != nil {
		return nil, nil, err
	}

	rules := new([]VirtualFirewallNatRule)
	if err := json.Unmarshal(page.Content, &rules); err != nil {
		return nil, nil, err
	}

	return page, rules, err
}

func (c *VirtualFirewallServiceImpl) CreateNatRule(firewallId string, create VirtualFirewallNatRuleCreate) (*Reference, error) {
	response := new(Reference)
	err := c.client.Post(iaasBasePath+"/virtualfirewall/"+firewallId+"/natrules", create, &response)
	return response, err
}

func (c *VirtualFirewallServiceImpl) UpdateNatRule(firewallId string, id string, create VirtualFirewallNatRuleCreate) error {
	err := c.client.Put(iaasBasePath+"/virtualfirewall/"+firewallId+"/natrules/"+id, create, nil)
	return err
}

func (c *VirtualFirewallServiceImpl) DeleteNatRule(firewallId string, id string) error {
	err := c.client.Delete(iaasBasePath+"/virtualfirewall/"+firewallId+"/natrules/"+id, nil)
	return err
}
