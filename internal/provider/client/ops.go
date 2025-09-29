package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Getops - Returns list of ops (no auth required)
func (c *Client) GetOps() ([]Ops, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/op", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	ops := []Ops{}
	err = json.Unmarshal(body, &ops)
	if err != nil {
		return nil, err
	}

	return ops, nil
}

// Getop - Returns specific op (no auth required)
func (c *Client) GetOp(opID string) (*Ops, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/op/id/%s", c.HostURL, opID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	op := Ops{}
	err = json.Unmarshal(body, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

// Createop - Create new op
func (c *Client) CreateOps(op Ops) (*Ops, error) {
	rb, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/op", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newop := Ops{}
	err = json.Unmarshal(body, &newop)
	if err != nil {
		return nil, err
	}

	return &newop, nil
}

func (c *Client) UpdateOps(opID string, op Ops) (*Ops, error) {
	rb, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/op/%s", c.HostURL, opID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	var resp Ops
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DeleteOps(opID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/op/%s", c.HostURL, opID), nil)
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
