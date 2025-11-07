package mod

// WaifuAlbum is a public collection of files, it can be shared with others in a read-only fashion
type WaifuAlbum struct {
	// Token is the private token of this album
	Token string `json:"token"`

	// BucketToken is the token of the bucket this album belongs to
	BucketToken string `json:"bucketToken"`

	// PublicToken is the public token used to share the album with others
	PublicToken *string `json:"publicToken"`

	// Name is the name of the album
	Name string `json:"name"`

	// Files are the files contained in this album
	Files []WaifuResponse[int] `json:"files"`

	// DateCreated is the date this album was created (epoch timestamp)
	DateCreated int64 `json:"dateCreated"`
}
