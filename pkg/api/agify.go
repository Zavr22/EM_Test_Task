package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AgifyResponse struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Count  int    `json:"count"`
	Errors bool   `json:"errors"`
}

func GetAgifyAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var agifyResponse AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&agifyResponse); err != nil {
		return 0, err
	}

	if agifyResponse.Errors {
		return 0, fmt.Errorf("Agify API returned an error")
	}

	return agifyResponse.Age, nil
}
