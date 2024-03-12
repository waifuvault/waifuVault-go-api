package mod

type Waifuvalt interface {
	UploadFile(options WaifuvaultPutOpts) (*WaifuResponse[string], error)
	FileInfo(token string) (*WaifuResponse[string], error)
	FileInfoFormatted(token string) (*WaifuResponse[int], error)
	DeleteFile(token string) (bool, error)
	GetFile(options GetFileInfo) ([]byte, error)
}
