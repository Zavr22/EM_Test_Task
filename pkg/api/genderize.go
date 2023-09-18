package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GenderizeResponse struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Count       int     `json:"count"`
	Probability float64 `json:"probability"`
}

func GetGenderizeGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var genderizeResponse GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderizeResponse); err != nil {
		return "", err
	}

	if genderizeResponse.Count == 0 {
		return "", fmt.Errorf("Genderize API could not determine gender")
	}

	return genderizeResponse.Gender, nil
}
