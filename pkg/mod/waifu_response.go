package mod

// WaifuResponse is the response from the api for files and uploads
type WaifuResponse[T string | int] struct {

	// Token for the uploaded file
	Token string `json:"token"`

	// URL to the uploaded file
	URL string `json:"url"`

	// Options The options for this upload
	Options WaifuResponseOptions `json:"options"`

	// RetentionPeriod is a string or a number that represents
	// when the file will expire, if called with `format` true, then
	// this will be a string like "332 days 7 hours 18 minutes 8 seconds"
	RetentionPeriod T `json:"retentionPeriod"`

	// Bucket - the bucket this belongs to
	Bucket string `json:"bucket"`

	// ID is the public ID of this file
	ID int `json:"id"`

	// Views is how many people have downloaded this file
	Views int `json:"views"`
}
