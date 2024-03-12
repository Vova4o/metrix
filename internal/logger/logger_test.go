package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid file name",
			args: args{
				name: "test.log",
			},
			wantErr: false,
		},
		{
			name: "Invalid file name",
			args: args{
				name: "/nonexistent/test.log",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Logger(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				defer got.Close()
			}
		})
	}
}
