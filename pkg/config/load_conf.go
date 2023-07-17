package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	logging "github.com/sirupsen/logrus"

	"k8s-image-check-admission-controller/pkg/k8simageadmissioncontroller"
)

// I'm declaring as vars so I can test easier, I recommend declaring these as constants
var (
	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --port is bound to K8S_IMAGE_ADMISSION_CONTROLLER_PORT.
	envPrefix = "K8S_IMAGE_ADMISSION_CONTROLLER"
	// Replace hyphenated flag names with camelCase in the config file
	replaceHyphenWithCamelCase = false
)

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	// Store the result of binding cobra flags and viper config. In a
	// real application these would be data structures, most likely
	// custom structs per command. This is simplified for the demo app and is
	// not recommended that you use one-off variables. The point is that we
	// aren't retrieving the values directly from viper or flags, we read the values
	// from standard Go data structures.
	defaultConfigFilename := ""
	environment := ""

	hostname := ""
	port := 0
	certFile := ""
	keyFile := ""

	compressedImageSizeLimit := int64(0)

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "k8s-image-admission-controller",
		Short: "A validating webhook for K8S images",
		Long:  `A k8S dynamic admission webhook controller to check some caracteristics on container's images`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd, defaultConfigFilename)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			// out := cmd.OutOrStdout()
			// Setup logging
			if environment == "prd" || environment == "prod" || environment == "production" {
				logging.SetFormatter(&logging.JSONFormatter{})
			} else {
				// The TextFormatter is default, you don't actually have to do this.
				logging.SetFormatter(&logging.TextFormatter{})
				logging.Warn("Server in development mode, do not run in production!")
			}
			// Print the final resolved value from binding cobra flags and viper config
			logging.WithFields(logging.Fields{
				"environment":              environment,
				"config":                   defaultConfigFilename,
				"hostname":                 hostname,
				"port":                     port,
				"certFile":                 certFile,
				"keyFile":                  keyFile,
				"compressedImageSizeLimit": compressedImageSizeLimit,
			}).Info("Variables loaded from env, configuration files and cli args")

			k8simageadmissioncontroller.RunWebhookServer(hostname, port, certFile, keyFile, compressedImageSizeLimit)
		},
	}

	// Define cobra flags, the default value has the lowest (least significant) precedence
	rootCmd.Flags().StringVarP(&environment, "env", "e", "prd", "Environment")
	rootCmd.Flags().StringVarP(&defaultConfigFilename, "config", "c", "config/configuration", "Configuration file")

	rootCmd.Flags().StringVar(&hostname, "hostname", "localhost", "Hostname")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port")
	rootCmd.Flags().StringVar(&certFile, "cert", "cert/server.pem", "Certiticate file for TLS server")
	rootCmd.Flags().StringVar(&keyFile, "key", "cert/server.key", "Key file for TLS server")

	rootCmd.Flags().Int64VarP(&compressedImageSizeLimit, "image-size-limit", "l", 1_000_000_000, "Compressed image size limit")

	return rootCmd
}

func initializeConfig(cmd *cobra.Command, defaultConfigFilename string) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name
		// If using camelCase in the config file, replace hyphens with a camelCased string.
		// Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
		if replaceHyphenWithCamelCase {
			configName = strings.ReplaceAll(f.Name, "-", "")
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				logging.Errorf("Error binding flags to variables %v", err)
			}
		}
	})
}
