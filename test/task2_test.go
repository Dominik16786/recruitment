package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"problem_2/internal"
	"testing"
)

func TestTopPairsHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/top-pairs?min=10&max=50&limit=5", nil)
	w := httptest.NewRecorder()

	internal.TopPairsHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var results []struct {
		Character1 struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"character1"`
		Character2 struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"character2"`
		Episodes int `json:"episodes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected some pairs, got 0")
	}

	for _, r := range results {
		if r.Episodes < 10 || r.Episodes > 50 {
			t.Fatalf("episodes out of range: %d", r.Episodes)
		}
	}
}