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
		UserID   string    `json:"userId"`
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Deadline time.Time `json:"deadline"`
	}

	UpdateInput struct {
		ID       string    `json:"id"`
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Deadline time.Time `json:"deadline"`
	}

	SortBy uint

	GetAllInput struct {
		PageSize          int    `json:"pageSize"`
		Page              int    `json:"page"`
		ShowOnlyCompleted bool   `json:"showOnlyCompleted"`
		SortBy            SortBy `json:"sortBy"`
	}
)
