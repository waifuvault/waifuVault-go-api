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

	// CreateBucket - create a new bucket, buckets are bound to your IP, so you may only have one bucket per IP
	CreateBucket(ctx context.Context) (*WaifuBucket, error)

	// GetBucket - Get a bucket and all the files it contains
	GetBucket(ctx context.Context, token string) (*WaifuBucket, error)

	// DeleteBucket - Delete a bucket and all files it contains
	DeleteBucket(ctx context.Context, token string) (bool, error)

	// CreateAlbum - Create an album with the given name
	CreateAlbum(ctx context.Context, body WaifuAlbumCreateBody) (*WaifuAlbum, error)

	// AssociateFiles - Associate files with an album, the album must exist, and the files must be in the same bucket as the album
	AssociateFiles(ctx context.Context, albumToken string, filesToAssociate []string) (*WaifuAlbum, error)

	// DisassociateFiles - Remove files from the album
	DisassociateFiles(ctx context.Context, albumToken string, filesToDisassociate []string) (*WaifuAlbum, error)

	// GetAlbum - Get an album and all of its files
	GetAlbum(ctx context.Context, albumToken string) (*WaifuAlbum, error)

	// DeleteAlbum - Deletes an album and optionally, deletes all associated files with the album
	DeleteAlbum(ctx context.Context, albumToken string, deleteFiles bool) (*GenericSuccess, error)

	// ShareAlbum - Sharing an album makes it so others can see it in a read-only view, returns the public URL
	ShareAlbum(ctx context.Context, albumToken string) (string, error)

	// RevokeAlbum - Revoking an album invalidates the URL used to view it and makes it private
	RevokeAlbum(ctx context.Context, albumToken string) (*GenericSuccess, error)

	// DownloadAlbum - Download an album or selected files from an album, returns a ZIP file as bytes
	DownloadAlbum(ctx context.Context, albumToken string, files []int) ([]byte, error)
}
