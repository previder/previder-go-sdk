package client

import (
	"encoding/json"
)

type KubernetesClusterService interface {
	Page(request PageRequest) (*Page, *[]KubernetesCluster, error)
	Get(id string) (*KubernetesClusterExt, error)
	Create(create KubernetesClusterCreate) error
	Delete(id string) error
	Update(id string, update KubernetesClusterUpdate) error
	GetKubeConfig(id string, endpoint string) (KubernetesClusterKubeConfigResponse, error)
}

type KubernetesClusterServiceImpl struct {
	client *BaseClient
}

type KubernetesCluster struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Version string `json:"version"`
}

type KubernetesClusterExt struct {
	KubernetesCluster
	Vips                      []string `json:"vips"`
	Endpoints                 []string `json:"endpoints"`
	MinimalNodes              int      `json:"minimalNodes"`
	MaximalNodes              int      `json:"maximalNodes"`
	AutoUpdate                bool     `json:"autoUpdate"`
	AutoScaleEnabled          bool     `json:"autoScaleEnabled"`
	ControlPlaneCpuCores      int      `json:"controlPlaneCpuCores"`
	ControlPlaneMemoryGb      int      `json:"controlPlaneMemoryGb"`
	ControlPlaneStorageGb     int      `json:"controlPlaneStorageGb"`
	NodeCpuCores              int      `json:"nodeCpuCores"`
	NodeMemoryGb              int      `json:"nodeMemoryGb"`
	NodeStorageGb             int      `json:"nodeStorageGb"`
	ComputeCluster            string   `json:"computeCluster"`
	CNI                       string   `json:"cni"`
	HighAvailableControlPlane bool     `json:"highAvailableControlPlane"`
	Network                   string   `json:"network"`
	Reference                 string   `json:"reference"`
}

type KubernetesClusterCreate struct {
	KubernetesClusterUpdate
	Vips      []string `json:"vips"`
	Endpoints []string `json:"endpoints,omitempty"`
	CNI       string   `json:"cni"`
	Network   string   `json:"network"`
}

type KubernetesClusterUpdate struct {
	Name                      string `json:"name"`
	Version                   string `json:"version,omitempty"`
	MinimalNodes              int    `json:"minimalNodes"`
	MaximalNodes              int    `json:"maximalNodes,omitempty"`
	AutoUpdate                bool   `json:"autoUpdate"`
	AutoScaleEnabled          bool   `json:"autoScaleEnabled"`
	ControlPlaneCpuCores      int    `json:"controlPlaneCpuCores"`
	ControlPlaneMemoryGb      int    `json:"controlPlaneMemoryGb"`
	ControlPlaneStorageGb     int    `json:"controlPlaneStorageGb"`
	NodeCpuCores              int    `json:"nodeCpuCores"`
	NodeMemoryGb              int    `json:"nodeMemoryGb"`
	NodeStorageGb             int    `json:"nodeStorageGb"`
	ComputeCluster            string `json:"computeCluster"`
	HighAvailableControlPlane bool   `json:"highAvailableControlPlane"`
}

type KubernetesClusterKubeConfigRequest struct {
	Endpoint string `json:"endpoint"`
}

type KubernetesClusterKubeConfigResponse struct {
	Config string `json:"config"`
}

func (c *KubernetesClusterServiceImpl) Page(request PageRequest) (*Page, *[]KubernetesCluster, error) {
	page := new(Page)
	err := c.client.Get(kubernetesBasePath+"cluster", page, &request)
	if err != nil {
		return nil, nil, err
	}

	clusters := new([]KubernetesCluster)
	if err := json.Unmarshal(page.Content, &clusters); err != nil {
		return nil, nil, err
	}

	return page, clusters, err
}

func (c *KubernetesClusterServiceImpl) Get(id string) (*KubernetesClusterExt, error) {
	cluster := new(KubernetesClusterExt)
	err := c.client.Get(kubernetesBasePath+"cluster/"+id, cluster, nil)
	return cluster, err
}

func (c *KubernetesClusterServiceImpl) Create(create KubernetesClusterCreate) error {
	err := c.client.Post(kubernetesBasePath+"cluster", create, nil)
	return err
}

func (c *KubernetesClusterServiceImpl) Update(id string, update KubernetesClusterUpdate) error {
	err := c.client.Put(kubernetesBasePath+"cluster/"+id, update, nil)
	return err
}

func (c *KubernetesClusterServiceImpl) Delete(id string) error {
	err := c.client.Delete(kubernetesBasePath+"cluster/"+id, nil)
	return err
}

func (c *KubernetesClusterServiceImpl) GetKubeConfig(id string, endpoint string) (KubernetesClusterKubeConfigResponse, error) {
	requestKubeConfig := KubernetesClusterKubeConfigRequest{Endpoint: endpoint}
	var response KubernetesClusterKubeConfigResponse
	err := c.client.Post(kubernetesBasePath+"cluster/"+id+"/config", requestKubeConfig, &response)
	return response, err
}
