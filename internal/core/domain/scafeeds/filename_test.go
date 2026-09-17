package scafeeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeFileName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1.2.3", SanitizeFileName("1.2.3"))
	assert.Equal(t, "a_b_c", SanitizeFileName("a/b\\c"))
	assert.Equal(t, "x_y", SanitizeFileName("x y"))
	assert.Equal(t, "", SanitizeFileName(""))
}

func TestDefaultArchiveName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "sca-feeds-1.2.3.zip", DefaultArchiveName("1.2.3"))
	assert.Equal(t, "sca-feeds-a_b.zip", DefaultArchiveName("a/b"))
}
