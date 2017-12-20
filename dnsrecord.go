package client

import "fmt"

const (
	dnsRecordBasePath = "dns/"
)

type DNSRecordService interface {
	Get(id string) (*DomainZone, error)
	Update(id string, update *DomainZoneUpdate) error
}

type DNSRecordServiceOp struct {
	client *BaseClient
}

type DomainZone struct {
	Name    string         `json:"name,omitempty"`
	Records []DomainRecord `json:"records,omitempty"`
}

type Domain struct {
	Id         string `json:"id,omitempty"`
	DomainName string `json:"name,omitempty"`
}

type DomainRecord struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	Type    string `json:"type,omitempty"`
	Ttl     int    `json:"ttl,omitempty"`
	Prio    int    `json:"prio,omitempty"`
}

type DomainZoneUpdate struct {
	Add    []DomainRecord `json:"add"`
	Update []DomainRecord `json:"update"`
	Remove []DomainRecord `json:"remove"`
}

func (c *DNSRecordServiceOp) getDomainId(name string) (string, error) {
	domains := new([]Domain)
	err := c.client.Get("/drs/domain", domains)
	if err != nil {
		return "", err
	}
	for _, domain := range *domains {
		if domain.DomainName == name || (domain.DomainName+".") == name {
			return domain.Id, nil
		}
	}
	return "", fmt.Errorf("Domain %s not found", name)
}

func (c *DNSRecordServiceOp) Get(name string) (*DomainZone, error) {
	id, err := c.getDomainId(name)
	if err != nil {
		return nil, err
	}
	domainZone := new(DomainZone)
	err = c.client.Get(dnsRecordBasePath+id, domainZone)
	return domainZone, err
}

func (c *DNSRecordServiceOp) Update(name string, update *DomainZoneUpdate) error {
	id, err := c.getDomainId(name)
	if err != nil {
		return err
	}
	return c.client.Post(dnsRecordBasePath+id, update, nil)
}
