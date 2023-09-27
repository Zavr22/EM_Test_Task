package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetNationalizeNationality(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/?name=John" {
			t.Errorf("Expected URL path to be '/?name=John', got '%s'", r.URL.Path)
		}
		response := `{"name": "John", "nationality": "by", "count": 1, "errors": false}`
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
	nationality, err := GetNationalizeNationality("John")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	var strT string
	if reflect.TypeOf(nationality) != reflect.TypeOf(strT) {
		t.Errorf("Expected JSON response to have naionality, got %T", nationality)
	}
}
