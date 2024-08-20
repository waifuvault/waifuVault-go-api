package waifuVault

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
)

const baseUrl = "https://waifuvault.moe"

type api struct {
	client http.Client
}

func NewWaifuvaltApi(client http.Client) mod.Waifuvalt {
	return &api{
		client: client,
	}
}

func (re *api) UploadFile(ctx context.Context, options mod.WaifuvaultPutOpts) (*mod.WaifuResponse[string], error) {
	if options.File != nil && options.Bytes != nil && options.Url != "" || options.File == nil && options.Bytes == nil && options.Url == "" {
		return nil, errors.New("you can only supply buffer, file or url")
	}
	body := bytes.Buffer{}
	var writer *multipart.Writer
	if options.File != nil || options.Bytes != nil {

		var fileBytes *bytes.Buffer
		var fileFormWriter io.Writer
		var err error

		writer = multipart.NewWriter(&body)
		if options.File != nil {
			fileBytes = bytes.NewBuffer(nil)
			_, err = io.Copy(fileBytes, options.File)
			if err != nil {
				return nil, err
			}
			fileFormWriter, err = writer.CreateFormFile("file", filepath.Base(options.File.Name()))

		} else {
			fileBytes = bytes.NewBuffer(*options.Bytes)
			if options.FileName == "" {
				return nil, errors.New("FileName must be set if bytes is used")
			}
			fileFormWriter, err = writer.CreateFormFile("file", options.FileName)
		}

		if options.Password != "" {
			passwordFormWriter, err := writer.CreateFormField("password")
			if err != nil {
				return nil, err
			}
			if _, err = passwordFormWriter.Write([]byte(options.Password)); err != nil {
				return nil, err
			}
		}

		if err != nil {
			return nil, err
		}

		if _, err = fileFormWriter.Write(fileBytes.Bytes()); err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
	} else if options.Url != "" {
		var bodyUrl string
		if options.Password != "" {
			bodyUrl = fmt.Sprintf(`{"url": "%s", "password": "%s"}`, options.Url, options.Password)
		} else {
			bodyUrl = fmt.Sprintf(`{"url": "%s"}`, options.Url)
		}
		body = *bytes.NewBuffer([]byte(bodyUrl))
	}

	uploadUrl := getUrl(map[string]any{
		"expires":           options.Expires,
		"hide_filename":     options.HideFilename,
		"one_time_download": options.OneTimeDownload,
	}, options.BucketToken)

	r, err := re.createRequest(ctx, http.MethodPut, uploadUrl, &body, writer)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return getResponse[string](resp)
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

func (re *api) FileInfo(ctx context.Context, token string) (*mod.WaifuResponse[int], error) {
	resp, err := re.createGetRequestForFileInfo(ctx, token, false)
	if err != nil {
		return nil, err
	}
	return getResponse[int](resp)
}

func (re *api) FileInfoFormatted(ctx context.Context, token string) (*mod.WaifuResponse[string], error) {
	resp, err := re.createGetRequestForFileInfo(ctx, token, true)
	if err != nil {
		return nil, err
	}
	return getResponse[string](resp)
}

func (re *api) createGetRequestForFileInfo(ctx context.Context, token string, isFormatted bool) (*http.Response, error) {
	getUrl := getUrl(map[string]any{"formatted": isFormatted}, token)
	r, err := re.createRequest(ctx, http.MethodGet, getUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	return re.client.Do(r)
}

func (re *api) DeleteFile(ctx context.Context, token string) (bool, error) {
	deleteUrl := getUrl(nil, token)

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

func (re *api) GetFile(ctx context.Context, options mod.GetFileInfo) ([]byte, error) {

	if options.Filename == "" && options.Token == "" {
		return nil, errors.New("please supply a token or a filename")
	}
	var fileUrl string
	if options.Filename != "" {
		fileUrl = fmt.Sprintf("%s/f/%s", baseUrl, options.Filename)
	} else {
		fileInfo, err := re.FileInfo(ctx, options.Token)
		if err != nil {
			return nil, err
		}
		fileUrl = fileInfo.URL
	}

	r, err := re.createRequest(ctx, http.MethodGet, fileUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	if options.Password != "" {
		r.Header.Set("x-password", options.Password)
	}

	resp, err := re.client.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New("password is incorrect")
	}

	err = checkError(resp)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func (re *api) ModifyFile(ctx context.Context, token string, options mod.ModifyEntryPayload) (*mod.WaifuResponse[int], error) {
	uploadUrl := getUrl(nil, token)

	jsonData, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	r, err := re.createRequest(ctx, http.MethodPatch, uploadUrl, bytes.NewBuffer(jsonData), nil)
	if err != nil {
		return nil, err
	}
	resp, err := re.client.Do(r)
	if err != nil {
		return nil, err
	}
	return getResponse[int](resp)
}

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

func getBucketResponse(response *http.Response) (*mod.WaifuBucket, error) {
	return readResponse(response, mod.WaifuBucket{})
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
