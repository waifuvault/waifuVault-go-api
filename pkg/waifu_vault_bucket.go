package waifuVault

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
)

func (re *api) CreateBucket(ctx context.Context) (*mod.WaifuBucket, error) {
	restUrl := baseUrl + "/rest/bucket/create"
	r, err := re.createRequest(ctx, http.MethodGet, restUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return getBucketResponse(resp)
}

func (re *api) GetBucket(ctx context.Context, token string) (*mod.WaifuBucket, error) {
	restUrl := baseUrl + "/rest/bucket/get"
	type payload struct {
		BucketToken string `json:"bucket_token"`
	}
	jsonData, err := json.Marshal(&payload{token})
	r, err := re.createRequest(ctx, http.MethodPost, restUrl, bytes.NewBuffer(jsonData), nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return getBucketResponse(resp)
}

func (re *api) DeleteBucket(ctx context.Context, token string) (bool, error) {
	deleteUrl := baseUrl + "/rest/bucket/" + token
	r, err := re.createRequest(ctx, http.MethodDelete, deleteUrl, nil, nil)
	if err != nil {
		return false, err
	}
	resp, err := re.client.Do(r)
	defer resp.Body.Close()
	err = checkError(resp)
	if err != nil {
		return false, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	return bodyString == "true", nil
}

func getBucketResponse(response *http.Response) (*mod.WaifuBucket, error) {
	return readResponse(response, mod.WaifuBucket{})
}
