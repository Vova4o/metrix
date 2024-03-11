package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"TestEmptyFlags", []string{}, "localhost:8080"},
		{"TestSingleFlag", []string{"-a=localhost:8090"}, "localhost:8090"},
		{"TestSingleFlagEnv", []string{"-a=http://localhost:8095"}, "localhost:8095"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append(os.Args, tt.args...)
			parseFlags()
			assert.Equal(t, *ServerAddress, tt.expected)
			os.Args = os.Args[:len(os.Args)-len(tt.args)]
		})
	}
}
