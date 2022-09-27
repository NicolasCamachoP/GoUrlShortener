package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name             string
		want             *Config
		ephimeralEnvVars map[string]string
		wantErr          bool
	}{
		{
			name: "TestGetConfig - Ok",
			want: &Config{
				"TEST",
				123,
				"TEST",
				"TEST",
				"TEST",
				"TEST",
				123,
				"TEST",
			},
			ephimeralEnvVars: map[string]string{
				"DB_HOST":            "TEST",
				"DB_PORT_NUMBER":     "123",
				"DB_PASSWORD":        "TEST",
				"DB_USERNAME":        "TEST",
				"DB_DATABASE":        "TEST",
				"DB_COLLECTION":      "TEST",
				"SERVER_PORT_NUMBER": "123",
				"SERVER_BASEPATH":    "TEST",
			},
			wantErr: false,
		},
		{
			name: "TestGetConfig - incomplete env vars",
			want: nil,
			ephimeralEnvVars: map[string]string{
				"DB_PORT_NUMBER":     "123",
				"DB_PASSWORD":        "TEST",
				"SERVER_PORT_NUMBER": "123",
				"SERVER_BASEPATH":    "TEST",
			},
			wantErr: true,
		},
		{
			name: "Wrong env var type",
			want: nil,
			ephimeralEnvVars: map[string]string{
				"DB_DOMAIN":          "TEST",
				"DB_PORT_NUMBER":     "TEST",
				"DB_PASSWORD":        "TEST",
				"SERVER_PORT_NUMBER": "123",
				"SERVER_BASEPATH":    "TEST",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.ephimeralEnvVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}
			got, err := GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
			for key := range tt.ephimeralEnvVars {
				os.Unsetenv(key)
			}
		})
	}
}
