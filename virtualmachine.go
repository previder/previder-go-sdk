package client

import (
	"errors"
	"strings"
	"fmt"
)

//noinspection GoUnusedConst
const (
	virtualMachineBasePath = "virtualmachine/"
	VmActionPowerOn = "POWERON"
	VmActionPowerOff = "POWEROFF"
	VmActionShutdown = "SHUTDOWN"
	VmActionReboot = "REBOOT"
	VmActionSuspend = "SUSPEND"
	VmActionReset = "RESET"
)

type VirtualMachineService interface {
	List() (*[]VirtualMachine, error)
	Get(id string) (*VirtualMachine, error)
	Create(vm *VirtualMachineCreate) (*Task, error)
	Delete(id string) (*Task, error)
	Update(id string, vm *VirtualMachineUpdate) (*Task, error)
	Control(id string, action string) (*Task, error)

	CreateDisk(id string, disk *VirtualDisk) (*Task, error)
	UpdateDisk(id string, diskId string, disk *VirtualDisk) (*Task, error)
	DeleteDisk(id string, diskId string) (*Task, error)

	CreateNetworkInterface(id string, nic *NetworkInterface) (*Task, error)
	UpdateNetworkInterface(id string, nicId string, nic *NetworkInterface) (*Task, error)
	DeleteNetworkInterface(id string, nicId string) (*Task, error)

	ComputeClusterList() (*[]ComputeCluster, error)
	VirtualMachineTemplateList() (*[]VirtualMachineTemplate, error)
	ComputeClusterGet(computeClusterName string) (*ComputeCluster, error)
	VirtualMachineTemplateGet(templateName string) (*VirtualMachineTemplate, error)
}

type VirtualMachineServiceOp struct {
	client *BaseClient
}

type VirtualMachine struct {
	Id                    string                 `json:"id,omitempty"`
	Name                  string                 `json:"name,omitempty"`
	MemoryMb              int                    `json:"memoryMb,omitempty"`
	CpuCores              int                    `json:"cpuCores,omitempty"`
	VirtualDisks          []VirtualDisk          `json:"virtualDisks,omitempty"`
	NetworkInterfaces     []NetworkInterface     `json:"networkInterfaces,omitempty"`
	UserData              string                 `json:"userData,omitempty"`
	ProvisioningType      string                 `json:"provisioningType,omitempty"`
	Template              VirtualMachineTemplate `json:"template,omitempty"`
	ComputeCluster        ComputeCluster         `json:"computeCluster,omitempty"`
	State                 string                 `json:"state,omitempty"`
	Hostname              string                 `json:"hostname,omitempty"`
	TerminationProtection bool                   `json:"terminationProtection,omitempty"`
	InitialPassword		  string				 `json:"initialPassword,omitempty"`
}

type VirtualMachineTemplate struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version int    `json:"version,omitempty"`
}

type ComputeCluster struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type VirtualDisk struct {
	Id         string `json:"id,omitempty"`
	DiskSizeMb int    `json:"diskSize,omitempty"`
}

type NetworkInterface struct {
	Id                         string              `json:"id,omitempty"`
	Network                    VirtualNetwork      `json:"network,omitempty"`
	Connected                  bool                `json:"connected,omitempty"`
	AddressAssignments         []AddressAssignment `json:"addressAssignments,omitempty"`
	FirstIPv4AddressAssignment AddressAssignment   `json:"firstIPv4AddressAssignment,omitempty"`
	FirstIPv6AddressAssignment AddressAssignment   `json:"firstIPv6AddressAssignment,omitempty"`
}

type AddressAssignment struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
}

type VirtualMachineCreate struct {
	Name              string                 `json:"name,omitempty"`
	MemoryMb          int                    `json:"memoryMb,omitempty"`
	CpuCores          int                    `json:"cpuCores,omitempty"`
	VirtualDisks      []VirtualDisk          `json:"virtualDisks,omitempty"`
	NetworkInterfaces []NetworkInterface     `json:"networkInterfaces,omitempty"`
	UserData          string                 `json:"userData,omitempty"`
	ProvisioningType  string                 `json:"provisioningType,omitempty"`
	Template          VirtualMachineTemplate `json:"template,omitempty"`
	ComputeCluster    ComputeCluster         `json:"computeCluster,omitempty"`
	SSHKeys           []SSHKey               `json:"sshKeys,omitempty"`
}

type VirtualMachineUpdate struct {
	Name           string `json:"name,omitempty"`
	MemoryMb       int    `json:"memoryMb,omitempty"`
	CpuCores       int    `json:"cpuCores,omitempty"`
	// Some fields for V1 api
	ComputeCluster ComputeCluster `json:"computeCluster,omitempty"`
	Hostname       string         `json:"hostname,omitempty"`
}

func (c *VirtualMachineServiceOp) ComputeClusterList() (*[]ComputeCluster, error) {
	computeClusters := new([]ComputeCluster)
	err := c.client.Get(virtualMachineBasePath + "cluster", computeClusters)
	return computeClusters, err
}

func (c *VirtualMachineServiceOp) VirtualMachineTemplateList() (*[]VirtualMachineTemplate, error) {
	virtualMachineTemplates := new([]VirtualMachineTemplate)
	err := c.client.Get(virtualMachineBasePath + "template", virtualMachineTemplates)
	return virtualMachineTemplates, err
}

func (c *VirtualMachineServiceOp) List() (*[]VirtualMachine, error) {
	virtualMachines := new([]VirtualMachine)
	err := c.client.Get(virtualMachineBasePath, virtualMachines)
	return virtualMachines, err
}

func (c *VirtualMachineServiceOp) Get(id string) (*VirtualMachine, error) {
	virtualMachine := new(VirtualMachine)
	err := c.client.Get(virtualMachineBasePath + id, virtualMachine)
	return virtualMachine, err
}

func (c *VirtualMachineServiceOp) Create(vm *VirtualMachineCreate) (*Task, error) {
	task := new(Task)
	err := c.client.Post(virtualMachineBasePath, vm, task)
	return task, err
}

func (c *VirtualMachineServiceOp) Update(id string, vm *VirtualMachineUpdate) (*Task, error) {
	task := new(Task)
	err := c.client.Put(virtualMachineBasePath + id, vm, task)
	return task, err
}

func (c *VirtualMachineServiceOp) Control(id string, action string) (*Task, error) {
	task := new(Task)
	err := c.client.Post(virtualMachineBasePath + id + "/control/" + action, nil, task)
	return task, err
}

func (c *VirtualMachineServiceOp) Delete(id string) (*Task, error) {
	task := new(Task)
	err := c.client.Delete(virtualMachineBasePath + id, task)
	return task, err
}

func (c *VirtualMachineServiceOp) CreateDisk(id string, disk *VirtualDisk) (*Task, error) {
	task := new(Task)
	err := c.client.Post(virtualMachineBasePath + id + "/disk", disk, task)
	return task, err
}

func (c *VirtualMachineServiceOp) UpdateDisk(id string, diskId string, disk *VirtualDisk) (*Task, error) {
	task := new(Task)
	err := c.client.Put(virtualMachineBasePath + id + "/disk/" + diskId, disk, task)
	return task, err
}

func (c *VirtualMachineServiceOp) DeleteDisk(id string, diskId string) (*Task, error) {
	task := new(Task)
	err := c.client.Delete(virtualMachineBasePath + id + "/disk/" + diskId, task)
	return task, err
}

func (c *VirtualMachineServiceOp) CreateNetworkInterface(id string, nic *NetworkInterface) (*Task, error) {
	task := new(Task)
	err := c.client.Post(virtualMachineBasePath + id + "/nic", nic, task)
	return task, err
}

func (c *VirtualMachineServiceOp) UpdateNetworkInterface(id string, nicId string, nic *NetworkInterface) (*Task, error) {
	task := new(Task)
	err := c.client.Put(virtualMachineBasePath + id + "/nic/" + nicId, nic, task)
	return task, err
}

func (c *VirtualMachineServiceOp) DeleteNetworkInterface(id string, nicId string) (*Task, error) {
	task := new(Task)
	err := c.client.Delete(virtualMachineBasePath + id + "/nic/" + nicId, task)
	return task, err
}

func (c *VirtualMachineServiceOp) ComputeClusterGet(computeClusterName string) (*ComputeCluster, error) {
	computeClusters, err := c.ComputeClusterList()
	if err != nil {
		return nil, err
	}
	for _, cc := range *computeClusters {
		if strings.EqualFold(cc.Name, computeClusterName) {
			return &cc, nil
		}
	}
	return nil, errors.New("ComputeCluster not found")
}

func (c *VirtualMachineServiceOp) VirtualMachineTemplateGet(selector string) (*VirtualMachineTemplate, error) {
	templates, err := c.VirtualMachineTemplateList()
	if err != nil {
		return nil, err
	}

	scores := make([]int, len(*templates))

	selector = strings.ToLower(selector)

	// Calculate a match score on each template
	for i, template := range *templates {
		templateName := strings.ToLower(template.Name)
		scores[i] = 0
		for _, attr := range strings.Split(templateName, " ") {
			if strings.HasPrefix(attr, "(") && strings.HasSuffix(attr, ")") {
				feature := strings.TrimPrefix(strings.TrimSuffix(attr, ")"), "(")
				if !strings.Contains(selector, feature) {
					scores[i] = 0
					break
				} else {
					scores[i]++
				}
			}
			if strings.Contains(selector, attr) {
				scores[i]++
			}
		}
	}

	// Get the template with the highest score and version
	highestScore := 0
	highestVersion := 0
	var bestTemplateMatch *VirtualMachineTemplate
	for i := range *templates {
		if scores[i] > highestScore {
			highestScore = scores[i]
			highestVersion = (*templates)[i].Version
			bestTemplateMatch = &(*templates)[i]
		}
		if highestScore > 0 && scores[i] == highestScore && (*templates)[i].Version > highestVersion {
			highestVersion = (*templates)[i].Version
			bestTemplateMatch = &(*templates)[i]
		}
	}

	if bestTemplateMatch == nil {
		return nil, errors.New("No template found")
	}
	fmt.Printf("Selecting template %s based on %s\n", bestTemplateMatch.Name, selector)

	return bestTemplateMatch, nil
}
