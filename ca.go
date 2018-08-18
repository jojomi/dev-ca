package main

// CARequest is a request for creation of a Certificate Authority with given parameters
type CARequest struct {
	RSAKeySize int
	Days       int

	Subj string
	Name string

	FilenameKey string
	FilenamePem string
}
