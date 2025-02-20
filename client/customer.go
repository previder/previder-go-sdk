package client

import "encoding/json"

type CustomerService interface {
	Page(request PageRequest) (*Page, *[]Customer, error)
	Get(id string) (*CustomerExt, error)
	Create(customerCreate CustomerCreate) (*Customer, error)
	Delete(id string) error
	Update(id string, customerUpdate CustomerCreate) (*Customer, error)
}

type CustomerServiceImpl struct {
	client *PreviderClient
}

type Customer struct {
	Id                  string `json:"id,omitempty"`
	Name                string `json:"name"`
	AccountName         string `json:"accountName,omitempty"`
	VerificationStatus  string `json:"verificationStatus,omitempty"`
	City                string `json:"city,omitempty"`
	CountryCode         string `json:"countryCode,omitempty"`
	Language            string `json:"language,omitempty"`
	OcfId               int    `json:"ocfId,omitempty"`
	NfaId               string `json:"nfaId,omitempty"`
	Partner             bool   `json:"partner"`
	HidingPrices        bool   `json:"hidingPrices"`
	PurchaseOrderNumber string `json:"purchaseOrderNumber,omitempty"`
	InvoiceToPartner    bool   `json:"invoiceToPartner"`
}

type CustomerExt struct {
	Customer
	Address       string   `json:"address,omitempty"`
	AddressNumber string   `json:"addressNumber,omitempty"`
	AddressSuffix string   `json:"addressSuffix,omitempty"`
	Zipcode       string   `json:"zipCode,omitempty"`
	CocNumber     string   `json:"cocNumber,omitempty"`
	ParentId      string   `json:"parentId,omitempty"`
	ParentName    string   `json:"parentName,omitempty"`
	ParentTree    []string `json:"parentTree,omitempty"`
	ParentOcfId   int      `json:"parentOcfId,omitempty"`
	ParentNfaId   string   `json:"parentNfaId,omitempty"`
	AccountType   string   `json:"accountType,omitempty"`
	Audit         Audit    `json:"audit"`
}

type CustomerCreate struct {
	Name                string `json:"name"`
	Address             string `json:"address,omitempty"`
	AccountName         string `json:"accountName,omitempty"`
	AddressNumber       string `json:"addressNumber,omitempty"`
	AddressSuffix       string `json:"addressSuffix,omitempty"`
	Zipcode             string `json:"zipCode,omitempty"`
	City                string `json:"city,omitempty"`
	CountryCode         string `json:"countryCode,omitempty"`
	Language            string `json:"language,omitempty"`
	PurchaseOrderNumber string `json:"purchaseOrderNumber,omitempty"`
	CocNumber           string `json:"cocNumber,omitempty"`
	Partner             bool   `json:"partner,omitempty"`
	HidingPrices        bool   `json:"hidingPrices,omitempty"`
	InvoiceToPartner    bool   `json:"invoiceToPartner,omitempty"`
}

func (c CustomerServiceImpl) Page(request PageRequest) (*Page, *[]Customer, error) {
	page := new(Page)
	err := c.client.Get(coreBasePath+"customer", page, &request)
	if err != nil {
		return nil, nil, err
	}

	customers := new([]Customer)
	if err := json.Unmarshal(page.Content, &customers); err != nil {
		return nil, nil, err
	}

	return page, customers, err
}

func (c CustomerServiceImpl) Get(id string) (*CustomerExt, error) {
	customer := new(CustomerExt)
	err := c.client.Get(coreBasePath+"customer/"+id, customer, nil)
	return customer, err
}

func (c CustomerServiceImpl) Create(customerCreate CustomerCreate) (*Customer, error) {
	customer := new(Customer)
	err := c.client.Post(coreBasePath+"customer", customerCreate, customer)
	return customer, err
}

func (c CustomerServiceImpl) Delete(id string) error {
	err := c.client.Delete(coreBasePath+"customer/"+id, nil)
	return err
}

func (c CustomerServiceImpl) Update(id string, customerUpdate CustomerCreate) (*Customer, error) {
	customer := new(Customer)
	err := c.client.Put(coreBasePath+id, customerUpdate, customer)
	return customer, err
}
