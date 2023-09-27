package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetGenderizeGender(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/?name=John" {
			t.Errorf("Expected URL path to be '/?name=John', got '%s'", r.URL.Path)
		}
		response := `{"name": "John", "gender": "male", "count": 1, "errors": false}`
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
	gender, err := GetGenderizeGender("John")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	var strT string
	if reflect.TypeOf(gender) != reflect.TypeOf(strT) {
		t.Errorf("Expected JSON response to have gender type string, got %T", gender)
	}
}
