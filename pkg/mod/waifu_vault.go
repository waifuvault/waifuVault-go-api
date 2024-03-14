package mod

import "context"

type Waifuvalt interface {
	// UploadFile - Upload a file using a byte array, url or file
	UploadFile(options WaifuvaultPutOpts, ctx *context.Context) (*WaifuResponse[string], error)

	// FileInfo - Obtain file info such as URL and retention period (as an epoch timestamp)
	FileInfo(token string, ctx *context.Context) (*WaifuResponse[int], error)

	// FileInfoFormatted - Same as FileInfo, but instead returns the retention period as a human-readable string
	FileInfoFormatted(token string, ctx *context.Context) (*WaifuResponse[string], error)

	// DeleteFile - Delete a file given a token
	DeleteFile(token string, ctx *context.Context) (bool, error)

	// GetFile - Download the file given options and return a byte array of said file
	GetFile(options GetFileInfo, ctx *context.Context) ([]byte, error)
}
