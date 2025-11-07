package waifuVault

import "github.com/waifuvault/waifuVault-go-api/pkg/mod"

const (
	MockUUID1 = "8c3d4527-4cea-4cb8-8171-002b158693ab"
	MockUUID2 = "3ed6dafa-b56a-4207-b004-852abb6fea11"
	MockUUID3 = "c18ff3cb-442d-44e7-b5d1-c43a83b3a1a4"
)

var (
	MockUUID3Ptr = MockUUID3
)

var WaifuResponseMock1 = mod.WaifuResponse[int]{
	URL:   "https://waifuvault.moe/f/1710111505084/08.png",
	Token: "123-fake-street",
	Options: mod.WaifuResponseOptions{
		Protected:       false,
		OneTimeDownload: false,
		HideFilename:    false,
	},
	RetentionPeriod: 1234,
	Bucket:          "123-fake-street-bucket",
	ID:              1,
	Views:           0,
}

var WaifuResponseMock2 = mod.WaifuResponse[string]{
	URL:   "https://waifuvault.moe/f/1710111505084/08.png",
	Token: "123-fake-street-2",
	Options: mod.WaifuResponseOptions{
		Protected:       true,
		OneTimeDownload: false,
		HideFilename:    false,
	},
	RetentionPeriod: "1234",
	Bucket:          "",
	ID:              2,
	Views:           0,
}

var WaifuBucketMock1 = mod.WaifuBucket{
	Token:  "123-fake-street-bucket",
	Files:  []mod.WaifuResponse[int]{WaifuResponseMock1},
	Albums: []mod.AlbumStub{},
}

var WaifuErrorMock = mod.WaifuError{
	Status:  400,
	Message: "loser",
	Name:    "whore",
}

var GenericSuccessDeletedMock = mod.GenericSuccess{
	Description: "deleted",
	Success:     true,
}

var WaifuAlbumMock1 = mod.WaifuAlbum{
	BucketToken: MockUUID1,
	DateCreated: 0,
	Files:       []mod.WaifuResponse[int]{WaifuResponseMock1},
	Name:        "album1",
	PublicToken: &MockUUID3Ptr,
	Token:       MockUUID2,
}

var WaifuAlbumMock2 = mod.WaifuAlbum{
	BucketToken: MockUUID1,
	DateCreated: 0,
	Files:       []mod.WaifuResponse[int]{},
	Name:        "album2",
	PublicToken: &MockUUID3Ptr,
	Token:       MockUUID2,
}

var SharedFileMock1 = mod.GenericSuccess{
	Description: "sharedAlbum.foo",
	Success:     true,
}
