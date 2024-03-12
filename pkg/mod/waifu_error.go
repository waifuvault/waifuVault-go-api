package mod

type WaifuError struct {
	Name    string // The name of the HTTP status. e.g: Bad Request
	Message string // the message or reason why the request failed
	Status  int    // the http status returned
}
