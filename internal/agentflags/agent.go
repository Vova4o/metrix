package agentflags

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var flags = pflag.NewFlagSet("flags", pflag.ExitOnError)

func init() {
	// Define the flags and bind them to viper
	flags.StringP("ServerAddress", "a", "localhost:8080", "HTTP server network address")
	flags.IntP("ReportInterval", "r", 10, "Interval between fetching reportable metrics in seconds")
	flags.IntP("PollInterval", "p", 2, "Interval between polling metrics in seconds")

	// Parse the command-line flags
	flags.Parse(os.Args[1:])

	// Bind the flags to viper
	bindFlagToViper("ServerAddress")
	bindFlagToViper("ReportInterval")
	bindFlagToViper("PollInterval")

	// Set the environment variable names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnvToViper("ServerAddress", "ADDRESS")
	bindEnvToViper("ReportInterval", "REPORT_INTERVAL")
	bindEnvToViper("PollInterval", "POLL_INTERVAL")

	// Read the environment variables
	viper.AutomaticEnv()
}

func bindFlagToViper(flagName string) {
	if err := viper.BindPFlag(flagName, flags.Lookup(flagName)); err != nil {
		log.Println(err)
	}
}

func bindEnvToViper(viperKey, envKey string) {
	if err := viper.BindEnv(viperKey, envKey); err != nil {
		log.Println(err)
	}
}

func GetServerAddress() string {
	return viper.GetString("ServerAddress")
}

func GetReportInterval() int {
	return viper.GetInt("ReportInterval")
}

func GetPollInterval() int {
	return viper.GetInt("PollInterval")
}
