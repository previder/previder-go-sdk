package client

import "encoding/json"

type VirtualNetworkService interface {
	Page() (*Page, *[]VirtualNetwork, error)
	Get(id string) (*VirtualNetwork, error)
	Create(vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error)
	Delete(id string) (*VirtualNetworkTask, error)
	Update(id string, vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error)
}

type VirtualNetworkServiceOp struct {
	client *BaseClient
}

type VirtualNetworkTask struct {
	Task
	VirtualNetwork     string
	VirtualNetworkName string
}

type VirtualNetwork struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Group   string `json:"group,omitempty"`
	Type    string `json:"type"`
	Managed bool   `json:"managed"`
	State   string `json:"state"`
}

type VirtualNetworkUpdate struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Group string `json:"group,omitempty"`
}

func (c *VirtualNetworkServiceOp) Page() (*Page, *[]VirtualNetwork, error) {
	page := new(Page)
	err := c.client.Get(iaasBasePath+"virtualnetwork", page)
	if err != nil {
		return nil, nil, err
	}

	virtualNetworks := new([]VirtualNetwork)
	if err := json.Unmarshal([]byte(page.Content), &virtualNetworks); err != nil {
		return nil, nil, err
	}

	return page, virtualNetworks, err
}

func (c *VirtualNetworkServiceOp) Get(id string) (*VirtualNetwork, error) {
	virtualNetwork := new(VirtualNetwork)
	err := c.client.Get(iaasBasePath+"virtualnetwork/"+id, virtualNetwork)
	return virtualNetwork, err
}

func (c *VirtualNetworkServiceOp) Create(vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error) {
	task := new(VirtualNetworkTask)
	err := c.client.Post(iaasBasePath+"virtualnetwork", vn, task)
	return task, err
}

func (c *VirtualNetworkServiceOp) Update(id string, vn *VirtualNetworkUpdate) (*VirtualNetworkTask, error) {
	task := new(VirtualNetworkTask)
	err := c.client.Put(iaasBasePath+"virtualnetwork/"+id, vn, task)
	return task, err
}

func (c *VirtualNetworkServiceOp) Delete(id string) (*VirtualNetworkTask, error) {
	task := new(VirtualNetworkTask)
	err := c.client.Delete(iaasBasePath+"virtualnetwork/"+id, task)
	return task, err
}
