package utility

import (
	"backend-blog/internal/model/entity"
	"github.com/adrium/goheif"
	"image/jpeg"
	"io"
	"os"
)

// ImgUtil 读取图片中的信息
type ImgUtil interface {
	// ReadExif 读取图片 EXIF 信息的通用方法
	ReadExif(path string, fileName string, fileType string) *entity.BlogImage
	// ReadFujiInfo 读取富士相机独有的参数
	ReadFujiInfo(imgInfo entity.BlogImage, exifData map[string]interface{})
}

// Skip Writer for exif writine
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

// ConvertHeicToJpg takes in an input file (of heic format) and converts
// it to a jpeg format, named as the output parameters.
func ConvertHeicToJpg(input, output string) error {

	fileInput, err := os.Open(input)
	if err != nil {
		return err
	}
	defer fileInput.Close()

	// Extract exif to add back in after conversion
	exif, err := goheif.ExtractExif(fileInput)
	if err != nil {
		return err
	}

	img, err := goheif.Decode(fileInput)
	if err != nil {
		return err
	}

	fileOutput, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fileOutput.Close()

	// Write both convert file + exif data back
	w, _ := newWriterExif(fileOutput, exif)
	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return err
	}
	return nil
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		appMarker := 0xe1
		markerLen := 2 + len(exif)
		marker := []byte{0xff, uint8(appMarker), uint8(markerLen >> 8), uint8(markerLen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}
