package mod

import "context"

type Waifuvalt interface {
	// UploadFile - Upload a file using a byte array, url or file
	UploadFile(ctx context.Context, options WaifuvaultPutOpts) (*WaifuResponse[string], error)

	// FileInfo - Obtain file info such as URL and retention period (as an epoch timestamp)
	FileInfo(ctx context.Context, token string) (*WaifuResponse[int], error)

	// FileInfoFormatted - Same as FileInfo, but instead returns the retention period as a human-readable string
	FileInfoFormatted(ctx context.Context, token string) (*WaifuResponse[string], error)

	// DeleteFile - Delete a file given a token
	DeleteFile(ctx context.Context, token string) (bool, error)

	// GetFile - Download the file given options and return a byte array of said file
	GetFile(ctx context.Context, options GetFileInfo) ([]byte, error)

	// ModifyFile - modify an entry
	ModifyFile(ctx context.Context, token string, options ModifyEntryPayload) (*WaifuResponse[int], error)
}
