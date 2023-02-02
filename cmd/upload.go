package cmd

import (
	"fmt"
	"log"

	"github.com/dicomweb-cli/core"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload <filepath> [<filepath>...]",
	Short: "Upload a DICOM file from disk",
	Long: `Upload a DICOM file from disk to the DICOM server.

Example:
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("nothing to upload")
			return
		}

		url := cmd.Flag("url").Value.String()

		for _, arg := range args {
			err := core.UploadStudy(arg, url)
			if err != nil {
				log.Fatalf("failed to upload %s:\n%s", arg, err)
			}
			fmt.Println("Uploaded ", arg)
		}
	},
	PersistentPreRun: ensureConfigured,
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.PersistentFlags().String("url", "", "DICOMweb server URL. If provided overrides default config.")
}
