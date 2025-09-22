package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetEngineers - Returns list of engineers (no auth required)
func (c *Client) GetEngineers() ([]Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineers := []Engineer{}
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

// GetEngineer - Returns specific engineer (no auth required)
func (c *Client) GetEngineer(engineerID string) ([]Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/%s", c.HostURL, engineerID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineers := []Engineer{}
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

// CreateEngineer - Create new Engineer
func (c *Client) CreateEngineer(engineer Engineer) (*Engineer, error) {
	rb, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newEngineer := Engineer{}
	err = json.Unmarshal(body, &newEngineer)
	if err != nil {
		return nil, err
	}

	return &newEngineer, nil
}
