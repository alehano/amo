package amo

import "errors"

type Contact struct {
	Name              string        `json:"name,omitempty"`
	CreatedAt         int           `json:"created_at,omitempty"`
	UpdatedAt         int           `json:"updated_at,omitempty"`
	ResponsibleUserID int           `json:"responsible_user_id,omitempty"`
	CreatedBy         int           `json:"created_by,omitempty"`
	CompanyName       string        `json:"company_name,omitempty"`
	Tags              string        `json:"tags,omitempty"`
	LeadsID           string        `json:"leads_id,omitempty"`
	CustomersID       string        `json:"customers_id,omitempty"`
	CompanyID         string        `json:"company_id,omitempty"`
	CustomFields      []CustomField `json:"custom_fields,omitempty"`
}

type ContactAction struct {
	Add []Contact `json:"add,omitempty"`
}

func (c *Client) AddContact(contact Contact) (int, error) {
	if contact.Name == "" {
		return 0, errors.New("Name is empty")
	}
	url := c.SetURL("/api/v2/contacts", nil)
	return c.DoPostWithReturnID(url, ContactAction{Add: []Contact{contact}})
}
