package serverflags

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
	flags.IntP("StoreInterval", "i", 300, "Interval in seconds to store the current server readings to disk")
	flags.StringP("FileStoragePath", "f", "/tmp/metrics-db.json", "Full filename where current values are saved")
	flags.BoolP("Restore", "r", true, "Whether to load previously saved values from the specified file at server startup")
	flags.StringP("DatabaseDSN", "d", "postgres://metrix:metrix@localhost:5432/metrix?sslmode=disable", "Database DSN")

	// Parse the command-line flags
	flags.Parse(os.Args[1:])

	// Bind the flags to viper
	bindFlagToViper("ServerAddress")
	bindFlagToViper("StoreInterval")
	bindFlagToViper("FileStoragePath")
	bindFlagToViper("Restore")
	bindFlagToViper("DatabaseDSN")

	// Set the environment variable names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnvToViper("ServerAddress", "ADDRESS")
	bindEnvToViper("StoreInterval", "STORE_INTERVAL")
	bindEnvToViper("FileStoragePath", "FILE_STORAGE_PATH")
	bindEnvToViper("Restore", "RESTORE")
	bindEnvToViper("DatabaseDSN", "DATABASE_DSN")

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

func ServerAddress() string {
	return viper.GetString("ServerAddress")
}

func StoreInterval() int {
	return viper.GetInt("StoreInterval")
}

func FileStoragePath() string {
	path := viper.GetString("FileStoragePath")
	if path == "=" {
		return ""
	}
	return path
}

func Restore() bool {
	return viper.GetBool("Restore")
}

func DatabaseDSN() string {
	return viper.GetString("DatabaseDSN")
}
