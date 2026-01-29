package media

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// CompressImage compresses an image to a maximum dimension and quality.
// It also handles automatic rotation based on EXIF orientation.
func CompressImage(srcPath, dstPath string) error {
	// 1. Open and decode image with orientation handling
	src, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// 2. Apply Auto-Orientation
	src = AutoOrient(src, srcPath)

	// 3. Determine dimensions after orientation
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 4. Max dimension (e.g., 1920px)
	const maxDim = 1920
	if width > maxDim || height > maxDim {
		if width > height {
			src = imaging.Resize(src, maxDim, 0, imaging.Lanczos)
		} else {
			src = imaging.Resize(src, 0, maxDim, imaging.Lanczos)
		}
	}

	// 5. Save optimized version
	err = imaging.Save(src, dstPath, imaging.JPEGQuality(80))
	if err != nil {
		return fmt.Errorf("failed to save compressed image: %w", err)
	}

	return nil
}

// AutoOrient reads EXIF orientation and rotates the image accordingly.
func AutoOrient(img image.Image, srcPath string) image.Image {
	f, err := os.Open(srcPath)
	if err != nil {
		return img
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return img
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return img
	}

	orient, err := tag.Int(0)
	if err != nil {
		return img
	}

	switch orient {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img) // Rotate 90 CW (Go imaging: Rotate270 is 90 CW in standard logic if 0 is top)
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img) // Rotate 270 CW (Rotate90 CCW)
	}

	return img
}

// CompressVideo uses ffmpeg to generate a compressed version of a video.
func CompressVideo(srcPath, dstPath string) error {
	// Check if ffmpeg is available
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		slog.Warn("ffmpeg not found, skipping video compression")
		return fmt.Errorf("ffmpeg not found")
	}

	// ffmpeg command for web-friendly compression
	// -i: input file
	// -vcodec libx264: H.264 video codec
	// -crf 28: Constant Rate Factor (0-51, 23 is default, 28-30 is good for web)
	// -preset faster: encoding speed vs compression ratio
	// -vf "scale='min(1080,iw)':-2": scale to max 1080p height, preserving aspect ratio (must be even height)
	// -acodec aac: AAC audio codec
	// -movflags +faststart: allows video to start playing before fully downloaded
	args := []string{
		"-i", srcPath,
		"-vcodec", "libx264",
		"-crf", "28",
		"-preset", "faster",
		"-vf", "scale='min(1080,iw)':-2",
		"-acodec", "aac",
		"-movflags", "+faststart",
		"-y", // overwrite output
		dstPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg compression failed: %w (output: %s)", err, string(output))
	}

	return nil
}

// GetCompressedPath generates a sibling path with a suffix for the compressed version.
func GetCompressedPath(originalPath string) string {
	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(originalPath, ext)
	return base + "_compressed" + ext
}
