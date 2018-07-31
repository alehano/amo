package amo

import "errors"

type Company struct {
	Name              string        `json:"name,omitempty"`
	CreatedAt         int           `json:"created_at,omitempty"`
	UpdatedAt         int           `json:"updated_at,omitempty"`
	ResponsibleUserID int           `json:"responsible_user_id,omitempty"`
	CreatedBy         int           `json:"created_by,omitempty"`
	Tags              string        `json:"tags,omitempty"`
	LeadsID           string        `json:"leads_id,omitempty"`
	CustomersID       string        `json:"customers_id,omitempty"`
	ContactsID        string        `json:"contacts_id,omitempty"`
	CustomFields      []CustomField `json:"custom_fields,omitempty"`
}

type CompanyAction struct {
	Add []Company `json:"add,omitempty"`
}

func (c *Client) AddCompany(comp Company) (int, error) {
	if comp.Name == "" {
		return 0, errors.New("Name is empty")
	}
	url := c.SetURL("/api/v2/companies", nil)
	return c.DoPostWithReturnID(url, CompanyAction{Add: []Company{comp}})
}
