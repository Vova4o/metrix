package allflags

// import (
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestParseFlags(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		args     []string
// 		expected string
// 	}{
// 		{"TestEmptyFlags", []string{}, "localhost:8080"},
// 		{"TestSingleFlag", []string{"-a=localhost:8090"}, "localhost:8090"},
// 		{"TestSingleFlagEnv", []string{"-a=http://localhost:8095"}, "localhost:8095"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Save current os.Args
// 			oldArgs := os.Args
// 			// Set os.Args for the test
// 			os.Args = append([]string{"cmd"}, tt.args...)
// 			ParseFlags()
// 			assert.Equal(t, *ServerAddress, tt.expected)
// 			// Restore original os.Args
// 			os.Args = oldArgs
// 		})
// 	}
// }
