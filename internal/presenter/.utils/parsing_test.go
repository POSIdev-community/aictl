package _utils

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

func TestParseUUIDs(t *testing.T) {
	id1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	id2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("valid", func(t *testing.T) {
		got, err := ParseUUIDs([]string{id1.String(), id2.String()})
		require.NoError(t, err)
		require.Equal(t, []uuid.UUID{id1, id2}, got)
	})

	t.Run("empty", func(t *testing.T) {
		got, err := ParseUUIDs(nil)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := ParseUUIDs([]string{"not-a-uuid"})
		require.Error(t, err)

		var fieldErr *validation.FieldError
		require.ErrorAs(t, err, &fieldErr)
		require.Equal(t, "not-a-uuid", fieldErr.Field)
	})
}

func TestReadArgsFromReader(t *testing.T) {
	t.Run("passthrough_without_dash", func(t *testing.T) {
		args := []string{"a", "b"}
		require.Equal(t, args, ReadArgsFromReader(args, strings.NewReader("ignored")))
	})

	t.Run("reads_stdin_when_dash", func(t *testing.T) {
		got := ReadArgsFromReader([]string{"-"}, strings.NewReader("line1\nline2\n"))
		require.Equal(t, []string{"line1\nline2"}, got)
	})
}
