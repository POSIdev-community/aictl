package _utils

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

var (
	uri         string
	token       string
	tlsSkip     bool
	cacert      string
	verboseFlag bool
	debugFlag   bool
	logPath     string
)

func AddConnectionPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&uri, "uri", "u", "", "AI server URI (overrides context)")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "AI server access token (overrides context)")
	cmd.PersistentFlags().BoolVar(&tlsSkip, "tls-skip", false, "Skip TLS certificate verification")
	cmd.PersistentFlags().StringVar(&cacert, "cacert", "", "Path to PEM file with CA certificate(s) to trust (appended to system roots)")
	cmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output (operational details)")
	cmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "V", false, "Debug output (includes error chains)")
	cmd.PersistentFlags().StringVarP(&logPath, "log-path", "l", "", "Log file path")
}

func VerboseLevel() int {
	level := logger.LevelQuiet
	if verboseFlag {
		level = logger.LevelVerbose
	}
	if debugFlag {
		level = logger.LevelDebug
	}
	return level
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

	if tlsSkip && cacert != "" {
		return fmt.Errorf("cannot use 'cacert' together with 'tls-skip'")
	}

	if tlsSkip {
		if cfg.CACertPath() != "" {
			return fmt.Errorf("cannot use 'tls-skip' together with configured cacert")
		}
		cfg.SetTLSSkip(tlsSkip)
	}

	if cacert != "" {
		if err := cfg.SetCACertPath(cacert); err != nil {
			return fmt.Errorf("set cacert error: %w", err)
		}
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
	l, err := logger.NewLogger(VerboseLevel(), logPath)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}
	ctx := logger.ContextWithLogger(cmd.Context(), l)

	cmd.SetContext(ctx)

	return nil
}
