package resthttp_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

func TestUsersSignUpSignIn(t *testing.T) {
	type Response struct {
		Errors []string           `json:"errors"`
		Data   users.SignInOutput `json:"data"`
	}

	var req bytes.Buffer

	req.Write([]byte(`
		{
			"email":"john@gmail.com",
			"password":supapassword",
			"username":"John Doe"
		}
	`))
	resp, err := http.Post(apiUrl+"/users/auth/signup", "application/json", &req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error("Incorrect status code", resp.StatusCode)
	}
	res := []byte("")
	n, err := resp.Body.Read(res)
	if err != nil {
		t.Error(err)
	}
	if n == 0 {
		t.Error("Empty response")
	}
}
