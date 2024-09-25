package client

import "encoding/json"

const (
	VirtualNetworkStateNew   = "NEW"
	VirtualNetworkStateReady = "READY"
)

type VirtualNetworkService interface {
	Page(request PageRequest) (*Page, *[]VirtualNetwork, error)
	Get(id string) (*VirtualNetwork, error)
	Create(vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error)
	Delete(id string) (*VirtualNetworkTask, error)
	Update(id string, vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error)
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
	Id        string `json:"id"`
	Name      string `json:"name"`
	Group     string `json:"group,omitempty"`
	GroupName string `json:"groupName,omitempty"`
	Type      string `json:"type"`
	Managed   bool   `json:"managed"`
	State     string `json:"state"`
}

type VirtualNetworkUpdate struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Group string `json:"group,omitempty"`
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
