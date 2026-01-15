package _utils

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	uri         string
	token       string
	tlsSkip     bool
	verboseFlag bool
	logPath     string
)

func AddConnectionPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&uri, "uri", "u", "", "AI server uri")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "AI server access token")
	cmd.PersistentFlags().BoolVar(&tlsSkip, "tls-skip", false, "Skip certificate verification")
	cmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringVarP(&logPath, "log-path", "l", "", "log file path")
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
			return fmt.Errorf("update context: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate cfg: %w", err)
		}

		return nil
	}
}

type RunE = func(cmd *cobra.Command, args []string) error

func ChainRunE(funcs ...RunE) RunE {
	return func(cmd *cobra.Command, args []string) error {
		if funcs != nil {
			for _, f := range funcs {
				if err := f(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func InitializeLogger(cmd *cobra.Command, _ []string) error {
	l, _ := logger.NewLogger(verboseFlag, logPath)
	ctx := logger.ContextWithLogger(cmd.Context(), l)

	cmd.SetContext(ctx)

	return nil
}
