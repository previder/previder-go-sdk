package client

import (
	"encoding/json"
	"fmt"
)

const sshKeyBasePath = "sshkey/"

type SSHKeyService interface {
	Page(pageNr int, size int) (*Page, []SSHKey, error)
	Get(id string) (*SSHKey, error)
	GetByFingerprint(fingerprint string) (*SSHKey, error)
	Create(create *SSHKeyCreate) (*SSHKey, error)
	Generate(create *SSHKeyCreate) (*SSHKey, error)
	Delete(id string) error
}

type SSHKeyServiceOp struct {
	client *BaseClient
}

type SSHKey struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PublicKey   string `json:"publicKey,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PrivateKey  string `json:"privateKey,omitempty"`
}

type SSHKeyCreate struct {
	Name      string `json:"name,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

func (c *SSHKeyServiceOp) Page(pageNr int, size int) (*Page, []SSHKey, error) {
	url := fmt.Sprintf("%s?page=%d&size=%d", sshKeyBasePath, pageNr, size)
	page := new(Page)

	if err := c.client.Get(url, page); err != nil {
		return nil, nil, err
	}

	sshKeys := new([]SSHKey)
	if err := json.Unmarshal([]byte(page.Content), &sshKeys); err != nil {
		return nil, nil, err
	}

	return page, *sshKeys, nil
}

func (c *SSHKeyServiceOp) Get(id string) (*SSHKey, error) {
	sshKey := new(SSHKey)
	err := c.client.Get(sshKeyBasePath+id, sshKey)
	return sshKey, err
}

func (c *SSHKeyServiceOp) GetByFingerprint(fingerprint string) (*SSHKey, error) {
	sshKey := new(SSHKey)
	err := c.client.Get(sshKeyBasePath+"fingerprint/"+fingerprint, sshKey)
	return sshKey, err
}

func (c *SSHKeyServiceOp) Create(create *SSHKeyCreate) (*SSHKey, error) {
	sshKey := new(SSHKey)
	err := c.client.Post(sshKeyBasePath, create, sshKey)
	return sshKey, err
}

func (c *SSHKeyServiceOp) Generate(create *SSHKeyCreate) (*SSHKey, error) {
	sshKey := new(SSHKey)
	err := c.client.Post(sshKeyBasePath+"generate", create, sshKey)
	return sshKey, err
}

func (c *SSHKeyServiceOp) Delete(id string) error {
	return c.client.Delete(sshKeyBasePath+id, nil)
}
