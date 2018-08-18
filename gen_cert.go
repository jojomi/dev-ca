package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	textTemplate "text/template"

	"github.com/jojomi/go-script"
	"github.com/jojomi/go-script/interview"
	"github.com/jojomi/go-script/print"
	"github.com/jojomi/strtpl"
	"github.com/spf13/cobra"
)

var (
	flagGenCertRSASize int
	flagGenCertDays    int
	flagGenCertName    string
	flagGenCertDomains string

	certFolder = "certs"
)

func genCertCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "gen-cert",
		Short: "Generate a certificate using the Certificate Authority",
		Run:   genCertHandler,
	}
	f := c.PersistentFlags()
	f.IntVarP(&flagGenCertRSASize, "rsa-size", "r", 2048, "RSA key size")
	f.IntVarP(&flagGenCertDays, "days", "", 365*10, "days of validity")
	f.StringVarP(&flagGenCertName, "name", "n", "", "Name of the certificate, if empty first --domain name given")
	f.StringVarP(&flagGenCertDomains, "domains", "d", "", "domains to issue the certificate for, comma separated")
	return &c
}

func genCertHandler(cmd *cobra.Command, args []string) {
	c := script.NewContext()
	c.MustCommandExist("openssl")
	outputFolder := certFolder
	err := c.EnsureDirExists(outputFolder, 0700)
	checkFail(err)
	if flagGenCertDomains == "" {
		fmt.Println("no domains given, see --domains flag")
		os.Exit(1)
	}
	domains := strings.Split(flagGenCertDomains, ",")
	// Create a new private key if one doesnt exist, or use the existing one if it does
	basename := flagGenCertName
	if basename == "" {
		basename = strings.Replace(domains[0], "*", "star", -1)
	}
	certificateRequest := &CertificateRequest{
		RSAKeySize:         flagGenCertRSASize,
		Days:               flagGenCertDays,
		Domains:            domains,
		FilenameCAPem:      filepath.Join(caFolder, "rootCA.pem"),
		FilenameCAKey:      filepath.Join(caFolder, "rootCA.key"),
		FilenameCsr:        filepath.Join(outputFolder, basename+".csr"),
		FilenameCrt:        filepath.Join(outputFolder, basename+".crt"),
		FilenamePrivateKey: filepath.Join(outputFolder, basename+".key"),
	}

	if !c.FileExists(certificateRequest.FilenameCAKey) || !c.FileExists(certificateRequest.FilenameCAPem) {
		print.Errorln("No CA found, please use gen-ca command to create one first.")
		os.Exit(1)
	}

	if c.FileExists(certificateRequest.FilenamePrivateKey) || c.FileExists(certificateRequest.FilenameCrt) {
		overwrite, errConfirm := interview.Confirm("Files exist, overwrite?", false)
		if errConfirm != nil || !overwrite {
			os.Exit(1)
		}
	}

	fullCommand := strtpl.MustEval(`openssl req -new -newkey rsa:{{.RSAKeySize}} -sha256 -nodes -keyout "{{.FilenamePrivateKey}}" -out "{{.FilenameCsr}}" -subj "/C=/ST=/L=/O=/CN={{index .Domains 0}}"`, certificateRequest)
	exec(c, fullCommand)

	templateString := `
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
{{ range $i, $d := .Domains -}}
	DNS.{{ add $i 1 }} = {{ $d }}
{{ end }}
`
	funcMap := textTemplate.FuncMap{
		"add": func(summands ...int) string {
			sum := 0
			for _, summand := range summands {
				sum += summand
			}
			return strconv.Itoa(sum)
		},
	}
	configString := strtpl.MustEvalWithFuncMap(templateString, funcMap, certificateRequest)
	if flagVerbose {
		print.Boldln("Config:")
		fmt.Println(configString)
	}
	f, err := c.TempFile()
	if err != nil {
		fmt.Println("error writing config")
		os.Exit(2)
	}
	_, err = f.WriteString(configString)
	checkFail(err)
	err = f.Close()
	checkFail(err)
	certificateRequest.FilenameConfig = f.Name()

	// openssl x509 -req -in device.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out device.crt -days $NUM_OF_DAYS -sha256 -extfile /tmp/__v3.ext
	fullCommand = strtpl.MustEval(`openssl x509 -req -in "{{.FilenameCsr}}" -CA "{{.FilenameCAPem}}" -CAkey "{{.FilenameCAKey}}" -CAcreateserial -out "{{.FilenameCrt}}" -days {{.Days}} -sha256 -extfile "{{.FilenameConfig}}"`, certificateRequest)
	exec(c, fullCommand)

	// remove csr file
	err = os.Remove(certificateRequest.FilenameCsr)
	if err != nil {
		fmt.Println(err)
	}

	// openssl x509 -text -noout -in certs/zgo.dev.crt (show certificate)
	fullCommand = strtpl.MustEval(`openssl x509 -text -noout -in "{{.FilenameCrt}}"`, certificateRequest)
	execOpen(c, fullCommand)

	// how to integrate to Apache
	hintCertApache(certificateRequest)
}
