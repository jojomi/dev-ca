package main

// CertificateRequest is a request for a new certificate to be issued by a Certificate Authority with given parameters
type CertificateRequest struct {
	Days       int
	RSAKeySize int

	Domains []string

	FilenameCAPem      string
	FilenameCAKey      string
	FilenameCsr        string
	FilenamePrivateKey string
	FilenameCrt        string
	FilenameConfig     string
}
