# waifuvault-node-api

This contains the official API bindings for uploading, deleting and obtaining files
with [waifuvault.moe](https://waifuvault.moe/). Contains a full up-to-date API for interacting with the service

## Installation

```sh
go get github.com/waifuvault/waifuVault-go-api@v1.0.0
```

## Usage

This API contains 4 interactions:

1. Upload
2. Delete
3. get file info
4. get file

The package is namespaced to `waifuVault`, so to import it, simply:

```go
import "github.com/waifuvault/waifuVault-go-api/pkg"

// then init the API

api := waifuVault.NewWaifuvaltApi(http.Client{})
```

### Contexts

Each function on the Waifuvalt API takes a context. This can be `nil`

```go
package main

import (
	"context"
	"fmt"
	waifuVault "github.com/waifuvault/waifuVault-go-api/pkg"
	waifuMod "github.com/waifuvault/waifuVault-go-api/pkg/mod" // namespace mod
	"net/http"
)

func main() {
	cx, cancel := context.WithCancel(context.Background())
	api := waifuVault.NewWaifuvaltApi(http.Client{})
	file, err := api.UploadFile(waifuMod.WaifuvaultPutOpts{
		Url: "https://waifuvault.moe/assets/custom/images/08.png",
	}, &cx)
	if err != nil {
		return
	}
	api.DeleteFile(file.Token, nil)
	fmt.Printf(file.URL) // the URL
}
```

### Upload File

To Upload a file, use the `UploadFile` function. This function takes the following options as struct:

| Option         | Type       | Description                                                               | Required                                       | Extra info                                                                        |
|----------------|------------|---------------------------------------------------------------------------|------------------------------------------------|-----------------------------------------------------------------------------------|
| `File`         | `*os.File` | The file to upload. This is an *os.File                                   | true only if `Url` or `Bytes` is not supplied  | If `Url` or `Bytes` is supplied, this prop can't be set                           |
| `Url`          | `string`   | The URL to a file that exists on the internet                             | true only if `File` or `Bytes` is not supplied | If `File` or `Bytes` is supplied, this prop can't be set                          |
| `Bytes`        | `*[]byte`  | The raw Bytes to of the file to upload.                                   | true only if `File` or `Url` is not supplied   | If `File` or `Url` is supplied, this prop can't be set and `FileName` MUST be set |
| `Expires`      | `string`   | A string containing a number and a unit (1d = 1day)                       | false                                          | Valid units are `m`, `h` and `d`                                                  |
| `HideFilename` | `bool`     | If true, then the uploaded filename won't appear in the URL               | false                                          | Defaults to `false`                                                               |
| `Password`     | `string`   | If set, then the uploaded file will be encrypted                          | false                                          |                                                                                   |
| `FileName`     | `string`   | Only used if `Bytes` is set, this will be the filename used in the upload | true only if `Bytes` is set                    |                                                                                   |

Using a URL:

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	waifuMod "github.com/waifuvault/waifuVault-go-api/pkg/mod" // namespace mod
	"net/http"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})
	file, err := api.UploadFile(waifuMod.WaifuvaultPutOpts{
		Url: "https://waifuvault.moe/assets/custom/images/08.png",
	}, nil)
	if err != nil {
		return
	}
	fmt.Printf(file.URL) // the URL
}
```

Using Bytes:

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	waifuMod "github.com/waifuvault/waifuVault-go-api/pkg/mod"
	"net/http"
	"os"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})

	b, err := os.ReadFile("myCoolFile.jpg")
	if err != nil {
		fmt.Print(err)
	}

	file, err := api.UploadFile(waifuMod.WaifuvaultPutOpts{
		Bytes:    &b,
		FileName: "myCoolFile.jpg", // make sure you supply the extension
	}, nil)
	if err != nil {
		return
	}
	fmt.Printf(file.URL) // the URL
}
```

Using a file:

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	waifuMod "github.com/waifuvault/waifuVault-go-api/pkg/mod"
	"net/http"
	"os"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})

	fileStruc, err := os.Open("myCoolFile.jpg")
	if err != nil {
		fmt.Print(err)
	}

	file, err := api.UploadFile(waifuMod.WaifuvaultPutOpts{
		File: fileStruc,
	}, nil)
	if err != nil {
		return
	}
	fmt.Printf(file.URL) // the URL
}
```

### File Info

If you have a token from your upload. Then you can get file info. This results in the following info:

* token
* url
* protected
* retentionPeriod

Use the `FileInfo` function. This function takes the following options as parameters:

| Option  | Type     | Description             | Required | Extra info |
|---------|----------|-------------------------|----------|------------|
| `token` | `string` | The token of the upload | true     |            |

If you want the `retentionPeriod` to be a human readble string and not a epoch, you can use `FileInfoFormatted` that
takes the same parameters

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	"net/http"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})
	info, err := api.FileInfo("token", nil)
	if err != nil {
		return
	}
	fmt.Print(info.RetentionPeriod) // the retention period as epoch number
	fmt.Print(info.URL)             // the URL
}
```

Human-readable timestamp:

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	"net/http"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})
	info, err := api.FileInfoFormatted("token", nil)
	if err != nil {
		return
	}
	fmt.Print(info.RetentionPeriod) // the retention period as a string 
	fmt.Print(info.URL)             // the URL
}
```

### Delete File

To delete a file, you must supply your token to the `DeleteFile` function.

This function takes the following options as parameters:

| Option  | Type     | Description                              | Required | Extra info |
|---------|----------|------------------------------------------|----------|------------|
| `token` | `string` | The token of the file you wish to delete | true     |            |

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	"net/http"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})
	success, err := api.DeleteFile("token", nil)
	if err != nil {
		return
	}
	fmt.Print(success) // true or false
}
```

### Get File

This lib also supports obtaining a file from the API as a Byte Array by supplying either the token or the unique
identifier
of the file (epoch/filename).

Use the `GetFile` function. This function takes the following options an object:

| Option     | Type     | Description                                                                                      | Required                           | Extra info                                               |
|------------|----------|--------------------------------------------------------------------------------------------------|------------------------------------|----------------------------------------------------------|
| `Token`    | `string` | The token of the file you want to download                                                       | true only if `filename` is not set | if `filename` is set, then this can not be used          |
| `FileName` | `string` | The Unique identifier of the file, this is the epoch time stamp it was uploaded and the filename | true only if `token` is not set    | if `token` is set, then this can not be used             |
| `Password` | `string` | The password for the file if it is protected                                                     | false                              | Must be supplied if the file is uploaded with `password` |

> **Important!** The Unique identifier filename is the epoch/filename only if the file uploaded did not have a hidden
> filename, if it did, then it's just the epoch.
> For example: `1710111505084/08.png` is the Unique identifier for a standard upload of a file called `08.png`, if this
> was uploaded with hidden filename, then it would be `1710111505084.png`

Obtain an encrypted file

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	waifuMod "github.com/waifuvault/waifuVault-go-api/pkg/mod"
	"net/http"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})

	// upload the file
	file, err := api.UploadFile(waifuMod.WaifuvaultPutOpts{
		Url:      "https://waifuvault.moe/assets/custom/images/08.png",
		Password: "foobar",
	}, nil)

	// download the file
	bytes, err := api.GetFile(waifuMod.GetFileInfo{
		Password: "foobar",
		Token:    file.Token,
	}, nil)
	if err != nil {
		return
	}
	fmt.Print(bytes) // byte array
}
```

Obtain a file from Unique identifier

```go
package main

import (
	"fmt"
	"github.com/waifuvault/waifuVault-go-api/pkg"
	waifuMod "github.com/waifuvault/waifuVault-go-api/pkg/mod"
	"net/http"
)

func main() {
	api := waifuVault.NewWaifuvaltApi(http.Client{})

	bytes, err := api.GetFile(waifuMod.GetFileInfo{
		Filename: "/1710111505084/08.png",
	}, nil)
	if err != nil {
		return
	}
	fmt.Print(bytes) // byte array
}
```
