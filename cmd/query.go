package cmd

import (
	"fmt"
	"log"

	"github.com/dicomweb-cli/core"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query for a DICOM study metadata",
	Long: `Query a DICOM study metadata by StudyInstanceUID. The result will be saved to a JSON file.

Example:
dicomweb-cli query 1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6 --output test.json`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("nothing to query")
			return
		}

		url := cmd.Flag("url").Value.String()
		destPath := cmd.Flag("output").Value.String()

		for _, arg := range args {
			err := core.QueryStudy(arg, url, destPath)
			if err != nil {
				log.Fatalf("failed to query %s:\n%s", arg, err)
			}
			fmt.Println("Metadata saved: ", arg)
		}
	},
	PersistentPreRun: ensureConfigured,
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.PersistentFlags().String("url", "", "DICOMweb server URL. If provided overrides default config.")
	queryCmd.PersistentFlags().String("output", "", "Destination path for downloaded metadata.")
}
