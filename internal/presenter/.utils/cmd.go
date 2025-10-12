package _utils

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

var (
	uri     string
	token   string
	tlsSkip bool
)

func AddConnectionPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&uri, "uri", "u", "", "AI server uri")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "AI server access token")
	cmd.PersistentFlags().BoolVar(&tlsSkip, "tls-skip", false, "Skip certificate verification")
}

func UpdateConnectionConfig(cfg *config.Config) error {
	if uri != "" {
		err := cfg.SetURI(uri)
		if err != nil {
			return fmt.Errorf("set uri error: %w", err)
		}
	}

	if token != "" {
		err := cfg.SetToken(token)
		if err != nil {
			return fmt.Errorf("set toker error: %w", err)
		}
	}

	if tlsSkip {
		cfg.SetTLSSkip(tlsSkip)
	}

	return nil
}

func UpdateConfig(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := UpdateConnectionConfig(cfg)
		if err != nil {
			return fmt.Errorf("could not update context: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate cfg error: %w", err)
		}

		return nil
	}
}
