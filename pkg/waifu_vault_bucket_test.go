package waifuVault

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateBucket(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a new bucket", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET request, got %s", r.Method)
			}
			if !strings.HasSuffix(r.URL.Path, "/rest/bucket/create") {
				t.Errorf("Expected path to end with /rest/bucket/create, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuBucketMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.CreateBucket(ctx)

		if err != nil {
			t.Fatalf("CreateBucket failed: %v", err)
		}
		if result.Token != WaifuBucketMock1.Token {
			t.Errorf("Expected token %s, got %s", WaifuBucketMock1.Token, result.Token)
		}
		if len(result.Files) != len(WaifuBucketMock1.Files) {
			t.Errorf("Expected %d files, got %d", len(WaifuBucketMock1.Files), len(result.Files))
		}
	})

	t.Run("should handle error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(WaifuErrorMock)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		_, err := api.CreateBucket(ctx)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestGetBucket(t *testing.T) {
	ctx := context.Background()

	t.Run("should get a bucket", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if !strings.HasSuffix(r.URL.Path, "/rest/bucket/get") {
				t.Errorf("Expected path to end with /rest/bucket/get, got %s", r.URL.Path)
			}

			// Verify request body
			body, _ := io.ReadAll(r.Body)
			var payload struct {
				BucketToken string `json:"bucket_token"`
			}
			json.Unmarshal(body, &payload)

			if payload.BucketToken != WaifuBucketMock1.Token {
				t.Errorf("Expected bucket_token %s, got %s", WaifuBucketMock1.Token, payload.BucketToken)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuBucketMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.GetBucket(ctx, WaifuBucketMock1.Token)

		if err != nil {
			t.Fatalf("GetBucket failed: %v", err)
		}
		if result.Token != WaifuBucketMock1.Token {
			t.Errorf("Expected token %s, got %s", WaifuBucketMock1.Token, result.Token)
		}
	})

	t.Run("should handle error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(WaifuErrorMock)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		_, err := api.GetBucket(ctx, WaifuBucketMock1.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestDeleteBucket(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete a bucket", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Expected DELETE request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, WaifuBucketMock1.Token) {
				t.Errorf("Expected path to contain bucket token %s", WaifuBucketMock1.Token)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("true"))
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DeleteBucket(ctx, WaifuBucketMock1.Token)

		if err != nil {
			t.Fatalf("DeleteBucket failed: %v", err)
		}
		if !result {
			t.Error("Expected true but got false")
		}
	})

	t.Run("should handle error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(WaifuErrorMock)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		_, err := api.DeleteBucket(ctx, WaifuBucketMock1.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}
