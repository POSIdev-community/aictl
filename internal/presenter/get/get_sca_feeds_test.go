package get

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeListScaFeedsUC struct {
	called   int
	statuses []scafeeds.Status
}

func (f *fakeListScaFeedsUC) Execute(_ context.Context, statuses []scafeeds.Status) error {
	f.called++
	f.statuses = statuses
	return nil
}

type fakeDownloadScaFeedsUC struct {
	called  int
	version string
	outPath string
}

func (f *fakeDownloadScaFeedsUC) Execute(_ context.Context, version, outPath string) error {
	f.called++
	f.version = version
	f.outPath = outPath
	return nil
}

// cmdtest.Execute swaps os.Stdout/Stderr — keep subtests sequential to avoid races.
func TestGetScaFeedsCmd(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		listUC := &fakeListScaFeedsUC{}
		dlUC := &fakeDownloadScaFeedsUC{}
		cmd := NewGetScaFeedsCmd(listUC, dlUC)
		require.NoError(t, cmdtest.Execute(t, cmd.Command))
		require.Equal(t, 1, listUC.called)
		require.Equal(t, 0, dlUC.called)
	})

	t.Run("list_with_status", func(t *testing.T) {
		listUC := &fakeListScaFeedsUC{}
		dlUC := &fakeDownloadScaFeedsUC{}
		cmd := NewGetScaFeedsCmd(listUC, dlUC)
		require.NoError(t, cmdtest.Execute(t, cmd.Command, "--status", "current", "--status", "active"))
		require.Equal(t, []scafeeds.Status{scafeeds.StatusCurrent, scafeeds.StatusActive}, listUC.statuses)
	})

	t.Run("download", func(t *testing.T) {
		listUC := &fakeListScaFeedsUC{}
		dlUC := &fakeDownloadScaFeedsUC{}
		cmd := NewGetScaFeedsCmd(listUC, dlUC)
		require.NoError(t, cmdtest.Execute(t, cmd.Command, "1.2.3", "-o", "./out.zip"))
		require.Equal(t, 1, dlUC.called)
		require.Equal(t, "1.2.3", dlUC.version)
		require.Equal(t, "./out.zip", dlUC.outPath)
		require.Equal(t, 0, listUC.called)
	})

	t.Run("download_without_output", func(t *testing.T) {
		dlUC := &fakeDownloadScaFeedsUC{}
		cmd := NewGetScaFeedsCmd(&fakeListScaFeedsUC{}, dlUC)
		require.NoError(t, cmdtest.Execute(t, cmd.Command, "2.0.0"))
		require.Equal(t, 1, dlUC.called)
		require.Equal(t, "2.0.0", dlUC.version)
		require.Equal(t, "", dlUC.outPath)
	})

	t.Run("status_with_version_error", func(t *testing.T) {
		cmd := NewGetScaFeedsCmd(&fakeListScaFeedsUC{}, &fakeDownloadScaFeedsUC{})
		require.Error(t, cmdtest.Execute(t, cmd.Command, "1.2.3", "--status", "current"))
	})

	t.Run("unknown_status", func(t *testing.T) {
		cmd := NewGetScaFeedsCmd(&fakeListScaFeedsUC{}, &fakeDownloadScaFeedsUC{})
		require.Error(t, cmdtest.Execute(t, cmd.Command, "--status", "nope"))
	})
}
