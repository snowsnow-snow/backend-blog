package media

import (
	"backend-blog/internal/model/entity"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

// VideoInfo extracted from video file
type VideoInfo struct {
	Asset    *entity.MediaAsset
	Metadata *entity.VideoMetadata
}

// ReadVideoInfo extracts metadata from a video file using ffprobe
func ReadVideoInfo(fullPath string) (*VideoInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use binary from system path by default
	data, err := ffprobe.ProbeURL(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	if len(data.Streams) == 0 {
		return nil, fmt.Errorf("no streams found in video")
	}

	// Find the first video stream
	var videoStream *ffprobe.Stream
	for _, s := range data.Streams {
		if s.CodecType == "video" {
			videoStream = s
			break
		}
	}

	if videoStream == nil {
		return nil, fmt.Errorf("no video stream found")
	}

	asset := &entity.MediaAsset{
		Width:    videoStream.Width,
		Height:   videoStream.Height,
		Duration: int(data.Format.Duration().Seconds()),
	}

	metadata := &entity.VideoMetadata{
		Codec:     videoStream.CodecName,
		Bitrate:   parseBitrate(videoStream.BitRate),
		FrameRate: parseFrameRate(videoStream.RFrameRate),
	}

	// Try to extract device info
	extractDeviceInfo(data, asset)

	return &VideoInfo{
		Asset:    asset,
		Metadata: metadata,
	}, nil
}

func parseBitrate(br string) int {
	if br == "" {
		return 0
	}
	val, _ := strconv.Atoi(br)
	return val
}

func parseFrameRate(fr string) int {
	if fr == "" {
		return 0
	}
	parts := strings.Split(fr, "/")
	if len(parts) == 1 {
		f, _ := strconv.ParseFloat(parts[0], 64)
		return int(f)
	}
	if len(parts) == 2 {
		num, _ := strconv.ParseFloat(parts[0], 64)
		den, _ := strconv.ParseFloat(parts[1], 64)
		if den != 0 {
			return int(num / den)
		}
	}
	return 0
}

func extractDeviceInfo(data *ffprobe.ProbeData, asset *entity.MediaAsset) {
	// Common tags for Make
	makeTags := []string{"com.apple.quicktime.make", "make", "Make", "manufacturer"}
	modelTags := []string{"com.apple.quicktime.model", "model", "Model"}

	for _, tag := range makeTags {
		if val, err := data.Format.TagList.GetString(tag); err == nil && val != "" {
			asset.DeviceMake = val
			break
		}
	}

	for _, tag := range modelTags {
		if val, err := data.Format.TagList.GetString(tag); err == nil && val != "" {
			asset.DeviceModel = val
			break
		}
	}

	// Fujifilm specific in 'comment' tag
	if asset.DeviceMake == "" {
		if comment, err := data.Format.TagList.GetString("comment"); err == nil && strings.HasPrefix(comment, "FUJIFILM") {
			parts := strings.Split(comment, " ")
			asset.DeviceMake = "FUJIFILM"
			if len(parts) > 1 {
				asset.DeviceModel = parts[len(parts)-1]
			}
		}
	}
}
