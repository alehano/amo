package amo

import "errors"

type Note struct {
	ElementID         int    `json:"element_id,omitempty"`
	ElementType       int    `json:"element_type,omitempty"`
	Text              string `json:"text,omitempty"`
	NoteType          int    `json:"note_type,omitempty"`
	CreatedAt         int    `json:"created_at,omitempty"`
	UpdatedAt         int    `json:"updated_at,omitempty"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty"`
}

type NoteAction struct {
	Add []Note `json:"add,omitempty"`
}

func (c *Client) AddNote(note Note) (int, error) {
	if note.Text == "" {
		return 0, errors.New("Text is empty")
	}
	if note.ElementID == 0 {
		return 0, errors.New("ElementID is empty")
	}
	if note.ElementType == 0 {
		return 0, errors.New("ElementType is empty")
	}
	if note.NoteType == 0 {
		return 0, errors.New("NoteType is empty")
	}
	url := c.SetURL("/api/v2/notes", nil)
	return c.DoPostWithReturnID(url, NoteAction{Add: []Note{note}})
}
