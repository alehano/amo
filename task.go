package amo

import "errors"

type Task struct {
	ElementID         int    `json:"element_id,omitempty"`
	ElementType       int    `json:"element_type,omitempty"`
	CompleteTillAt    int    `json:"complete_till_at,omitempty"`
	TaskType          int    `json:"task_type,omitempty"`
	Text              string `json:"text,omitempty"`
	CreatedAt         int    `json:"created_at,omitempty"`
	UpdatedAt         int    `json:"updated_at,omitempty"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty"`
	IsCompleted       bool   `json:"is_completed,omitempty"`
	CreatedBy         int    `json:"created_by,omitempty"`
}

type TaskAction struct {
	Add []Task `json:"add,omitempty"`
}


func (c *Client) AddTask(task Task) (int, error) {
	if task.Text == "" {
		return 0, errors.New("Text is empty")
	}
	if task.ResponsibleUserID == 0 {
		return 0, errors.New("ResponsibleUserID is empty")
	}
	if task.CompleteTillAt == 0 {
		return 0, errors.New("CompleteTillAt is empty")
	}
	url := c.SetURL("/api/v2/tasks", nil)
	return c.DoPostWithReturnID(url, TaskAction{Add: []Task{task}})
}
