package delete

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeDeleteProjectsUC struct {
	called int
	ids    []uuid.UUID
}

func (f *fakeDeleteProjectsUC) Execute(_ context.Context, projectIds []uuid.UUID) error {
	f.called++
	f.ids = projectIds
	return nil
}

func TestDeleteProjectsCommand(t *testing.T) {
	id1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	id2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	cfg := cmdtest.MustCfg(t)

	t.Run("requires_args", func(t *testing.T) {
		uc := &fakeDeleteProjectsUC{}
		root := NewDeleteCmd(cfg, NewDeleteProjectsCommand(uc))
		require.Error(t, cmdtest.Execute(t, root.Command, "projects"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("invalid_uuid", func(t *testing.T) {
		uc := &fakeDeleteProjectsUC{}
		root := NewDeleteCmd(cfg, NewDeleteProjectsCommand(uc))
		require.Error(t, cmdtest.Execute(t, root.Command, "projects", "bad-id"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("passes_ids", func(t *testing.T) {
		uc := &fakeDeleteProjectsUC{}
		root := NewDeleteCmd(cfg, NewDeleteProjectsCommand(uc))
		require.NoError(t, cmdtest.Execute(t, root.Command, "projects", id1.String(), id2.String()))
		require.Equal(t, []uuid.UUID{id1, id2}, uc.ids)
	})

	t.Run("stdin_dash", func(t *testing.T) {
		uc := &fakeDeleteProjectsUC{}
		root := NewDeleteCmd(cfg, NewDeleteProjectsCommand(uc))
		cmdtest.WithStdin(t, id1.String()+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "projects", "-"))
		})
		require.Equal(t, []uuid.UUID{id1}, uc.ids)
	})
}
