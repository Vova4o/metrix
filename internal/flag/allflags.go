package allflags

import (
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// Set the default values
	viper.SetDefault("ServerAddress", "localhost:8080")
	viper.SetDefault("ReportInterval", 10)
	viper.SetDefault("PollInterval", 2)

	// Define the flags and bind them to viper
	pflag.StringP("ServerAddress", "a", viper.GetString("ServerAddress"), "HTTP server network address")
	pflag.IntP("ReportInterval", "r", viper.GetInt("ReportInterval"), "Interval between fetching reportable metrics in seconds")
	pflag.IntP("PollInterval", "p", viper.GetInt("PollInterval"), "Interval between polling metrics in seconds")

	// Parse the command-line flags
	pflag.Parse()

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
	if err := viper.BindPFlag(flagName, pflag.Lookup(flagName)); err != nil {
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
