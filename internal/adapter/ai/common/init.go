package common

import (
	"context"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

type InitState struct {
	Version      version.Version
	VersionKnown bool
}

type Initializer interface {
	TryInitialize(
		ctx context.Context,
		base *BaseClient,
		cfg *config.Config,
		state InitState,
	) (client ClientAi, nextState InitState, matched bool, err error)
}

type versionedInitializer interface {
	Initialize(ctx context.Context, cfg *config.Config) error
	GetVersion(ctx context.Context) (version.Version, error)
}

type VersionRangeInitializer struct {
	MinVersion version.Version
	MaxVersion version.Version
	NewClient  func(*BaseClient) ClientAi
}

func MatchesVersionRange(ver, min, max version.Version) bool {
	return !ver.Less(min) && ver.Less(max)
}

func (i VersionRangeInitializer) TryInitialize(
	ctx context.Context,
	base *BaseClient,
	cfg *config.Config,
	state InitState,
) (ClientAi, InitState, bool, error) {
	client := i.NewClient(base)
	initializable, ok := client.(versionedInitializer)
	if !ok {
		return nil, state, false, nil
	}

	if err := initializable.Initialize(ctx, cfg); err != nil {
		// Soft failure: caller may try the next versioned client, but must surface the error
		// if none match (otherwise TLS/network problems become "no compatible client").
		return nil, state, false, err
	}

	state, err := EnsureVersion(ctx, state, initializable.GetVersion)
	if err != nil {
		return nil, state, false, err
	}

	if !MatchesVersionRange(state.Version, i.MinVersion, i.MaxVersion) {
		return nil, state, false, nil
	}

	return client, state, true, nil
}

func EnsureVersion(
	ctx context.Context,
	state InitState,
	getVersion func(context.Context) (version.Version, error),
) (InitState, error) {
	if state.VersionKnown {
		return state, nil
	}

	ver, err := getVersion(ctx)
	if err != nil {
		return state, err
	}

	return InitState{Version: ver, VersionKnown: true}, nil
}
