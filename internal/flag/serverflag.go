package allflags

// import (
// 	"flag"
// 	"os"
// 	"strings"
// )

// // parseFlags parses the flags and sets the serverAddress variable
// var ServerAddress = flag.String("a", "localhost:8080", "HTTP server address")

// func ParseFlags() {
// 	flag.Parse()

// 	// Check if ServerAddress starts with http://
// 	if strings.HasPrefix(*ServerAddress, "http://") {
// 		// Remove http:// from ServerAddress
// 		*ServerAddress = strings.TrimPrefix(*ServerAddress, "http://")
// 	}

// 	// fmt.Println("from flag:", *ServerAddress)

// 	*ServerAddress = getenvWithDefault("ADDRESS", *ServerAddress)

// 	// fmt.Println("from env:", *ServerAddress)
// }

// func getenvWithDefault(name, defaultValue string) string {
// 	val := os.Getenv(name)
// 	if val == "" {
// 		val = defaultValue
// 	}

// 	return val
// }
