package mod

// ModifyEntryPayload A modify entry request to change aspects of the entry
type ModifyEntryPayload struct {

	// The new password.
	// If the file is not currently encrypted, then this will encrypt it with the new password if it is encrypted, Then this will change the password (`previousPassword` will need to be set in this case)
	// set this to an empty string `""` to remove protection and decrypt the file
	Password *string `json:"password"`

	// If changing a password, then this will need to be set
	PreviousPassword *string `json:"previousPassword"`

	// same as WaifuvaultPutOpts.Expires
	CustomExpiry *string `json:"customExpiry"`

	// hide the filename. use the new URL in the response to get the new URL to use
	HideFilename *bool `json:"hideFilename"`
}
