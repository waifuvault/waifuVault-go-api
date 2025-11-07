package mod

// AlbumStub is an album but with files omitted
type AlbumStub struct {
	// Token is the private token of this album
	Token string `json:"token"`

	// Bucket is the token of the bucket this album belongs to
	Bucket string `json:"bucket"`

	// PublicToken is the public token used to share the album with others
	PublicToken *string `json:"publicToken"`

	// Name is the name of the album
	Name string `json:"name"`

	// DateCreated is the date this album was created (epoch timestamp)
	DateCreated int64 `json:"dateCreated"`
}
