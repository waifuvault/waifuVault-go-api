package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"waifuVault-go-api/pkg/mod"
)

const baseUrl = "https://waifuvault.moe"

var client = &http.Client{}
var WaifuVault = waifuvalt{}

func main() {
	/*url := "http://localhost:3001/"
	// don't worry about errors
	response, _ := http.Get(url)
	var target = mod.WaifuResponse[string]{}
	err := json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		return
	}
	fmt.Print(target)*/
	/*fileDir, _ := os.Getwd()
	fileName := "main.go"
	filePath := path.Join(fileDir, "pkg", fileName)
	b, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
		return
	}
	fileNames := "main.go"*/

	opts := mod.WaifuvaultPutOpts{
		Url:          "https://victorique.moe/img/slider/Quotes.jpg",
		HideFilename: true,
		Expires:      "1h",
	}
	file, err := WaifuVault.UploadFile(opts)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(file)
}

type waifuvalt struct {
}

func (re *waifuvalt) UploadFile(options mod.WaifuvaultPutOpts) (*mod.WaifuResponse[string], error) {
	if options.File != nil && options.Bytes != nil && options.Url != "" || options.File == nil && options.Bytes == nil && options.Url == "" {
		return nil, errors.New("you can only supply buffer, file or url")
	}

	body := bytes.Buffer{}
	var writer *multipart.Writer
	if options.File != nil || options.Bytes != nil {

		var fileBytes *bytes.Buffer
		var w io.Writer
		var err error
		writer = multipart.NewWriter(&body)
		if options.File != nil {
			fileBytes = bytes.NewBuffer(nil)
			_, err = io.Copy(fileBytes, options.File)
			if err != nil {
				return nil, err
			}
			w, err = writer.CreateFormFile("file", filepath.Base(options.File.Name()))
		} else {
			fileBytes = bytes.NewBuffer(*options.Bytes)
			if options.FileName == "" {
				return nil, errors.New("FileName must be set if bytes is used")
			}
			w, err = writer.CreateFormFile("file", options.FileName)
		}

		if err != nil {
			return nil, err
		}

		if _, err = w.Write(fileBytes.Bytes()); err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
	} else if options.Url != "" {
		bodyUrl := fmt.Sprintf(`{"url": "%s"}`, options.Url)
		body = *bytes.NewBuffer([]byte(bodyUrl))
	}

	uploadUrl := getUrl(map[string]any{
		"expires":       options.Expires,
		"hide_filename": options.HideFilename,
		"password":      options.Password,
	}, nil)

	r, err := createRequest(http.MethodPut, uploadUrl, &body, writer)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return getResponse[string](resp)
}

func createRequest(method, url string, body io.Reader, writer *multipart.Writer) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
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

func (re *waifuvalt) FileInfo(token string) (*mod.WaifuResponse[string], error) {
	return nil, nil
}

func (re *waifuvalt) FileInfoFormatted(token string) (*mod.WaifuResponse[int], error) {
	return nil, nil
}

func (re *waifuvalt) DeleteFile(token string) (bool, error) {
	return false, nil
}

func (re *waifuvalt) GetFile(options mod.GetFileInfo) ([]byte, error) {
	return nil, nil
}

func getUrl(obj map[string]any, path *string) string {
	baseRestUrl := fmt.Sprintf("%s/rest", baseUrl)
	if path != nil {
		baseRestUrl += "/" + *path
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
	err := checkError(response)
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(response.Body)
	var target = &mod.WaifuResponse[T]{}
	err = json.Unmarshal(bodyBytes, target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func checkError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(response.Body)
		errStr := string(bodyBytes)

		var respErrorJson mod.WaifuError
		jsonErr := json.Unmarshal(bodyBytes, &respErrorJson)

		if jsonErr != nil {
			errStr = fmt.Sprintf("Error %d (%s): %s", respErrorJson.Status, respErrorJson.Name, respErrorJson.Message)
		}

		return errors.New(errStr)
	}
	return nil
}

func ToPtr[T any](x T) *T {
	return &x
}
