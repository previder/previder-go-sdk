package client

import (
	"encoding/json"
)

// noinspection GoUnusedConst
const (
	VmActionPowerOn  = "POWERON"
	VmActionPowerOff = "POWEROFF"
	VmActionShutdown = "SHUTDOWN"
	VmActionReboot   = "REBOOT"
	VmActionSuspend  = "SUSPEND"
	VmActionReset    = "RESET"

	VmStateNew        = "NEW"
	VmStateDeploying  = "DEPLOYING"
	VmStatePoweredOff = "POWEREDOFF"
	VmStatePoweredOn  = "POWEREDON"
)

type VirtualServerService interface {
	ComputeClusterList() (*[]ComputeCluster, error)
	VirtualMachineTemplateList() (*[]VirtualMachineTemplate, error)
	Page(request PageRequest) (*Page, *[]VirtualMachine, error)
	Get(id string) (*VirtualMachineExt, error)
	Create(vm *VirtualMachineCreate) (*VirtualMachineTask, error)
	Delete(id string) (*VirtualMachineTask, error)
	Update(id string, vm *VirtualMachineUpdate) (*VirtualMachineTask, error)
	Control(id string, action string) (*VirtualMachineTask, error)
	OpenConsole(id string) (*OpenConsoleResult, error)
}

type VirtualServerServiceImpl struct {
	client *PreviderClient
}

type VirtualMachineTask struct {
	Task
	VirtualMachine     string
	VirtualMachineName string
}

type VirtualMachine struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name"`
	Group            string `json:"group,omitempty"`
	GroupName        string `json:"groupName,omitempty"`
	ComputeCluster   string `json:"computeCluster"`
	CpuCores         int    `json:"cpuCores"`
	Memory           uint64 `json:"memory"`
	Template         string `json:"template"`
	GuestId          string `json:"guestId"`
	State            string `json:"state"`
	TotalDiskSize    int    `json:"totalDiskSize"`
	HasSnapshots     bool   `json:"hasSnapshots"`
	MarkedAsTemplate bool   `json:"markedAsTemplate"`
	Managed          bool   `json:"managed"`
}

type VirtualMachineExt struct {
	VirtualMachine
	Hostname                     string             `json:"hostname"`
	Tags                         []string           `json:"tags"`
	Disks                        []Disk             `json:"disks"`
	NetworkInterfaces            []NetworkInterface `json:"networkInterfaces"`
	TerminationProtectionEnabled bool               `json:"terminationProtectionEnabled"`
	Flavor                       string             `json:"flavor,omitempty"`
	GuestToolsStatus             string             `json:"guestToolsStatus"`
	InitialUsername              string             `json:"initialUsername"`
	InitialPassword              string             `json:"initialPassword"`
	CreatedAt                    int                `json:"createdAt"`
	CreatedBy                    string             `json:"createdBy"`
	LastModifiedAt               int                `json:"lastModifiedAt"`
	LastModifiedBy               string             `json:"lastModifiedBy"`
}

type VirtualMachineUpdate struct {
	VirtualMachine
	Tags                         []string                 `json:"tags"`
	Disks                        []DiskUpdate             `json:"disks"`
	NetworkInterfaces            []NetworkInterfaceUpdate `json:"networkInterfaces,"`
	TerminationProtectionEnabled bool                     `json:"terminationProtectionEnabled"`
	Flavor                       string                   `json:"flavor,omitempty"`
}

type Disk struct {
	Id    string `json:"id,omitempty"`
	Size  uint64 `json:"size"`
	Uuid  string `json:"uuid,omitempty"`
	Label string `json:"label,omitempty"`
}

type DiskUpdate struct {
	Id     string `json:"id,omitempty"`
	Size   uint64 `json:"size"`
	Uuid   string `json:"uuid,omitempty"`
	Label  string `json:"label,omitempty"`
	Delete bool   `json:"delete,omitempty"`
}

type NetworkInterface struct {
	Id                  string   `json:"id,omitempty"`
	Network             string   `json:"network"`
	NetworkName         string   `json:"networkName,omitempty"`
	Connected           bool     `json:"connected"`
	MacAddress          string   `json:"macAddress,omitempty"`
	DiscoveredAddresses []string `json:"discoveredAddresses,omitempty"`
	AssignedAddresses   []string `json:"assignedAddresses,omitempty"`
	Primary             bool     `json:"primary,omitempty"`
	Label               string   `json:"label,omitempty"`
	Type                string   `json:"type,omitempty"`
}

type NetworkInterfaceUpdate struct {
	Id        string `json:"id,omitempty"`
	Network   string `json:"network"`
	Connected bool   `json:"connected"`
	Label     string `json:"label,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
}

type VirtualMachineCreate struct {
	VirtualMachineExt
	Template             string `json:"template,omitempty"`
	SourceVirtualMachine string `json:"sourceVirtualMachine,omitempty"`
	UserData             string `json:"userData,omitempty"`
	GuestId              string `json:"guestId,omitempty"`
	ProvisioningType     string `json:"provisioningType,omitempty"`
	PowerOnAfterClone    bool   `json:"powerOnAfterClone,omitempty"`
}

type VirtualMachineTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     int    `json:"version"`
	Category    string `json:"category,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type ComputeCluster struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OpenConsoleResult struct {
	ConsoleUrl string `json:"consoleUrl,omitempty"`
}

func (c *VirtualServerServiceImpl) ComputeClusterList() (*[]ComputeCluster, error) {
	computeClusters := new([]ComputeCluster)
	err := c.client.Get(iaasBasePath+"computecluster", computeClusters, nil)
	return computeClusters, err
}

func (c *VirtualServerServiceImpl) VirtualMachineTemplateList() (*[]VirtualMachineTemplate, error) {
	virtualMachineTemplates := new([]VirtualMachineTemplate)
	err := c.client.Get(iaasBasePath+"template", virtualMachineTemplates, nil)
	return virtualMachineTemplates, err
}

func (c *VirtualServerServiceImpl) Page(request PageRequest) (*Page, *[]VirtualMachine, error) {
	page := new(Page)
	err := c.client.Get(iaasBasePath+"virtualmachine", page, &request)
	if err != nil {
		return nil, nil, err
	}

	virtualMachines := new([]VirtualMachine)
	if err := json.Unmarshal(page.Content, &virtualMachines); err != nil {
		return nil, nil, err
	}

	return page, virtualMachines, err
}

func (c *VirtualServerServiceImpl) Get(id string) (*VirtualMachineExt, error) {
	virtualMachine := new(VirtualMachineExt)
	err := c.client.Get(iaasBasePath+"virtualmachine/"+id, virtualMachine, nil)
	return virtualMachine, err
}

func (c *VirtualServerServiceImpl) Create(vm *VirtualMachineCreate) (*VirtualMachineTask, error) {
	task := new(VirtualMachineTask)
	err := c.client.Post(iaasBasePath+"virtualmachine", vm, task)
	return task, err
}

func (c *VirtualServerServiceImpl) Update(id string, vm *VirtualMachineUpdate) (*VirtualMachineTask, error) {
	task := new(VirtualMachineTask)
	err := c.client.Put(iaasBasePath+"virtualmachine/"+id, vm, task)
	return task, err
}

func (c *VirtualServerServiceImpl) Delete(id string) (*VirtualMachineTask, error) {
	task := new(VirtualMachineTask)
	err := c.client.Delete(iaasBasePath+"virtualmachine/"+id, task)
	return task, err
}

func (c *VirtualServerServiceImpl) Control(id string, action string) (*VirtualMachineTask, error) {
	task := new(VirtualMachineTask)
	err := c.client.Post(iaasBasePath+"virtualmachine/"+id+"/action/"+action, nil, task)
	return task, err
}

func (c *VirtualServerServiceImpl) OpenConsole(id string) (*OpenConsoleResult, error) {
	res := new(OpenConsoleResult)
	err := c.client.Post(iaasBasePath+"virtualmachine/"+id+"/console", nil, res)
	return res, err
}
