package port

import "github.com/POSIdev-community/aictl/internal/core/domain/config"

type Config interface {
	ClearCurrentContext() error

	StoreContext(cfg *config.Config) error

	String(*config.Config) (string, error)
	StringJson(*config.Config) (string, error)
	StringYaml(*config.Config) (string, error)
}
