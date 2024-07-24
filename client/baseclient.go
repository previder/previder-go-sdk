package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultBaseURL     = "https://portal.previder.com/api/"
	iaasBasePath       = "v2/iaas/"
	kubernetesBasePath = "v2/kubernetes/"
	staasBasePath      = "v2/storage/staas/"
	jsonEncoding       = "application/json; charset=utf-8"
	customerHeader     = "X-CustomerId"
)

type BaseClient struct {
	httpClient        *http.Client
	clientOptions     *ClientOptions
	Task              TaskService
	VirtualMachine    VirtualMachineService
	VirtualNetwork    VirtualNetworkService
	KubernetesCluster KubernetesClusterService
	STaaSEnvironment  STaaSEnvironmentService
}

type ApiInfo struct {
	Version string `json:"version,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ApiError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ApiErrorResponseBody struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Path    string `json:"path,omitempty"`
}

type ClientOptions struct {
	Token      string
	BaseUrl    string
	CustomerId string
}

type OwnerReference struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

// noinspection GoUnusedExportedFunction
func New(options *ClientOptions) (*BaseClient, error) {
	if options.Token == "" {
		return nil, fmt.Errorf("missing token")
	}
	if options.BaseUrl == "" {
		options.BaseUrl = defaultBaseURL
	}

	c := &BaseClient{httpClient: http.DefaultClient, clientOptions: options}

	c.Task = &TaskServiceOp{client: c}
	c.VirtualMachine = &VirtualMachineServiceOp{client: c}
	c.VirtualNetwork = &VirtualNetworkServiceOp{client: c}
	c.KubernetesCluster = &KubernetesClusterServiceOp{client: c}
	c.STaaSEnvironment = &STaaSEnvironmentServiceOp{client: c}
	return c, nil
}

func (c *BaseClient) Get(url string, responseBody interface{}, requestParams *PageRequest) error {
	return c.request("GET", url, nil, requestParams, &responseBody)
}

func (c *BaseClient) Delete(url string, responseBody interface{}) error {
	return c.request("DELETE", url, nil, nil, &responseBody)
}

func (c *BaseClient) Post(url string, requestBody, responseBody interface{}) error {
	return c.request("POST", url, &requestBody, nil, &responseBody)
}

func (c *BaseClient) Put(url string, requestBody, responseBody interface{}) error {
	return c.request("PUT", url, &requestBody, nil, &responseBody)
}

func (c *BaseClient) request(method string, url string, requestBody interface{}, pageRequest *PageRequest, responseBody interface{}) error {

	// content will be empty with GET, so can be sent anyway
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.clientOptions.BaseUrl+url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", jsonEncoding)

	req.Header.Set("X-Auth-Token", c.clientOptions.Token)

	req.Header.Set("Accept", jsonEncoding)

	if c.clientOptions.CustomerId != "" && len(c.clientOptions.CustomerId) == 24 {
		req.Header.Set(customerHeader, c.clientOptions.CustomerId)
	}
	if pageRequest != nil {
		q := req.URL.Query()
		q.Add("size", strconv.Itoa(pageRequest.Size))
		q.Add("page", strconv.Itoa(pageRequest.Page))
		q.Add("sort", pageRequest.Sort)
		q.Add("query", pageRequest.Query)
		req.URL.RawQuery = q.Encode()
	}
	res, requestErr := c.httpClient.Do(req)
	if requestErr != nil {
		log.Printf("[ERROR] [Previder API] Error from Previder API received")
		return requestErr
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		apiError := new(ApiError)

		var apiErrorResponseBody ApiErrorResponseBody
		temp, err := io.ReadAll(res.Body)
		err = json.Unmarshal(temp, &apiErrorResponseBody)
		var tmpBuffer any
		json.Unmarshal(temp, &tmpBuffer)
		fmt.Println(tmpBuffer)

		if err != nil {
			log.Printf("[ERROR] [Previder API] Could not parse error result:" + string(temp))
			return err
		}
		log.Printf("[ERROR] [Previder API] Error while executing the request to " + apiErrorResponseBody.Path + ": [" + strconv.Itoa(apiErrorResponseBody.Status) + "] " + apiErrorResponseBody.Message)
		apiError.Code = res.StatusCode
		apiError.Message = "[Previder API] " + apiErrorResponseBody.Message
		if apiErrorResponseBody.Error != "" {
			apiError.Message = apiError.Message + " - " + apiErrorResponseBody.Error
		}
		return apiError
	}

	if responseBody != nil {
		err := json.NewDecoder(res.Body).Decode(&responseBody)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}

func (c *BaseClient) ApiInfo() (*ApiInfo, error) {
	apiInfo := new(ApiInfo)
	err := c.Get("", apiInfo, nil)
	return apiInfo, err
}
