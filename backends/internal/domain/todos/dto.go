package todos

const (
	SortByCreationASC SortBy = iota
	SortByCreationDESC
	SortByDeadlineASC
	SortByDeadlineDESC
)

type (
	SortBy uint

	GetAllInput struct {
		PageSize          int    `json:"pageSize"`
		Page              int    `json:"page"`
		ShowOnlyCompleted bool   `json:"showOnlyCompleted"`
		SortBy            SortBy `json:"sortBy"`
	}
)
