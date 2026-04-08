package video

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetContainer_InvalidPath(t *testing.T) {
	container, err := GetContainer("/path/does/not/exist.webm")
	require.Error(t, err)
	require.Empty(t, container)
}

func TestGetCodec_InvalidPath(t *testing.T) {
	codec, err := GetCodec("/path/does/not/exist.webm")
	require.Error(t, err)
	require.Empty(t, codec)
}

func TestTranscodeToWebm_InvalidInput(t *testing.T) {
	err := TranscodeToWebm("/path/does/not/exist.mp4", "/tmp/out.webm")
	require.Error(t, err)
}
