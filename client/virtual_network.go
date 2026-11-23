package client

import "encoding/json"

const (
	VirtualNetworkStateNew   = "NEW"
	VirtualNetworkStateReady = "READY"

	AddressPoolTypeIPv4 = "IPV4"
	AddressPoolTypeIPv6 = "IPV6"
)

type VirtualNetworkService interface {
	Page(request PageRequest) (*Page, *[]VirtualNetwork, error)
	Get(id string) (*VirtualNetwork, error)
	Create(vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error)
	Delete(id string) (*VirtualNetworkTask, error)
	Update(id string, vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error)

	ListAddressPools(networkId string) (*[]AddressPool, error)
	GetAddressPool(networkId string, id string) (*AddressPool, error)
	CreateAddressPool(networkId string, pool *AddressPoolUpdate) (*AddressPool, error)
	UpdateAddressPool(networkId string, id string, pool *AddressPoolUpdate) (*AddressPool, error)
	DeleteAddressPool(networkId string, id string) error
}

type VirtualNetworkServiceImpl struct {
	client *PreviderClient
}

type VirtualNetworkTask struct {
	Task
	VirtualNetwork     string
	VirtualNetworkName string
}

type VirtualNetwork struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group,omitempty"`
	GroupName  string `json:"groupName,omitempty"`
	Type       string `json:"type"`
	Managed    bool   `json:"managed"`
	State      string `json:"state"`
	Region     string `json:"region,omitempty"`
	RegionName string `json:"regionName,omitempty"`
}

type VirtualNetworkUpdate struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Group  string `json:"group,omitempty"`
	Region string `json:"region,omitempty"`
}

type AddressPool struct {
	Id             string   `json:"id,omitempty"`
	VirtualNetwork string   `json:"virtualNetwork,omitempty"`
	Start          string   `json:"start"`
	End            string   `json:"end"`
	Mask           string   `json:"mask"`
	Gateway        string   `json:"gateway,omitempty"`
	Type           string   `json:"type"`
	NameServers    []string `json:"nameServers,omitempty"`
	Region         string   `json:"region,omitempty"`
	RegionName     string   `json:"regionName,omitempty"`
}

type AddressPoolUpdate struct {
	Start       string   `json:"start"`
	End         string   `json:"end"`
	Mask        string   `json:"mask"`
	Gateway     string   `json:"gateway,omitempty"`
	Type        string   `json:"type"`
	NameServers []string `json:"nameServers,omitempty"`
}

func (c *VirtualNetworkServiceImpl) Page(request PageRequest) (*Page, *[]VirtualNetwork, error) {
	page := new(Page)
	err := c.client.Get(iaasBasePath+"virtualnetwork", page, &request)
	if err != nil {
		return nil, nil, err
	}

	virtualNetworks := new([]VirtualNetwork)
	if err := json.Unmarshal(page.Content, &virtualNetworks); err != nil {
		return nil, nil, err
	}

	return page, virtualNetworks, err
}

func (c *VirtualNetworkServiceImpl) Get(id string) (*VirtualNetwork, error) {
	virtualNetwork := new(VirtualNetwork)
	err := c.client.Get(iaasBasePath+"virtualnetwork/"+id, virtualNetwork, nil)
	return virtualNetwork, err
}

func (c *VirtualNetworkServiceImpl) Create(vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error) {
	task := new(VirtualNetworkTask)
	err := c.client.Post(iaasBasePath+"virtualnetwork", vn, task)
	return task, err
}

func (c *VirtualNetworkServiceImpl) Update(id string, vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error) {
	task := new(VirtualNetworkTask)
	err := c.client.Put(iaasBasePath+"virtualnetwork/"+id, vn, task)
	return task, err
}

func (c *VirtualNetworkServiceImpl) Delete(id string) (*VirtualNetworkTask, error) {
	task := new(VirtualNetworkTask)
	err := c.client.Delete(iaasBasePath+"virtualnetwork/"+id, task)
	return task, err
}

func (c *VirtualNetworkServiceImpl) ListAddressPools(networkId string) (*[]AddressPool, error) {
	pools := new([]AddressPool)
	err := c.client.Get(iaasBasePath+"virtualnetwork/"+networkId+"/addresspool", pools, nil)
	return pools, err
}

func (c *VirtualNetworkServiceImpl) GetAddressPool(networkId string, id string) (*AddressPool, error) {
	pool := new(AddressPool)
	err := c.client.Get(iaasBasePath+"virtualnetwork/"+networkId+"/addresspool/"+id, pool, nil)
	return pool, err
}

func (c *VirtualNetworkServiceImpl) CreateAddressPool(networkId string, pool *AddressPoolUpdate) (*AddressPool, error) {
	created := new(AddressPool)
	err := c.client.Post(iaasBasePath+"virtualnetwork/"+networkId+"/addresspool", pool, created)
	return created, err
}

func (c *VirtualNetworkServiceImpl) UpdateAddressPool(networkId string, id string, pool *AddressPoolUpdate) (*AddressPool, error) {
	updated := new(AddressPool)
	err := c.client.Put(iaasBasePath+"virtualnetwork/"+networkId+"/addresspool/"+id, pool, updated)
	return updated, err
}

func (c *VirtualNetworkServiceImpl) DeleteAddressPool(networkId string, id string) error {
	return c.client.Delete(iaasBasePath+"virtualnetwork/"+networkId+"/addresspool/"+id, nil)
}
