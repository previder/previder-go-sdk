package client

import "encoding/json"

const secureVaultPath = kubernetesBasePath + "securevault"

type SecureVaultService interface {
	Page(request PageRequest) (*Page, *[]SecureVault, error)
	Get(id string) (*SecureVaultExt, error)
	Create(create SecureVaultCreate) (*SecureVaultResponse, error)
	Update(id string, update SecureVaultCreate) error
	Delete(id string) error
}

type SecureVaultServiceImpl struct {
	client *PreviderClient
}

type SecureVault struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Group     string `json:"group,omitempty"`
	GroupName string `json:"groupName,omitempty"`
	State     string `json:"state"`
}

type SecureVaultExt struct {
	SecureVault
	TerminationProtected bool `json:"terminationProtected"`
}

type SecureVaultCreate struct {
	Name                 string `json:"name"`
	Group                string `json:"group,omitempty"`
	TerminationProtected bool   `json:"terminationProtected"`
}

type SecureVaultResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
	Token string `json:"token"`
}

func (c *SecureVaultServiceImpl) Page(request PageRequest) (*Page, *[]SecureVault, error) {
	page := new(Page)
	err := c.client.Get(secureVaultPath, page, &request)
	if err != nil {
		return nil, nil, err
	}

	vaults := new([]SecureVault)
	if err := json.Unmarshal(page.Content, &vaults); err != nil {
		return nil, nil, err
	}

	return page, vaults, err
}

func (c *SecureVaultServiceImpl) Get(id string) (*SecureVaultExt, error) {
	vault := new(SecureVaultExt)
	err := c.client.Get(secureVaultPath+"/"+id, vault, nil)
	return vault, err
}

func (c *SecureVaultServiceImpl) Create(create SecureVaultCreate) (*SecureVaultResponse, error) {
	response := new(SecureVaultResponse)
	err := c.client.Post(secureVaultPath, create, response)
	return response, err
}

func (c *SecureVaultServiceImpl) Update(id string, update SecureVaultCreate) error {
	err := c.client.Put(secureVaultPath+"/"+id, update, nil)
	return err
}

func (c *SecureVaultServiceImpl) Delete(id string) error {
	err := c.client.Delete(secureVaultPath+"/"+id, nil)
	return err
}
