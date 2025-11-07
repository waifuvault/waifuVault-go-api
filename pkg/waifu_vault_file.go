package waifuVault

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/waifuvault/waifuVault-go-api/pkg/mod"
)

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
