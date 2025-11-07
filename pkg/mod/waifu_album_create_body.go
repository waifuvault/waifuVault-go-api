package mod

// WaifuAlbumCreateBody is the payload for creating a new album
type WaifuAlbumCreateBody struct {
	// Name is the name of the album
	Name string `json:"name"`

	// BucketToken is the bucket this album will be created in
	BucketToken string `json:"bucketToken"`
}
