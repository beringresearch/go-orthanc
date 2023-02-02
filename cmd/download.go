package cmd

import (
	"fmt"
	"log"

	"github.com/dicomweb-cli/core"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <studyId> [<studyId>...]",
	Short: "Download a DICOM study to disk",
	Long: `Download a DICOM study from the DICOMweb server to disk using its StudyInstanceUID.

Example:

dicomweb-cli download 1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6 --output test.dcm
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("nothing to download")
			return
		}

		url := cmd.Flag("url").Value.String()
		destPath := cmd.Flag("output").Value.String()

		for _, arg := range args {
			err := core.DownloadStudy(arg, destPath, url)
			if err != nil {
				log.Fatalf("failed to download %s: %s\n", arg, err)
			}
			fmt.Println("Downloaded ", arg)
		}
	},
	PersistentPreRun: ensureConfigured,
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.PersistentFlags().String("url", "", "DICOMweb server URL. If provided overrides default config.")
	downloadCmd.PersistentFlags().String("output", "", "Destination path for downloaded study.")
}
