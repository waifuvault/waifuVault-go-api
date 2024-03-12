package mod

type GetFileInfo struct {
	// Password for this file
	Password string
	// the file token
	Token string
	// the filename and the file upload epoch. for example, 1710111505084/08.png.
	// files with hidden filenames will only contain the epoch with ext. for example, 1710111505084.png
	Filename string
}
