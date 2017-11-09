package config

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Client ...
type Client struct {
	EndPoint string
	Client   http.Client
}

// List ...
func (c *Client) List(args ListArgs) (*ListReply, error) {
	type body struct {
		ID     int           `json:"id"`
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}
	buf, err := json.Marshal(body{
		ID:     1,
		Method: "ConfigService.List",
		Params: []interface{}{args},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.EndPoint, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	type response struct {
		ID     int       `json:"id"`
		Error  string    `json:"error"`
		Result ListReply `json:"result"`
	}
	tmp := response{}
	err = json.NewDecoder(resp.Body).Decode(&tmp)
	return &ListReply{
		Configs: tmp.Result.Configs,
	}, err
}
