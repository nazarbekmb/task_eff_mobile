package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type EnricherService struct {
	Client *http.Client
}

func NewEnricherService() *EnricherService {
	return &EnricherService{
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

type ageResponse struct {
	Age int `json:"age"`
}

type genderResponse struct {
	Gender string `json:"gender"`
}

type nationalizeResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func (s *EnricherService) Enrich(name string) (int, string, []string, error) {
	age, err := s.fetchAge(name)
	if err != nil {
		return 0, "", nil, err
	}
	gender, err := s.fetchGender(name)
	if err != nil {
		return 0, "", nil, err
	}
	nations, err := s.fetchNationalities(name)
	if err != nil {
		return 0, "", nil, err
	}
	return age, gender, nations, nil
}

func (s *EnricherService) fetchAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	log.WithField("url", url).Debug("Executing a request to AgifyAPI")

	resp, err := s.Client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data ageResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"name": name,
		"age":  data.Age,
	}).Info("Enrichment: age received")
	return data.Age, nil
}

func (s *EnricherService) fetchGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	log.WithField("url", url).Debug("Executing a request to GenderizeAPI")

	resp, err := s.Client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data genderResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"name":   name,
		"gender": data.Gender,
	}).Info("Enrichment: gender received")
	return data.Gender, nil
}

func (s *EnricherService) fetchNationalities(name string) ([]string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	log.WithField("url", url).Debug("Execute a request to NationalizeAPI")

	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data nationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var countries []string
	for _, c := range data.Country {
		countries = append(countries, c.CountryID)
		if len(countries) >= 2 {
			break
		}
	}
	log.WithFields(log.Fields{
		"name":          name,
		"nationalities": data.Country,
	}).Info("Enrichment: nationality received")
	return countries, nil
}
