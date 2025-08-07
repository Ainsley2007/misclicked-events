package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHiscoreDataSource_CheckIfPlayerExists(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("player") == "testplayer" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create datasource with test server URL
	ds := &hiscoreDS{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	tests := []struct {
		name     string
		username string
		want     bool
		wantErr  bool
	}{
		{
			name:     "player exists",
			username: "testplayer",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "player does not exist",
			username: "nonexistentplayer",
			want:     false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.CheckIfPlayerExists(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIfPlayerExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckIfPlayerExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHiscoreDataSource_FetchHiscore(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("player") == "testplayer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"skills": [
					{"id": 1, "name": "Overall", "rank": 1, "level": 99, "xp": 13034431}
				],
				"activities": [
					{"id": 1, "name": "Clue Scrolls (all)", "rank": 1, "score": 100}
				]
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create datasource with test server URL
	ds := &hiscoreDS{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid player",
			username: "testplayer",
			wantErr:  false,
		},
		{
			name:     "invalid player",
			username: "nonexistentplayer",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ds.FetchHiscore(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchHiscore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && data == nil {
				t.Errorf("FetchHiscore() returned nil data for valid player")
			}
			if !tt.wantErr && len(data.Skills) == 0 {
				t.Errorf("FetchHiscore() returned no skills")
			}
			if !tt.wantErr && len(data.Activities) == 0 {
				t.Errorf("FetchHiscore() returned no activities")
			}
		})
	}
}

func TestNewHiscoreDataSource(t *testing.T) {
	ds := NewHiscoreDataSource()
	if ds == nil {
		t.Error("NewHiscoreDataSource() returned nil")
	}
}
