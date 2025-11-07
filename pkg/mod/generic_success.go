package mod

// GenericSuccess represents a generic success response from the API
type GenericSuccess struct {
	// Success indicates if the operation was successful
	Success bool `json:"success"`

	// Description provides additional information about the result
	Description string `json:"description"`
}
