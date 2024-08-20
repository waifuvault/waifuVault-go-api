package mod

import "os"

type WaifuvaultPutOpts struct {

	// A string containing a number and a letter of `m` for mins, `h` for hours, `d` for days.
	// For example, `1h` would be 1 hour and `1d` would be 1 day.
	// Omit this if you want the file to exist, according to the retention policy
	Expires string

	// if set to true, then your filename will not appear in the URL. if false, then it will appear in the URL. defaults to false
	HideFilename bool

	// Setting a password will encrypt the file
	Password string

	//The file object on disk
	File *os.File

	// the raw bytes of the file
	Bytes *[]byte

	//An url to the file you want uploaded
	Url string

	// The filename if `Bytes` is used
	FileName string

	// If this is true, then the file will be deleted as soon as it is accessed
	OneTimeDownload bool `json:"oneTimeDownload"`

	// If supplied, this file will be associated to that bucket
	BucketToken string
}
