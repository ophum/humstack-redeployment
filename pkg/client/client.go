package client

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack-redeployment/pkg/api"
)

type RedeploymentClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type RedeploymentResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		Redeployment api.Redeployment `json:"redeployment"`
	} `json:"data"`
}

type RedeploymentListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		RedeploymentList []*api.Redeployment `json:"redeployments"`
	} `json:"data"`
}

const basePath = "api/v0"

func NewRedeploymentClient(scheme, apiServerAddress string, apiServerPort int32) *RedeploymentClient {
	return &RedeploymentClient{
		scheme:           scheme,
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,
		client:           resty.New(),
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

func (c *RedeploymentClient) Get(rdID string) (*api.Redeployment, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(rdID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	rdResp := RedeploymentResponse{}
	err = json.Unmarshal(body, &rdResp)
	if err != nil {
		return nil, err
	}

	return &rdResp.Data.Redeployment, nil
}

func (c *RedeploymentClient) List() ([]*api.Redeployment, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	rdResp := RedeploymentListResponse{}
	err = json.Unmarshal(body, &rdResp)
	if err != nil {
		return nil, err
	}
	return rdResp.Data.RedeploymentList, nil
}

func (c *RedeploymentClient) Create(rd *api.Redeployment) (*api.Redeployment, error) {
	body, err := json.Marshal(rd)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	rdResp := RedeploymentResponse{}
	err = json.Unmarshal(body, &rdResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", rdResp.Error)
	}

	return &rdResp.Data.Redeployment, nil
}

func (c *RedeploymentClient) Update(rd *api.Redeployment) (*api.Redeployment, error) {
	body, err := json.Marshal(rd)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(rd.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	rdResp := RedeploymentResponse{}
	err = json.Unmarshal(body, &rdResp)
	if err != nil {
		return nil, err
	}

	return &rdResp.Data.Redeployment, nil
}

func (c *RedeploymentClient) Delete(rdID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(rdID))
	if err != nil {
		return err
	}

	return nil
}

func (c *RedeploymentClient) getPath(path string) string {
	return fmt.Sprintf("%s://%s", c.scheme, filepath.Join(fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort), basePath, path))
}
