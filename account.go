package amo

// Account an information on the authorized account.
type Account struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Subdomain      string `json:"subdomain"`
	Currency       string `json:"currency"`
	Timezone       string `json:"timezone"`
	TimezoneOffset string `json:"timezoneoffset"`
	Language       string `json:"language"`
	CurrentUser    int    `json:"current_user"`
	Embedded struct {
		Users map[string]User
		CustomFields struct {
			Contacts  map[string]AccountCustomField `json:"contacts"`
			Leads     map[string]AccountCustomField `json:"leads"`
			Companies map[string]AccountCustomField `json:"companies"`
		} `json:"custom_fields"`
	} `json:"_embedded"`
}

func (c *Client) GetAccount(freeUsers bool, with []string) (Account, error) {
	addValues := map[string]string{}
	addValues["type"] = "json"
	if len(with) > 0 {
		addValues["with"] = MultiValuesString(with...)
	}
	if freeUsers {
		addValues["free_users"] = "Y"
	}
	url := c.SetURL("/api/v2/account", addValues)
	acc := Account{}
	err := c.DoGet(url, &acc)
	return acc, err
}
