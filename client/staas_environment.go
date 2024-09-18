package client

import "encoding/json"

type STaaSEnvironmentService interface {
	Page(request PageRequest) (*Page, *[]STaaSEnvironment, error)
	Get(id string) (*STaaSEnvironmentExt, error)
	Create(create STaaSEnvironmentCreate) error
	Delete(id string, delete STaaSEnvironmentDelete) error
	Update(id string, update STaaSEnvironmentUpdate) error
	CreateVolume(id string, create STaaSVolumeCreate) error
	UpdateVolume(id string, volumeId string, update STaaSVolumeUpdate) error
	DeleteVolume(id string, volumeId string, delete STaaSVolumeDelete) error
	CreateNetwork(id string, create STaaSNetworkCreate) error
	DeleteNetwork(id string, networkId string) error
}

type STaaSEnvironmentServiceImpl struct {
	client *BaseClient
}

type STaaSEnvironment struct {
	Id             string         `json:"id,omitempty"`
	Name           string         `json:"name"`
	State          string         `json:"state"`
	Cluster        string         `json:"cluster"`
	ClusterId      string         `json:"clusterId"`
	Type           string         `json:"type"`
	OwnerReference OwnerReference `json:"ownerReference,omitempty"`
}

type STaaSEnvironmentExt struct {
	STaaSEnvironment
	Volumes              []STaaSVolume     `json:"volumes,omitempty"`
	Networks             []STaaSNetwork    `json:"networks,omitempty"`
	Credentials          []STaaSCredential `json:"credentials,omitempty"`
	Routes               []STaaSRoute      `json:"routes,omitempty"`
	Windows              bool              `json:"windows,omitempty"`
	SynchronousClusterId string            `json:"synchronousClusterId,omitempty"`
}

type STaaSEnvironmentDelete struct {
	Force bool `json:"force"`
}

type STaaSVolume struct {
	Id                         string   `json:"id,omitempty"`
	Name                       string   `json:"name"`
	State                      string   `json:"state"`
	Type                       string   `json:"type"`
	SynchronousEnvironmentId   string   `json:"synchronousEnvironmentId,omitempty"`
	SynchronousEnvironmentName string   `json:"synchronousEnvironmentName,omitempty"`
	AllowedIpsRo               []string `json:"allowedIpsRo,omitempty"`
	AllowedIpsRw               []string `json:"allowedIpsRw,omitempty"`
	SizeMb                     int      `json:"sizeMb"`
}

type STaaSNetwork struct {
	Id          string   `json:"id,omitempty"`
	State       string   `json:"state"`
	NetworkName string   `json:"networkName"`
	NetworkId   string   `json:"networkId"`
	IpAddresses []string `json:"ipAddresses,omitempty"`
	Cidr        string   `json:"cidr,omitempty"`
}

type STaaSCredential struct {
	Id                 string `json:"id,omitempty"`
	AuthenticationType string `json:"authenticationType"`
	State              string `json:"state"`
	Initiator          string `json:"initiator"`
}

type STaaSRoute struct {
	Id          string `json:"id,omitempty"`
	Destination string `json:"destination"`
	State       string `json:"state"`
	Gateway     string `json:"gateway"`
}

type STaaSEnvironmentCreate struct {
	STaaSEnvironment
	Windows bool   `json:"windows,omitempty"`
	Type    string `json:"type"`
	Cluster string `json:"cluster"`
}

type STaaSEnvironmentUpdate struct {
	STaaSEnvironment
	Windows bool `json:"windows,omitempty"`
}

type STaaSNetworkCreate struct {
	Network string `json:"network"`
	Cidr    string `json:"cidr"`
}

type STaaSVolumeCreate struct {
	Name                       string   `json:"name"`
	Type                       string   `json:"type"`
	SynchronousEnvironmentId   string   `json:"synchronousEnvironmentId,omitempty"`
	SynchronousEnvironmentName string   `json:"synchronousEnvironmentName,omitempty"`
	AllowedIpsRo               []string `json:"allowedIpsRo,omitempty"`
	AllowedIpsRw               []string `json:"allowedIpsRw,omitempty"`
	SizeMb                     int      `json:"sizeMb"`
}

type STaaSVolumeUpdate struct {
	STaaSVolumeCreate
}

type STaaSVolumeDelete struct {
	Force bool `json:"force"`
}

func (c *STaaSEnvironmentServiceImpl) Page(request PageRequest) (*Page, *[]STaaSEnvironment, error) {
	page := new(Page)
	err := c.client.Get(staasBasePath+"/environment", page, &request)
	if err != nil {
		return nil, nil, err
	}

	environments := new([]STaaSEnvironment)
	if err := json.Unmarshal(page.Content, &environments); err != nil {
		return nil, nil, err
	}

	return page, environments, err
}

func (c *STaaSEnvironmentServiceImpl) Get(id string) (*STaaSEnvironmentExt, error) {
	environment := new(STaaSEnvironmentExt)
	err := c.client.Get(staasBasePath+"/environment/"+id, environment, nil)
	return environment, err
}

func (c *STaaSEnvironmentServiceImpl) Create(create STaaSEnvironmentCreate) error {
	err := c.client.Post(staasBasePath+"/environment", create, nil)
	return err
}

func (c *STaaSEnvironmentServiceImpl) Update(id string, update STaaSEnvironmentUpdate) error {
	err := c.client.Put(staasBasePath+"/environment/"+id, update, nil)
	return err
}

func (c *STaaSEnvironmentServiceImpl) Delete(id string, delete STaaSEnvironmentDelete) error {
	err := c.client.Delete(staasBasePath+"/environment/"+id, delete)
	return err
}

func (c *STaaSEnvironmentServiceImpl) CreateVolume(id string, create STaaSVolumeCreate) error {
	err := c.client.Post(staasBasePath+"/environment/"+id+"/volume", create, nil)
	return err
}

func (c *STaaSEnvironmentServiceImpl) UpdateVolume(id string, volumeId string, create STaaSVolumeUpdate) error {
	err := c.client.Put(staasBasePath+"/environment/"+id+"/volume/"+volumeId, create, nil)
	return err
}

func (c *STaaSEnvironmentServiceImpl) DeleteVolume(id string, volumeId string, delete STaaSVolumeDelete) error {
	err := c.client.Delete(staasBasePath+"/environment/"+id+"/volume/"+volumeId, delete)
	return err
}

func (c *STaaSEnvironmentServiceImpl) CreateNetwork(id string, create STaaSNetworkCreate) error {
	err := c.client.Post(staasBasePath+"/environment/"+id+"/network", create, nil)
	return err
}

func (c *STaaSEnvironmentServiceImpl) DeleteNetwork(id string, networkId string) error {
	err := c.client.Delete(staasBasePath+"/environment/"+id+"/network/"+networkId, nil)
	return err
}
