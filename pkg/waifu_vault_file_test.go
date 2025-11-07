package waifuVault

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
)

func TestUploadFile(t *testing.T) {
	ctx := context.Background()

	t.Run("should upload a file as bytes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Expected PUT request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuResponseMock2)
		}))
		defer server.Close()

		// Override baseUrl for testing
		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		fileBytes := []byte("test content")
		result, err := api.UploadFile(ctx, mod.WaifuvaultPutOpts{
			Bytes:    &fileBytes,
			FileName: "test.txt",
		})

		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}
		if result.Token != WaifuResponseMock2.Token {
			t.Errorf("Expected token %s, got %s", WaifuResponseMock2.Token, result.Token)
		}
	})

	t.Run("should upload a file from URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Expected PUT request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuResponseMock2)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.UploadFile(ctx, mod.WaifuvaultPutOpts{
			Url: "https://example.com/file.txt",
		})

		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}
		if result.Token != WaifuResponseMock2.Token {
			t.Errorf("Expected token %s, got %s", WaifuResponseMock2.Token, result.Token)
		}
	})

	t.Run("should upload with options", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("Expected PUT request, got %s", r.Method)
			}
			// Check query parameters
			query := r.URL.Query()
			if query.Get("expires") != "2d" {
				t.Errorf("Expected expires=2d, got %s", query.Get("expires"))
			}
			if query.Get("hide_filename") != "true" {
				t.Errorf("Expected hide_filename=true, got %s", query.Get("hide_filename"))
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuResponseMock2)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.UploadFile(ctx, mod.WaifuvaultPutOpts{
			Url:          "https://example.com/file.txt",
			Password:     "foo",
			HideFilename: true,
			Expires:      "2d",
		})

		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}
		if result.Token != WaifuResponseMock2.Token {
			t.Errorf("Expected token %s, got %s", WaifuResponseMock2.Token, result.Token)
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
		_, err := api.UploadFile(ctx, mod.WaifuvaultPutOpts{
			Url: "https://example.com/file.txt",
		})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
		expectedErr := "Error 400 (whore): loser"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("Expected error containing %s, got %s", expectedErr, err.Error())
		}
	})
}

func TestFileInfo(t *testing.T) {
	ctx := context.Background()

	t.Run("should get file info from token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuResponseMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.FileInfo(ctx, WaifuResponseMock1.Token)

		if err != nil {
			t.Fatalf("FileInfo failed: %v", err)
		}
		if result.Token != WaifuResponseMock1.Token {
			t.Errorf("Expected token %s, got %s", WaifuResponseMock1.Token, result.Token)
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
		_, err := api.FileInfo(ctx, WaifuResponseMock1.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestFileInfoFormatted(t *testing.T) {
	ctx := context.Background()

	t.Run("should get formatted file info", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("formatted") != "true" {
				t.Errorf("Expected formatted=true query parameter")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuResponseMock2)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.FileInfoFormatted(ctx, WaifuResponseMock2.Token)

		if err != nil {
			t.Fatalf("FileInfoFormatted failed: %v", err)
		}
		if result.RetentionPeriod != WaifuResponseMock2.RetentionPeriod {
			t.Errorf("Expected retention period %s, got %s", WaifuResponseMock2.RetentionPeriod, result.RetentionPeriod)
		}
	})
}

func TestDeleteFile(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete a file", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Expected DELETE request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("true"))
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.DeleteFile(ctx, WaifuResponseMock1.Token)

		if err != nil {
			t.Fatalf("DeleteFile failed: %v", err)
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
		_, err := api.DeleteFile(ctx, WaifuResponseMock1.Token)

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestGetFile(t *testing.T) {
	ctx := context.Background()

	t.Run("should get a file from token", func(t *testing.T) {
		fileContent := []byte("test file content")
		callCount := 0

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount == 1 {
				// First call to get file info
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(WaifuResponseMock1)
			} else {
				// Second call to download file
				w.WriteHeader(http.StatusOK)
				w.Write(fileContent)
			}
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		// Update mock URL to point to our test server
		WaifuResponseMock1.URL = server.URL + "/test.txt"

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.GetFile(ctx, mod.GetFileInfo{
			Token: WaifuResponseMock1.Token,
		})

		if err != nil {
			t.Fatalf("GetFile failed: %v", err)
		}
		if !bytes.Equal(result, fileContent) {
			t.Errorf("Expected file content %s, got %s", string(fileContent), string(result))
		}
	})

	t.Run("should get a file with password", func(t *testing.T) {
		fileContent := []byte("protected file")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			password := r.Header.Get("x-password")
			if password != "test-password" {
				t.Errorf("Expected x-password header to be test-password, got %s", password)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(fileContent)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		result, err := api.GetFile(ctx, mod.GetFileInfo{
			Filename: "1710111505084/08.png",
			Password: "test-password",
		})

		if err != nil {
			t.Fatalf("GetFile failed: %v", err)
		}
		if !bytes.Equal(result, fileContent) {
			t.Errorf("Expected file content %s, got %s", string(fileContent), string(result))
		}
	})

	t.Run("should handle incorrect password", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<div></div>"))
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		_, err := api.GetFile(ctx, mod.GetFileInfo{
			Filename: "1710111505084/08.png",
			Password: "wrong-password",
		})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if !strings.Contains(err.Error(), "password is incorrect") {
			t.Errorf("Expected password error, got: %s", err.Error())
		}
	})
}

func TestModifyFile(t *testing.T) {
	ctx := context.Background()

	t.Run("should modify a file", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				t.Errorf("Expected PATCH request, got %s", r.Method)
			}

			// Read and verify request body
			body, _ := io.ReadAll(r.Body)
			var payload mod.ModifyEntryPayload
			json.Unmarshal(body, &payload)

			if payload.CustomExpiry == nil || *payload.CustomExpiry != "2d" {
				t.Error("Expected customExpiry to be 2d")
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(WaifuResponseMock1)
		}))
		defer server.Close()

		origBaseUrl := baseUrl
		defer func() { baseUrl = origBaseUrl }()
		baseUrl = server.URL

		api := NewWaifuvaltApi(http.Client{})
		customExpiry := "2d"
		result, err := api.ModifyFile(ctx, WaifuResponseMock1.Token, mod.ModifyEntryPayload{
			CustomExpiry: &customExpiry,
		})

		if err != nil {
			t.Fatalf("ModifyFile failed: %v", err)
		}
		if result.Token != WaifuResponseMock1.Token {
			t.Errorf("Expected token %s, got %s", WaifuResponseMock1.Token, result.Token)
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
		customExpiry := "2d"
		_, err := api.ModifyFile(ctx, WaifuResponseMock1.Token, mod.ModifyEntryPayload{
			CustomExpiry: &customExpiry,
		})

		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}
