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
		outPath := cmd.Flag("out").Value.String()
		err := core.CreateStructuredReport(args[0], outPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	reportCmd.Flags().StringP("out", "o", "derived-sr.dcm", "path to output the generated DICOM SR")
	rootCmd.AddCommand(reportCmd)
}
