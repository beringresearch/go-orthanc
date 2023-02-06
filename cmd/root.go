package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dicomweb-cli",
	Short: "CLI for interacting with DICOMweb services",
	Long: `A CLI tool for interacting with DICOMweb RESTful services.

Supports:
- Upload of DICOM files using STOW-RS
- Download of DICOM files using WADO-RS
- Querying of DICOM files using QUDO-RS
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
