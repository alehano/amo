package amo

type AccountCustomField struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	FieldType   int               `json:"field_type"`
	Sort        int               `json:"sort"`
	IsMultiple  bool              `json:"is_multiple"`
	IsSystem    bool              `json:"is_system"`
	IsEditable  bool              `json:"is_editable"`
	IsRequired  bool              `json:"is_required"`
	IsDeletable bool              `json:"is_deletable"`
	IsVisible   bool              `json:"is_visible"`
	Enums       map[string]string `json:"enums"`
}

type CustomField struct {
	ID     int       `json:"id"`
	Values []CSValue `json:"values"`
}

type CSValue struct {
	Value   string `json:"value"`
	Enum    string `json:"enum"`
	Subtype string `json:"subtype"`
}
