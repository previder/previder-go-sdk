package client

import (
	"errors"
	"strings"
)

const virtualNetworkBasePath = "virtualnetwork/"

type VirtualNetworkService interface {
	List() (*[]VirtualNetwork, error)
	Get(virtualNetworkName string) (*VirtualNetwork, error)
	Create(network VirtualNetworkCreate) (*Task, error)
	Update(id string, network VirtualNetworkUpdate) (*VirtualNetwork, error)
	Delete(id string) (*Task, error)
}

type VirtualNetworkServiceOp struct {
	client *BaseClient
}

type VirtualNetwork struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	VlanId    int    `json:"vlanId,omitempty"`
	PublicNet bool   `json:"publicNet,omitempty"`
	Type      string `json:"type,omitempty"`
	AddressPool AddressPool `json:"addressPool,omitempty"`
}

type AddressPool struct {
	Editable  bool	 `json:"editable,omitempty"`
	End  	  string `json:"end,omitempty"`
	Gateway	  string `json:"gateway,omitempty"`
	Mask	  string `json:"mask,omitempty"`
	NameServers	 []string `json:"nameServers,omitempty"`
	Start	  string `json:"start,omitempty"`
	Warning	  string `json:"warning,omitempty"`
}

type VirtualNetworkCreate struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type VirtualNetworkUpdate struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	VlanId    int    `json:"vlanId,omitempty"`
	PublicNet bool   `json:"publicNet,omitempty"`
	Type      string `json:"type,omitempty"`
	AddressPool AddressPool `json:"addressPool,omitempty"`
}

func (c *VirtualNetworkServiceOp) List() (*[]VirtualNetwork, error) {
	virtualNetwork := new([]VirtualNetwork)
	err := c.client.Get(virtualNetworkBasePath, virtualNetwork)
	return virtualNetwork, err
}

func (c *VirtualNetworkServiceOp) Get(virtualNetworkName string) (*VirtualNetwork, error) {
	virtualNetworks, err := c.List()

	if err != nil {
		return nil, err
	}

	for _, n := range *virtualNetworks {
		if strings.EqualFold(n.Name, virtualNetworkName) {
			return &n, nil
		}
	}
	return nil, errors.New("VirtualNetwork not found")
}

func (c *VirtualNetworkServiceOp) Create(network VirtualNetworkCreate) (*Task, error) {
	task := new(Task)
	err := c.client.Post(virtualNetworkBasePath, network, task)
	return task, err
}

func (c *VirtualNetworkServiceOp) Update(id string, network VirtualNetworkUpdate) (*VirtualNetwork, error) {
	virtualNetwork := new(VirtualNetwork)
	err := c.client.Put(virtualNetworkBasePath+id, network, nil)
	return virtualNetwork, err
}

func (c *VirtualNetworkServiceOp) Delete(id string) (*Task, error) {
	task := new(Task)
	err := c.client.Delete(virtualNetworkBasePath+id, task)
	return task, err
}
