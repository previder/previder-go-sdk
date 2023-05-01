package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultBaseURL = "https://portal.previder.com/api/"
	iaasBasePath   = "v2/iaas/"
	jsonEncoding   = "application/json; charset=utf-8"
	customerHeader = "X-CustomerId"
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

	req.Header.Set("X-Auth-Token", c.clientOptions.Token)

	req.Header.Set("Accept", jsonEncoding)

	if c.clientOptions.CustomerId != "" && len(c.clientOptions.CustomerId) == 24 {
		req.Header.Set(customerHeader, c.clientOptions.CustomerId)
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
		temp, err := ioutil.ReadAll(res.Body)

		err = json.Unmarshal(temp, &apiErrorResponseBody)
		if err != nil {
			log.Printf("[ERROR] [Previder API] Could not parse error result:" + string(temp))
			return err
		}
		log.Printf("[ERROR] [Previder API] Error while executing the request to " + apiErrorResponseBody.Path + ": [" + strconv.Itoa(apiErrorResponseBody.Status) + "] " + apiErrorResponseBody.Message)
		apiError.Code = res.StatusCode
		apiError.Message = "[Previder API] " + apiErrorResponseBody.Message
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
	err := c.Get("", apiInfo)
	return apiInfo, err
}
