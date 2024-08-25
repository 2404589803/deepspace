//go:build endpoint_custom && !endpoint_sg
// +build endpoint_custom,!endpoint_sg

package main

import (
	"strings"

	"github.com/spf13/cobra"
)

var endpoint string

func init() {
	flags := DeepSpace.PersistentFlags()
	flags.StringVar(&endpoint, "endpoint", "https://api.deepseek.com", "API endpoint")
	cobra.OnInitialize(func() {
		if !flags.Changed("endpoint") && DeepConfig.Endpoint != "" {
			endpoint = MoonConfig.Endpoint
		}
		endpoint = strings.TrimSuffix(endpoint, "/")
	})
}
