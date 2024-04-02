package mod

type WaifuResponseOptions struct {
	// HideFilename If the filename is hidden
	HideFilename bool `json:"hideFilename"`

	// OneTimeDownload  If this file will be deleted when it is accessed
	OneTimeDownload bool `json:"oneTimeDownload"`

	// Protected is if this file is protected-protected/encrypted
	Protected bool `json:"protected"`
}
