package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBug1IdenticalDiffKeepsEmptyHunksSlice(t *testing.T) {
	srv := httptest.NewServer(New().Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/diff", "application/json", bytes.NewBufferString(`{"old":"same\n","new":"same\n"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, data)
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatal(err)
	}
	if got := string(body["hunks"]); got != "[]" {
		t.Fatalf("hunks = %s, want []", got)
	}
}
