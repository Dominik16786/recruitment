package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type SearchResult struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type APIResponse struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func fetchResults(endpoint, term, resultType string) ([]SearchResult, error) {
	url := fmt.Sprintf("https://rickandmortyapi.com/api/%s/?name=%s", endpoint, term)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, _ := io.ReadAll(resp.Body)

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range apiResp.Results {
		results = append(results, SearchResult{
			Name: r.Name,
			Type: resultType,
			URL:  r.URL,
		})
	}

	return results, nil
}

// Handler HTTP dla /search
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	limitParam := r.URL.Query().Get("limit")

	if term == "" {
		http.Error(w, "Missing term parameter", http.StatusBadRequest)
		return
	}

	var limit int
	var err error
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	var allResults []SearchResult

	characters, _ := fetchResults("character", term, "character")
	locations, _ := fetchResults("location", term, "location")
	episodes, _ := fetchResults("episode", term, "episode")

	allResults = append(allResults, characters...)
	allResults = append(allResults, locations...)
	allResults = append(allResults, episodes...)

	if limit > 0 && limit < len(allResults) {
		allResults = allResults[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allResults)
}