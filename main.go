package main

import (
	"fmt"
	"os"

	script "github.com/jojomi/go-script"
	"github.com/jojomi/go-script/print"
	"github.com/spf13/cobra"
)

var (
	binaryName = "dev-ca"

	flagVerbose bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use: binaryName,
	}
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(genCACmd(), genCertCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func exec(c *script.Context, fullCommand string) {
	if flagVerbose {
		print.Boldln(fullCommand)
	}
	cm, params := script.SplitCommand(fullCommand)
	var execFunc func(string, ...string) (*script.ProcessResult, error)
	if flagVerbose {
		execFunc = c.ExecuteDebug
	} else {
		execFunc = c.ExecuteFullySilent
	}
	pr, err := execFunc(cm, params...)
	if err != nil || !pr.Successful() {
		fmt.Printf("Command failed (%s): %s\n", err, fullCommand)
		os.Exit(1)
	}
	if flagVerbose {
		fmt.Println()
		fmt.Println()
	}
}

func execOpen(c *script.Context, fullCommand string) {
	if flagVerbose {
		print.Boldln(fullCommand)
	}
	cm, params := script.SplitCommand(fullCommand)
	pr, err := c.ExecuteDebug(cm, params...)
	if err != nil || !pr.Successful() {
		fmt.Printf("Command failed (%s): %s\n", err, fullCommand)
		os.Exit(1)
	}
	if flagVerbose {
		fmt.Println()
		fmt.Println()
	}
}
