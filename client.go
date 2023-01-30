package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

type Request struct {
	Lights []*Light `json:"lights"`
}

type Light struct {
	On          *int `json:"on,omitempty"`
	Brightness  *int `json:"brightness,omitempty"`
	Temperature *int `json:"temperature,omitempty"` // TODO: custom type for marshalling Kelvins<->API values
}

var _ fmt.Stringer = (*Light)(nil)

func (l *Light) String() string {
	const format = "Light{On: %s, Brightness: %s, Temperature: %s}\n"
	return fmt.Sprintf(format, toString(l.On), toString(l.Brightness), toString(l.Temperature))
}

func toString(x *int) string {
	if x == nil {
		return "<nil>"
	}
	return strconv.Itoa(*x)
}

type Client struct {
	addr net.IP
}

func NewClient(addr net.IP) *Client {
	return &Client{addr: addr}
}

func (c *Client) GetLight() (*Light, error) {
	resp, err := http.Get(c.url())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r Request
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	if len(r.Lights) == 0 {
		return nil, errors.New("malformed response")
	}

	return r.Lights[0], nil
}

func (c *Client) UpdateLight(l *Light) error {
	b, err := json.Marshal(Request{Lights: []*Light{l}})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, c.url(), bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %s", resp.Status)
	}

	return nil
}

func (c *Client) url() string {
	return fmt.Sprintf("http://%s:9123/elgato/lights", c.addr.String())
}
