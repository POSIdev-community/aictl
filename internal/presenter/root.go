package presenter

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/context"
	"github.com/POSIdev-community/aictl/internal/presenter/create"
	"github.com/POSIdev-community/aictl/internal/presenter/delete"
	"github.com/POSIdev-community/aictl/internal/presenter/get"
	"github.com/POSIdev-community/aictl/internal/presenter/scan"
	"github.com/POSIdev-community/aictl/internal/presenter/set"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewRootCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		verbose bool
		logPath string
	)

	rootCmd := &cobra.Command{
		Use:               "aictl",
		Short:             "Application Inspector ConTroL tool",
		Long:              `aictl - api клиент для Application Inspector`,
		PersistentPreRunE: initialize(&verbose, &logPath),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					return fmt.Errorf("get help: %w", err)
				}
				return nil
			}

			// TODO: add old aisa behavior
			fmt.Println("TODO: add old aisa behavior ")

			return nil
		},
	}

	rootCmd.AddCommand(context.NewContextCmd(cfg, depsContainer))
	rootCmd.AddCommand(create.NewCreateCmd(cfg, depsContainer))
	rootCmd.AddCommand(delete.NewDeleteCmd(cfg, depsContainer))
	rootCmd.AddCommand(get.NewGetCmd(cfg, depsContainer))
	rootCmd.AddCommand(scan.NewScanCmd(cfg, depsContainer))
	rootCmd.AddCommand(set.NewSetCmd(cfg, depsContainer))
	rootCmd.AddCommand(update.NewUpdateCmd(cfg, depsContainer))

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&logPath, "log-path", "l", "", "log file path")

	return rootCmd
}

func initialize(verbose *bool, logPath *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true

		if !*verbose && *logPath == "" {
			l := zap.NewNop()
			ctx := logger.ContextWithLogger(cmd.Context(), l)
			cmd.SetContext(ctx)

			return nil
		}

		cfg := zap.NewDevelopmentConfig()
		cfg.OutputPaths = []string{}
		cfg.ErrorOutputPaths = []string{}

		if *logPath != "" {
			cfg.OutputPaths = append(cfg.OutputPaths, *logPath)
			cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, *logPath)
		}
		if *verbose {
			cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
			cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, "stderr")
		}

		if err := zap.RegisterEncoder("verbose", logger.NewVerboseEncoder); err != nil {
			return fmt.Errorf("register logger encoder failed: %v", err)
		}

		cfg.Encoding = "verbose"
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.EncoderConfig.ConsoleSeparator = "  "

		zapLevel := zapcore.DebugLevel
		cfg.Level = zap.NewAtomicLevelAt(zapLevel)

		l, err := cfg.Build()
		if err != nil {
			return fmt.Errorf("create logger: %v", err)
		}

		zap.ReplaceGlobals(l)
		zap.RedirectStdLog(l)

		ctx := cmd.Context()
		ctx = logger.ContextWithLogger(ctx, l)

		cmd.SetContext(ctx)

		return nil
	}
}
