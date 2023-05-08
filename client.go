package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/schema"
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
	on := "<nil>"
	if l.On != nil {
		on = strconv.Itoa(*l.On)
	}

	b := "<nil>"
	if l.Brightness != nil {
		b = strconv.Itoa(*l.Brightness)
	}

	t := "<nil>"
	if l.Brightness != nil {
		t = strconv.Itoa(APIToKelvin(*l.Temperature))
	}

	const format = "On: %s\nBrightness: %s%%\nTemperature: %sK"
	return fmt.Sprintf(format, on, b, t)
}

type Client struct {
	addr   net.IP
	client *http.Client
}

var _ http.Handler = (*Client)(nil)

func NewClient(addr net.IP) *Client {
	return &Client{
		addr: addr,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetLight() (*Light, error) {
	resp, err := c.client.Get(c.url())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseBody(resp.Body)
}

func (c *Client) UpdateLight(l *Light) (*Light, error) {
	b, err := json.Marshal(Request{Lights: []*Light{l}})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.url(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}
	defer resp.Body.Close()

	return parseBody(resp.Body)
}

func (c *Client) GetSettings() (string, error) {
	resp, err := c.client.Get(c.url() + "/settings")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// TODO: refactor to use single struct for all settings
func (c *Client) SetBypass(bypass int) error {
	settings := &struct {
		Battery struct {
			Bypass int `json:"bypass"`
		} `json:"battery"`
	}{}
	settings.Battery.Bypass = bypass

	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, c.url()+"/settings", bytes.NewReader(b))
	if err != nil {
		return err
	}

	if _, err := c.client.Do(req); err != nil {
		return err
	}

	return nil
}

func (c *Client) ToggleBypass() error {
	resp, err := c.client.Get(c.url() + "/settings")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	settings := &struct {
		Battery struct {
			Bypass int `json:"bypass"`
		} `json:"battery"`
	}{}
	if err := json.Unmarshal(b, settings); err != nil {
		return err
	}

	settings.Battery.Bypass = 1 - settings.Battery.Bypass
	b, err = json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, c.url()+"/settings", bytes.NewReader(b))
	if err != nil {
		return err
	}

	if _, err := c.client.Do(req); err != nil {
		return err
	}

	return nil
}

//go:embed index.html
var website string

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println("ERROR", err)
			return
		}

		d := schema.NewDecoder()
		d.SetAliasTag("json")

		var l Light
		if err := d.Decode(&l, r.Form); err != nil {
			log.Println("ERROR", err)
			return
		}

		temperature, err := KelvinToAPI(*l.Temperature)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
		l.Temperature = &temperature

		if _, err := c.UpdateLight(&l); err != nil {
			log.Println("ERROR", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	l, err := c.GetLight()
	if err != nil {
		log.Println("ERROR", err)
		err := fmt.Errorf("couldn't fetch light settings: %w", err)
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"Enabled":     func() bool { return *l.On == 1 },
		"APIToKelvin": APIToKelvin,
	}).Parse(website)
	if err != nil {
		log.Println("ERROR", err)
		return
	}

	if err := t.Execute(w, l); err != nil {
		log.Println("ERROR", err)
		return
	}
}

func (c *Client) url() string {
	return fmt.Sprintf("http://%s:9123/elgato/lights", c.addr.String())
}

func parseBody(r io.Reader) (*Light, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req Request
	if err := json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	if len(req.Lights) == 0 {
		return nil, errors.New("malformed response")
	}

	return req.Lights[0], nil
}
