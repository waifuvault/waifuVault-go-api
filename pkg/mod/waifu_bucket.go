package mod

type WaifuBucket struct {
	// the token of the bucket
	Token string `json:"token"`

	// The files contained in this bucket
	Files []WaifuResponse[int] `json:"files"`
}
