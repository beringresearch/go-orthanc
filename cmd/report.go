package cmd

import (
	"log"

	"github.com/dicomweb-cli/core"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a dervied DICOM",
	Long:  `Generate a new dicom file derived from the input DICOM file`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.CreateStructuredReport(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
