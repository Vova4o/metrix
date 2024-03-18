package allflags

import (
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
	viper.BindPFlag("ServerAddress", pflag.Lookup("ServerAddress"))
	viper.BindPFlag("ReportInterval", pflag.Lookup("ReportInterval"))
	viper.BindPFlag("PollInterval", pflag.Lookup("PollInterval"))

	// Set the environment variable names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("ServerAddress", "ADDRESS")
	viper.BindEnv("ReportInterval", "REPORT_INTERVAL")
	viper.BindEnv("PollInterval", "POLL_INTERVAL")

	// Read the environment variables
	viper.AutomaticEnv()
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
