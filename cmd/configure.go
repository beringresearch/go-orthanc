package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/dicomweb-cli/core"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the DICOMweb service settings",
	Long: `Configure the DICOMweb service settings through CLI.
Either a valid URL or an empty string can be passed.
If an empty string is passed, URL will be configured interactively.
When prompted enter the URL of the DICOMweb server in full, including scheme, port, etc.
Example:

dicomweb-cli configure
> Enter the DICOMweb server URL: http://localhost:8042/dicom-web
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			createConfig("")
		}
		createConfig(args[0])
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

// ensureConfigured makes sure the user has configured a DICOMweb server before running commands.
// It checks for an existing config and makes the user create one if it does not exist.
func ensureConfigured(cmd *cobra.Command, args []string) {
	// If URL not explicitly provided in flag need to check for config.
	if flag := cmd.Flag("url"); flag == nil || flag.Value.String() == "" {
		if !core.ConfigExists() {
			fmt.Println("No DICOMweb server config found, creating one...")
			fmt.Println()
			createConfig("")
		}

	}
}

// createConfig prompts the user for DICOMweb server settings and saves
// the config settings in the home directory
func createConfig(serverUrl string) {
	if serverUrl == "" {
		fmt.Print("> Enter the DICOMweb server URL: ")
		fmt.Scanln(&serverUrl)
	}

	err := validateURL(serverUrl)
	if err != nil {
		log.Fatal(err)
	}

	config := core.DICOMWebServer{
		QIDOEndpoint: serverUrl,
		WADOEndpoint: serverUrl,
		STOWEndpoint: serverUrl,
	}
	err = core.SaveServerConfig(config)
	if err != nil {
		log.Fatalln("failed to save config file: ", err)
	}

	fmt.Println("Config saved")
}

// validateURL catches invalid URLs and returns an error detailing the problem
func validateURL(urlString string) error {
	parsedUrl, err := url.ParseRequestURI(urlString)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme == "" {
		return errors.New("no protocol scheme provided in URL (e.g. HTTP)")
	}
	// If host is empty, the host is usually being parsed as the URL scheme
	if parsedUrl.Host == "" {
		return errors.New("no protocol scheme provided in URL (e.g. HTTP)")
	}
	if parsedUrl.Port() == "" {
		return errors.New("no port provided in URL (e.g. :8042)")
	}
	if parsedUrl.Path == "" {
		return errors.New("no path provided - if this is purposeful append a trailing '/' to the URL")
	}

	return nil
}
