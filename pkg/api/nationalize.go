package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NationalizeResponse struct {
	Name    string `json:"name"`
	Country []struct {
		CountryId   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func GetNationalizeNationality(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var nationalizeResponse NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationalizeResponse); err != nil {
		return "", err
	}

	if len(nationalizeResponse.Country) == 0 {
		return "", fmt.Errorf("Nationalize API could not determine nationality")
	}

	maxProbability := 0.0
	nationality := ""
	for _, country := range nationalizeResponse.Country {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			nationality = country.CountryId
		}
	}

	return nationality, nil
}
