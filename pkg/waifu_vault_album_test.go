package waifuVault

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
)

func TestCreateAlbum(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a new album in an existing bucket with files", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, WaifuBucketMock1.Token) {
				t.Errorf("Expected path to contain bucket token %s", WaifuBucketMock1.Token)
			}

			// Verify request body
			body, _ := io.ReadAll(r.Body)
			var payload mod.WaifuAlbumCreateBody
			json.Unmarshal(body, &payload)

			if payload.Name != WaifuAlbumMock1.Name {
				t.Errorf("Expected name %s, got %s", WaifuAlbumMock1.Name, payload.Name)
			}
			if payload.BucketToken != WaifuBucketMock1.Token {
				t.Errorf("Expected bucketToken %s, got %s", WaifuBucketMock1.Token, payload.BucketToken)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuAlbumMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.CreateAlbum(ctx, mod.WaifuAlbumCreateBody{
			Name:        WaifuAlbumMock1.Name,
			BucketToken: WaifuBucketMock1.Token,
		})

		if err != nil {
			t.Fatalf("CreateAlbum failed: %v", err)
		}
		if result.Token != WaifuAlbumMock1.Token {
			t.Errorf("Expected token %s, got %s", WaifuAlbumMock1.Token, result.Token)
		}
		if result.Name != WaifuAlbumMock1.Name {
			t.Errorf("Expected name %s, got %s", WaifuAlbumMock1.Name, result.Name)
		}
	})

	t.Run("should create a new album in an existing bucket without files", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuAlbumMock2)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.CreateAlbum(ctx, mod.WaifuAlbumCreateBody{
			Name:        WaifuAlbumMock2.Name,
			BucketToken: WaifuBucketMock1.Token,
		})

		if err != nil {
			t.Fatalf("CreateAlbum failed: %v", err)
		}
		if len(result.Files) != 0 {
			t.Errorf("Expected 0 files, got %d", len(result.Files))
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
		_, err := api.CreateAlbum(ctx, mod.WaifuAlbumCreateBody{
			Name:        WaifuAlbumMock1.Name,
			BucketToken: WaifuBucketMock1.Token,
		})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestAssociateFiles(t *testing.T) {
	ctx := context.Background()

	t.Run("should associate a file to an album", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, "/associate") {
				t.Errorf("Expected path to contain /associate")
			}

			// Verify request body
			body, _ := io.ReadAll(r.Body)
			var payload struct {
				FileTokens []string `json:"fileTokens"`
			}
			json.Unmarshal(body, &payload)

			if len(payload.FileTokens) != 1 || payload.FileTokens[0] != WaifuResponseMock1.Token {
				t.Errorf("Expected fileTokens to contain %s", WaifuResponseMock1.Token)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuAlbumMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.AssociateFiles(ctx, WaifuAlbumMock1.Token, []string{WaifuResponseMock1.Token})

		if err != nil {
			t.Fatalf("AssociateFiles failed: %v", err)
		}
		if result.Token != WaifuAlbumMock1.Token {
			t.Errorf("Expected token %s, got %s", WaifuAlbumMock1.Token, result.Token)
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
		_, err := api.AssociateFiles(ctx, WaifuAlbumMock1.Token, []string{WaifuResponseMock1.Token})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestDisassociateFiles(t *testing.T) {
	ctx := context.Background()

	t.Run("should disassociate a file from an album", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, "/disassociate") {
				t.Errorf("Expected path to contain /disassociate")
			}

			// Verify request body
			body, _ := io.ReadAll(r.Body)
			var payload struct {
				FileTokens []string `json:"fileTokens"`
			}
			json.Unmarshal(body, &payload)

			if len(payload.FileTokens) != 1 || payload.FileTokens[0] != WaifuResponseMock1.Token {
				t.Errorf("Expected fileTokens to contain %s", WaifuResponseMock1.Token)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuAlbumMock2)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DisassociateFiles(ctx, WaifuAlbumMock2.Token, []string{WaifuResponseMock1.Token})

		if err != nil {
			t.Fatalf("DisassociateFiles failed: %v", err)
		}
		if result.Token != WaifuAlbumMock2.Token {
			t.Errorf("Expected token %s, got %s", WaifuAlbumMock2.Token, result.Token)
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
		_, err := api.DisassociateFiles(ctx, WaifuAlbumMock2.Token, []string{WaifuResponseMock1.Token})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestGetAlbum(t *testing.T) {
	ctx := context.Background()

	t.Run("should get an album from a private token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuAlbumMock2)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.GetAlbum(ctx, WaifuAlbumMock2.Token)

		if err != nil {
			t.Fatalf("GetAlbum failed: %v", err)
		}
		if result.Token != WaifuAlbumMock2.Token {
			t.Errorf("Expected token %s, got %s", WaifuAlbumMock2.Token, result.Token)
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
		_, err := api.GetAlbum(ctx, WaifuAlbumMock2.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestDeleteAlbum(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete an album without deleting files", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Expected DELETE request, got %s", r.Method)
			}
			if r.URL.Query().Get("deleteFiles") != "false" {
				t.Errorf("Expected deleteFiles=false, got %s", r.URL.Query().Get("deleteFiles"))
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(GenericSuccessDeletedMock)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DeleteAlbum(ctx, WaifuAlbumMock2.Token, false)

		if err != nil {
			t.Fatalf("DeleteAlbum failed: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("should delete an album and delete files", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Expected DELETE request, got %s", r.Method)
			}
			if r.URL.Query().Get("deleteFiles") != "true" {
				t.Errorf("Expected deleteFiles=true, got %s", r.URL.Query().Get("deleteFiles"))
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(GenericSuccessDeletedMock)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DeleteAlbum(ctx, WaifuAlbumMock2.Token, true)

		if err != nil {
			t.Fatalf("DeleteAlbum failed: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
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
		_, err := api.DeleteAlbum(ctx, WaifuAlbumMock2.Token, false)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestShareAlbum(t *testing.T) {
	ctx := context.Background()

	t.Run("should share an album", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, "/share/") {
				t.Errorf("Expected path to contain /share/")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SharedFileMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.ShareAlbum(ctx, WaifuAlbumMock2.Token)

		if err != nil {
			t.Fatalf("ShareAlbum failed: %v", err)
		}
		if result != SharedFileMock1.Description {
			t.Errorf("Expected description %s, got %s", SharedFileMock1.Description, result)
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
		_, err := api.ShareAlbum(ctx, WaifuAlbumMock2.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestRevokeAlbum(t *testing.T) {
	ctx := context.Background()

	t.Run("should revoke a shared album", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, "/revoke/") {
				t.Errorf("Expected path to contain /revoke/")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SharedFileMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.RevokeAlbum(ctx, WaifuAlbumMock1.Token)

		if err != nil {
			t.Fatalf("RevokeAlbum failed: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
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
		_, err := api.RevokeAlbum(ctx, WaifuAlbumMock1.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestDownloadAlbum(t *testing.T) {
	ctx := context.Background()

	t.Run("should download a single file from an album", func(t *testing.T) {
		zipContent := []byte("fake zip content")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, "/download/") {
				t.Errorf("Expected path to contain /download/")
			}

			// Verify request body
			body, _ := io.ReadAll(r.Body)
			var payload []int
			json.Unmarshal(body, &payload)

			if len(payload) != 1 || payload[0] != WaifuResponseMock1.ID {
				t.Errorf("Expected file ID %d, got %v", WaifuResponseMock1.ID, payload)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(zipContent)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DownloadAlbum(ctx, WaifuAlbumMock1.Token, []int{WaifuResponseMock1.ID})

		if err != nil {
			t.Fatalf("DownloadAlbum failed: %v", err)
		}
		if string(result) != string(zipContent) {
			t.Errorf("Expected zip content %s, got %s", string(zipContent), string(result))
		}
	})

	t.Run("should download a whole album", func(t *testing.T) {
		zipContent := []byte("fake zip content for whole album")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request body is an empty array
			body, _ := io.ReadAll(r.Body)
			var payload []int
			json.Unmarshal(body, &payload)

			if len(payload) != 0 {
				t.Errorf("Expected empty array, got %v", payload)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(zipContent)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DownloadAlbum(ctx, WaifuAlbumMock1.Token, []int{})

		if err != nil {
			t.Fatalf("DownloadAlbum failed: %v", err)
		}
		if string(result) != string(zipContent) {
			t.Errorf("Expected zip content %s, got %s", string(zipContent), string(result))
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
		_, err := api.DownloadAlbum(ctx, WaifuAlbumMock1.Token, []int{})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}
