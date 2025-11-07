package waifuVault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
)

func (re *api) CreateAlbum(ctx context.Context, body mod.WaifuAlbumCreateBody) (*mod.WaifuAlbum, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/%s", body.BucketToken))
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r, err := re.createRequest(ctx, http.MethodPost, albumUrl, bytes.NewBuffer(jsonData), nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, mod.WaifuAlbum{})
}

func (re *api) AssociateFiles(ctx context.Context, albumToken string, filesToAssociate []string) (*mod.WaifuAlbum, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/%s/associate", albumToken))
	type payload struct {
		FileTokens []string `json:"fileTokens"`
	}
	jsonData, err := json.Marshal(payload{FileTokens: filesToAssociate})
	if err != nil {
		return nil, err
	}
	r, err := re.createRequest(ctx, http.MethodPost, albumUrl, bytes.NewBuffer(jsonData), nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, mod.WaifuAlbum{})
}

func (re *api) DisassociateFiles(ctx context.Context, albumToken string, filesToDisassociate []string) (*mod.WaifuAlbum, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/%s/disassociate", albumToken))
	type payload struct {
		FileTokens []string `json:"fileTokens"`
	}
	jsonData, err := json.Marshal(payload{FileTokens: filesToDisassociate})
	if err != nil {
		return nil, err
	}
	r, err := re.createRequest(ctx, http.MethodPost, albumUrl, bytes.NewBuffer(jsonData), nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, mod.WaifuAlbum{})
}

func (re *api) GetAlbum(ctx context.Context, albumToken string) (*mod.WaifuAlbum, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/%s", albumToken))
	r, err := re.createRequest(ctx, http.MethodGet, albumUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, mod.WaifuAlbum{})
}

func (re *api) DeleteAlbum(ctx context.Context, albumToken string, deleteFiles bool) (*mod.GenericSuccess, error) {
	albumUrl := getUrl(map[string]any{"deleteFiles": deleteFiles}, fmt.Sprintf("album/%s", albumToken))
	r, err := re.createRequest(ctx, http.MethodDelete, albumUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, mod.GenericSuccess{})
}

func (re *api) ShareAlbum(ctx context.Context, albumToken string) (string, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/share/%s", albumToken))
	r, err := re.createRequest(ctx, http.MethodGet, albumUrl, nil, nil)
	if err != nil {
		return "", err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return "", err
	}
	result, err := readResponse(resp, mod.GenericSuccess{})
	if err != nil {
		return "", err
	}
	return result.Description, nil
}

func (re *api) RevokeAlbum(ctx context.Context, albumToken string) (*mod.GenericSuccess, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/revoke/%s", albumToken))
	r, err := re.createRequest(ctx, http.MethodGet, albumUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, mod.GenericSuccess{})
}

func (re *api) DownloadAlbum(ctx context.Context, albumToken string, files []int) ([]byte, error) {
	albumUrl := getUrl(nil, fmt.Sprintf("album/download/%s", albumToken))
	jsonData, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}
	r, err := re.createRequest(ctx, http.MethodPost, albumUrl, bytes.NewBuffer(jsonData), nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = checkError(resp)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
