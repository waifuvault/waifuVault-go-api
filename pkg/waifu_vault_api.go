package waifuVault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
)

var baseUrl = "https://waifuvault.moe"

type api struct {
	client http.Client
}

func NewWaifuvaltApi(client http.Client) mod.Waifuvalt {
	return &api{
		client: client,
	}
}

func (re *api) createRequest(ctx context.Context, method, url string, body io.Reader, writer *multipart.Writer) (*http.Request, error) {

	r, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, err
	}
	if writer != nil {
		r.Header.Add("Content-Type", writer.FormDataContentType())
	} else {
		r.Header.Set("Content-Type", "application/json")
	}

	return r, nil
}

func getUrl(obj map[string]any, path string) string {
	baseRestUrl := fmt.Sprintf("%s/rest", baseUrl)
	if path != "" {
		baseRestUrl = fmt.Sprintf("%s/%s", baseRestUrl, path)
	}
	if obj == nil {
		return baseRestUrl
	}

	params := url.Values{}
	for key, val := range obj {
		if val == nil || val == "" {
			continue
		}
		params.Add(key, fmt.Sprintf("%v", val))
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", baseRestUrl, params.Encode())
	}

	return baseRestUrl
}

func getResponse[T string | int](response *http.Response) (*mod.WaifuResponse[T], error) {
	return readResponse(response, mod.WaifuResponse[T]{})
}

func readResponse[T any](response *http.Response, target T) (*T, error) {
	defer response.Body.Close()
	err := checkError(response)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &target)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func checkError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(response.Body)
		errStr := string(bodyBytes)

		var respErrorJson mod.WaifuError
		jsonErr := json.Unmarshal(bodyBytes, &respErrorJson)

		if jsonErr == nil {
			errStr = fmt.Sprintf("Error %d (%s): %s", respErrorJson.Status, respErrorJson.Name, respErrorJson.Message)
		}

		return errors.New(errStr)
	}
	return nil
}
