package todos

import "time"

const (
	SortByCreationASC SortBy = iota
	SortByCreationDESC
	SortByDeadlineASC
	SortByDeadlineDESC
)

type (
	CreateInput struct {
		UserID string `json:"userId" validate:"required"`
		Title  string `json:"title" validate:"gt=6,lt=100"`
		Body   string `json:"body" validate:"lt=2000"`
		// TODO: dk if i should allow deadlines in past
		Deadline time.Time `json:"deadline"`
	}

	UpdateInput struct {
		ID       string    `json:"id" validate:"required"`
		Title    string    `json:"title" validate:"gt=6,lt=100"`
		Body     string    `json:"body" validate:"lt=2000"`
		Deadline time.Time `json:"deadline"`
	}

	SortBy uint

	GetAllInput struct {
		UserID            string `json:"userID"`
		PageSize          int    `json:"pageSize"`
		Page              int    `json:"page"`
		ShowOnlyCompleted bool   `json:"showOnlyCompleted"`
		SortBy            SortBy `json:"sortBy"`
	}
)
