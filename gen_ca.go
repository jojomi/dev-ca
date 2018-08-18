package main

import (
	"os"
	"path/filepath"

	"github.com/jojomi/go-script/interview"

	"github.com/jojomi/go-script"
	"github.com/jojomi/strtpl"
	"github.com/spf13/cobra"
)

var (
	flagGenCARSASize int
	flagGenCADays    int
	flagGenCASubj    string
	flagGenCAName    string

	caFolder = "ca"
)

func genCACmd() *cobra.Command {
	c := cobra.Command{
		Use:   "gen-ca",
		Short: "Generate a Certificate Authority",
		Run:   genCAHandler,
	}
	f := c.PersistentFlags()
	f.IntVarP(&flagGenCARSASize, "rsa-size", "r", 2048, "RSA key size")
	f.IntVarP(&flagGenCADays, "days", "d", 365*10, "days of validity")
	f.StringVarP(&flagGenCASubj, "subj", "s", "", "value for openssl subj flag")
	f.StringVarP(&flagGenCAName, "name", "n", "", "name of the CA")
	return &c
}

func genCAHandler(cmd *cobra.Command, args []string) {
	c := script.NewContext()
	c.MustCommandExist("openssl")
	outputFolder := caFolder
	err := c.EnsureDirExists(outputFolder, 0700)
	checkFail(err)

	// openssl genrsa -out rootCA.key 2048
	caRequest := &CARequest{
		RSAKeySize:  flagGenCARSASize,
		Days:        flagGenCADays,
		Subj:        flagGenCASubj,
		Name:        flagGenCAName,
		FilenameKey: filepath.Join(outputFolder, "rootCA.key"),
		FilenamePem: filepath.Join(outputFolder, "rootCA.pem"),
	}

	if c.FileExists(caRequest.FilenameKey) || c.FileExists(caRequest.FilenamePem) {
		overwrite, err := interview.Confirm("Files exist, overwrite?", false)
		if err != nil || !overwrite {
			os.Exit(1)
		}
	}

	fullCommand := strtpl.MustEval(`openssl genrsa -out "{{ .FilenameKey }}" {{.RSAKeySize}}`, caRequest)
	exec(c, fullCommand)

	// openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -out rootCA.pem
	if caRequest.Name == "" {
		caRequest.Name = "Dev CA"
	}
	fullCommand = strtpl.MustEval(`openssl req -x509 -new -nodes -key "{{.FilenameKey}}" -sha256 -days {{.Days}} -out "{{.FilenamePem}}"
	{{- if .Subj }} -subj "{{.Subj}}"
	{{- else if .Name }} -subj "/C=/ST=/L=/O=/OU=/CN={{.Name}}"
	{{- end -}}`, caRequest)
	exec(c, fullCommand)

	// openssl x509 -text -noout -in certs/zgo.dev.crt (show certificate)
	fullCommand = strtpl.MustEval(`openssl x509 -text -noout -in "{{.FilenamePem}}"`, caRequest)
	execOpen(c, fullCommand)

	hintCAFirefox(caRequest)
	hintCAMacOSX(caRequest)
}
