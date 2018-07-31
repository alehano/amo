package amo

type User struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	LastName string            `json:"last_name"`
	Login    string            `json:"login"`
	Language string            `json:"language"`
	GroupID  int               `json:"group_id"`
	IsActive bool              `json:"is_active"`
	IsFree   bool              `json:"is_free"`
	IsAdmin  bool              `json:"is_admin"`
	Rights   map[string]string `json:"rights,omitempty"`
}
