package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"problem_2/internal"
	"strings"
	"testing"
)

func TestSearchHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/search?term=rick&limit=5", nil)
	w := httptest.NewRecorder()

	internal.SearchHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var results []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected some results, got 0")
	}

	found := false
	for _, r := range results {
		if strings.Contains(strings.ToLower(r.Name), "rick") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected at least one result with 'rick'")
	}
}