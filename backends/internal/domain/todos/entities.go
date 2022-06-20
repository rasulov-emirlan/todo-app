package todos

import (
	"time"

	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

type (
	Todo struct {
		ID     string      `json:"id"`
		Author *users.User `json:"author,omitempty"`

		Title string `json:"title"`
		Body  string `json:"body"`

		Completed bool      `json:"completed"`
		Deadline  time.Time `json:"deadline"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
