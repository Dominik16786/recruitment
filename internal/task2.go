package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

type Character struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PairResult struct {
	Character1 Character `json:"character1"`
	Character2 Character `json:"character2"`
	Episodes   int       `json:"episodes"`
}

type APICharacter struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Episode []string `json:"episode"`
	URL     string   `json:"url"`
}

type APICharacterResponse struct {
	Results []APICharacter `json:"results"`
	Info    struct {
		Next string `json:"next"`
	} `json:"info"`
}

func fetchAllCharacters() ([]APICharacter, error) {
	var all []APICharacter
	url := "https://rickandmortyapi.com/api/character"

	for url != "" {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching URL:", url, err)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Println("Non-200 status code from:", url, resp.Status)
			break
		}

		var data APICharacterResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			fmt.Println("Error decoding JSON from:", url, err)
			return nil, err
		}
		resp.Body.Close()

		all = append(all, data.Results...)
		url = data.Info.Next
	}

	return all, nil
}

func TopPairsHandler(w http.ResponseWriter, r *http.Request) {
	minParam := r.URL.Query().Get("min")
	maxParam := r.URL.Query().Get("max")
	limitParam := r.URL.Query().Get("limit")

	min := 0
	max := 1 << 30 
	limit := 20

	if minParam != "" {
		if v, err := strconv.Atoi(minParam); err == nil {
			min = v
		}
	}
	if maxParam != "" {
		if v, err := strconv.Atoi(maxParam); err == nil {
			max = v
		}
	}
	if limitParam != "" {
		if v, err := strconv.Atoi(limitParam); err == nil {
			limit = v
		}
	}

	chars, err := fetchAllCharacters()
	if err != nil {
		fmt.Println("fetchAllCharacters error:", err)
		http.Error(w, "Failed to fetch characters", http.StatusInternalServerError)
		return
	}

	type pairKey struct{ i, j int }
	pairsMap := make(map[pairKey]int)

	for i := 0; i < len(chars); i++ {
		for j := i + 1; j < len(chars); j++ {
			count := 0
			episodesI := make(map[string]bool)
			for _, ep := range chars[i].Episode {
				episodesI[ep] = true
			}
			for _, ep := range chars[j].Episode {
				if episodesI[ep] {
					count++
				}
			}
			if count >= min && count <= max {
				pairsMap[pairKey{i, j}] = count
			}
		}
	}

	var results []PairResult
	for k, v := range pairsMap {
		results = append(results, PairResult{
			Character1: Character{Name: chars[k.i].Name, URL: chars[k.i].URL},
			Character2: Character{Name: chars[k.j].Name, URL: chars[k.j].URL},
			Episodes:   v,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Episodes > results[j].Episodes
	})

	if limit < len(results) {
		results = results[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}