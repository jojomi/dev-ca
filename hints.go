package main

import (
	"fmt"

	"github.com/jojomi/go-script/print"
	"github.com/jojomi/strtpl"
)

func hintCertApache(cr *CertificateRequest) {
	print.Boldln("Apache configuration")

	fmt.Println(strtpl.MustEval(`To use these files on your Apache webserver, simply copy both {{.FilenameCrt}} and {{.FilenamePrivateKey}} to your webserver, and include them like this:
    SSLEngine On
    SSLCertificateFile    /path_to_your_files/{{.FilenameCrt}}"
    SSLCertificateKeyFile /path_to_your_files/{{.FilenamePrivateKey}}"`, cr))
}

func hintCAFirefox(cr *CARequest) {
	print.Boldln("Firefox configuration")

	fmt.Println(strtpl.MustEval(`To trust this CA in Firefox, go to Settings, Security, Certificates, Show Certificates…, Certificate Authorities, Import… and select {{.FilenamePem}}.`, cr))
}

func hintCAMacOSX(cr *CARequest) {
	print.Boldln("MacOS X configuration (Chrome)")

	fmt.Println(strtpl.MustEval(`To trust this CA in Mac OS X, go to Keychain, Import Objects… and select {{.FilenamePem}}. Then double click it and set to trusted.`, cr))
}
