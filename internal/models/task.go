package models

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
}

type TaskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Overdue     bool   `json:"overdue"`
	Completed   bool   `json:"completed"`
}

type CompleteRequest struct {
	Completed bool `json:"completed"`
}
