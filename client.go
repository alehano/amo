package amo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"bytes"
	"errors"
)

// Client AmoCRM API Client
type Client struct {
	userLogin         string
	userHash          string
	Timezone          string
	accountWebAddress *url.URL
	rateLimiter       RateLimiter
	timeout           time.Duration
}

// NewClient creates and initializes AmoCRM API.
// accountWebAddress is your AmoCRM account address like https://your-account.amocrm.com
func NewClient(accountWebAddress string, rateLimiter RateLimiter) (*Client, error) {
	c := &Client{rateLimiter: rateLimiter, timeout: 5 * time.Second}
	if c.rateLimiter == nil {
		c.rateLimiter = defaultRTLimiter
	}
	var err error
	c.accountWebAddress, err = url.Parse(accountWebAddress)
	return c, err
}

// AuthResponse information about Amo' user and account.
type AuthResponse struct {
	Response struct {
		Auth bool `json:"auth"`
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

// Authorize In case of successful authorization,
// this method returns the user’s session ID,
// which must be used while accessing to other API methods.
// Saves user's login and hash into the client for next usage.
func (c *Client) Authorize(userLogin, userHash string) (*AuthResponse, error) {
	c.userHash = userHash
	c.userLogin = userLogin
	reqURL := *c.accountWebAddress
	reqURL.Path = "/private/api/auth.php"
	values := reqURL.Query()
	values.Set("type", "json")
	reqURL.RawQuery = values.Encode()
	c.rateLimiter.WaitForRequest()
	client := &http.Client{
		Timeout: c.timeout,
	}
	resp, err := client.PostForm(reqURL.String(),
		url.Values{"USER_LOGIN": {c.userLogin}, "USER_HASH": {c.userHash}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, err
	}
	if len(authResponse.Response.Accounts) > 0 {
		c.Timezone = authResponse.Response.Accounts[0].Timezone
	}
	return &authResponse, nil
}

// Set URL and add default params
func (c *Client) SetURL(url string, addValues map[string]string) string {
	reqURL := *c.accountWebAddress
	reqURL.Path = url
	values := reqURL.Query()
	values.Set("USER_LOGIN", c.userLogin)
	values.Set("USER_HASH", c.userHash)
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
	resp, err := http.Get(url)
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
	return http.Post(url, "application/json", buf)
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
