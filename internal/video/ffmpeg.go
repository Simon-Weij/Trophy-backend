package video

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

func GetContainer(filePath string) (string, error) {
	out, err := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", filePath).Output()
	if err != nil {
		return "", err
	}
	var result struct {
		Format struct {
			FormatLongName string `json:"format_long_name"`
		} `json:"format"`
	}
	json.Unmarshal(out, &result)
	return result.Format.FormatLongName, nil
}

func GetCodec(filePath string) (string, error) {
	out, err := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", filePath).Output()
	if err != nil {
		return "", err
	}
	var result struct {
		Streams []struct {
			CodecName string `json:"codec_name"`
			CodecType string `json:"codec_type"`
		} `json:"streams"`
	}
	json.Unmarshal(out, &result)
	for _, s := range result.Streams {
		if s.CodecType == "video" {
			return s.CodecName, nil
		}
	}
	return "", fmt.Errorf("no video stream found")
}

func TranscodeToWebm(inputPath string, outputPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inputPath, "-c:v", "libvpx-vp9", "-y", outputPath)

	return cmd.Run()
}
