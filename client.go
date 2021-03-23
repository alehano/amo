package amo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"time"
)

// Client AmoCRM API Client
type Client struct {
	token             *oauth2.Token
	config            *oauth2.Config
	Timezone          string
	accountWebAddress *url.URL
	rateLimiter       RateLimiter
	timeout           time.Duration
}

// NewClient creates and initializes AmoCRM API.
// accountWebAddress is your AmoCRM account address like https://your-account.amocrm.com
func NewClient(accountWebAddress string, token *oauth2.Token, config *oauth2.Config, rateLimiter RateLimiter) (*Client, error) {
	c := &Client{rateLimiter: rateLimiter, timeout: 5 * time.Second}
	if c.rateLimiter == nil {
		c.rateLimiter = defaultRTLimiter
	}
	var err error
	c.accountWebAddress, err = url.Parse(accountWebAddress)
	c.token = token
	c.config = config
	return c, err
}

// AuthResponse information about Amo' user and account.
type AuthResponse struct {
	Response struct {
		Auth     bool `json:"auth"`
		Accounts []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Subdomain string `json:"subdomain"`
			Language  string `json:"language"`
			Timezone  string `json:"timezone"`
		} `json:"accounts"`
		ServerTime int `json:"server_time"`
	} `json:"response"`
}

func (c *Client) GetToken() *oauth2.Token {
	return c.token
}

// Set URL and add default params
func (c *Client) SetURL(url string, addValues map[string]string) string {
	reqURL := *c.accountWebAddress
	reqURL.Path = url
	values := reqURL.Query()
	if addValues != nil {
		for key, value := range addValues {
			values.Set(key, value)
		}
	}
	reqURL.RawQuery = values.Encode()
	return reqURL.String()
}

func (c *Client) DoGet(url string, result interface{}) error {
	c.rateLimiter.WaitForRequest()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization: Bearer", c.token.AccessToken)
	client := c.config.Client(context.TODO(), c.token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(&result)
}

func (c *Client) DoPost(url string, data interface{}) (*http.Response, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	c.rateLimiter.WaitForRequest()
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization: Bearer", c.token.AccessToken)
	client := c.config.Client(context.TODO(), c.token)
	return client.Do(req)
}

// Post with return ID on new entity
func (c *Client) DoPostWithReturnID(url string, data interface{}) (int, error) {
	type respID struct {
		Embedded struct {
			Items []struct {
				ID int `json:"id"`
			} `json:"items"`
		} `json:"_embedded"`
		Response struct {
			Error string `json:"error"`
		} `json:"response"`
	}
	resp, err := c.DoPost(url, data)
	if err != nil {
		return 0, err
	}
	result := respID{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return 0, err
	}
	if len(result.Embedded.Items) == 0 {
		if result.Response.Error != "" {
			return 0, errors.New(result.Response.Error)
		}
		return 0, errors.New("No Items")
	}
	return result.Embedded.Items[0].ID, nil
}
