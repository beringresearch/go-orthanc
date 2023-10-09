package cmd

import (
	"log"

	"github.com/dicomweb-cli/core"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var derivedCmd = &cobra.Command{
	Use:   "derived <dicom> <image>",
	Short: "Generate a dervied DICOM",
	Long:  `Generate a new dicom file derived from the input DICOM file using the provided image.`,
	Run: func(cmd *cobra.Command, args []string) {
		outPath := cmd.Flag("out").Value.String()
		err := core.CreateDerivedImage(args[0], args[1], outPath)
		if err != nil {
			log.Fatal(err)
		}
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	derivedCmd.Flags().StringP("out", "o", "derived.dcm", "path to output the generated DICOM")
	rootCmd.AddCommand(derivedCmd)
}
