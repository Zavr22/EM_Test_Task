package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetAgifyAge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/?name=John" {
			t.Errorf("Expected URL path to be '/?name=John', got '%s'", r.URL.Path)
		}
		response := `{"name": "John", "age": 30, "count": 1, "errors": false}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer server.Close()

	var agifyAPIURL string
	oldURL := agifyAPIURL
	agifyAPIURL = server.URL
	defer func() {
		agifyAPIURL = oldURL
	}()
	age, err := GetAgifyAge("John")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	var intT int
	if reflect.TypeOf(age) != reflect.TypeOf(intT) {
		t.Errorf("Expected JSON response to have age of type int, got %T", age)
	}
}
