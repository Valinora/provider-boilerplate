package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetDevs - Returns list of Devs (no auth required)
func (c *Client) GetDevs() ([]Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	devs := []Dev{}
	err = json.Unmarshal(body, &devs)
	if err != nil {
		return nil, err
	}

	return devs, nil
}

// GetDev - Returns specific Dev (no auth required)
func (c *Client) GetDev(devID string) (*Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev/id/%s", c.HostURL, devID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dev := Dev{}
	err = json.Unmarshal(body, &dev)
	if err != nil {
		return nil, err
	}

	return &dev, nil
}

// CreateDev - Create new Dev
func (c *Client) CreateDev(dev Dev) (*Dev, error) {
	rb, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dev", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newDev := Dev{}
	err = json.Unmarshal(body, &newDev)
	if err != nil {
		return nil, err
	}

	return &newDev, nil
}

func (c *Client) UpdateDev(devID string, dev Dev) (*Dev, error) {
	rb, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/dev/%s", c.HostURL, devID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	var resp Dev
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DeleteDev(DevID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/dev/%s", c.HostURL, DevID), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// Check if response contains "resource deleted"
	if !strings.Contains(string(body), "resource deleted") {
		return fmt.Errorf("unexpected response: resource deletion not confirmed")
	}

	return nil
}
