package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	defaultBaseURL = "https://portal.previder.nl/api/"
	iaasBasePath   = "v2/iaas/"
	jsonEncoding   = "application/json; charset=utf-8"
)

type BaseClient struct {
	httpClient     *http.Client
	clientOptions  *ClientOptions
	Task           TaskService
	VirtualMachine VirtualMachineService
	VirtualNetwork VirtualNetworkService
}

type ApiInfo struct {
	Version string `json:"result,omitempty"`
}

type ApiError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ClientOptions struct {
	Token   string
	BaseUrl string
}

type Page struct {
	TotalPages       int
	TotalElements    int
	NumberOfElements int
	Size             int
	Number           int
	Content          json.RawMessage
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

//noinspection GoUnusedExportedFunction
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
	return c, nil
}

func (c *BaseClient) Get(url string, responseBody interface{}) error {
	return c.request("GET", url, nil, &responseBody)
}

func (c *BaseClient) Delete(url string, responseBody interface{}) error {
	return c.request("DELETE", url, nil, &responseBody)
}

func (c *BaseClient) Post(url string, requestBody, responseBody interface{}) error {
	return c.request("POST", url, &requestBody, &responseBody)
}

func (c *BaseClient) Put(url string, requestBody, responseBody interface{}) error {
	return c.request("PUT", url, &requestBody, &responseBody)
}

func (c *BaseClient) request(method string, url string, requestBody, responseBody interface{}) error {

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

	req.Header.Set("Authorization", "Bearer "+c.clientOptions.Token)

	req.Header.Set("Accept", jsonEncoding)

	res, requestErr := c.httpClient.Do(req)
	if requestErr != nil {
		return requestErr
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		apiError := new(ApiError)
		_, err := ioutil.ReadAll(res.Body)
		if err != nil {

			return err
		}
		apiError.Code = res.StatusCode
		return apiError
	}

	if responseBody != nil {
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
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
	err := c.Get("", apiInfo)
	return apiInfo, err
}
