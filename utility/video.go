package utility

import (
	"backend-blog/config"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"context"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var makes = [...]string{"Apple", "iPhone", "FUJIFILM", "qt"}

func ReadInfo(path string) (*entity.BlogVideo, error) {
	path = config.GlobalConfig.File.Path.System + path
	ffprobe.SetFFProbeBinPath("/Users/snowsnowsnow/MyBlog/ffprobe/ffprobe")
	// 打开要解析的音视频文件
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	fileReader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := ffprobe.ProbeReader(ctx, fileReader)
	if err != nil {
		return nil, err
	}
	blogVideo := entity.BlogVideo{}
	setBrand(data, &blogVideo)

	return &blogVideo, nil
}

func setBrand(probeData *ffprobe.ProbeData, blogVideo *entity.BlogVideo) {
	for tag := range probeData.Format.TagList {
		getString, _ := probeData.Format.TagList.GetString(tag)
		if getString == "" {
			return
		}
		getString = strings.Trim(getString, " ")
		for index := range makes {
			if strings.Contains(getString, makes[index]) {
				handle := getHandleByBrand(makes[index])
				if handle != nil {
					handle.SetVideoInfo(probeData, blogVideo)
				}
			}
		}
	}
}
func getHandleByBrand(brand string) VideoInfoHandle {
	switch brand {
	case "Apple":
	case "iPhone":
		return &Apple{}
	case "FUJIFILM":
		return &Fujifilm{}
	default:
		return &Universal{}
	}
	return nil
}

type VideoInfoHandle interface {
	SetVideoInfo(probeData *ffprobe.ProbeData, blogVideo *entity.BlogVideo)
}

type Apple struct{}
type Fujifilm struct{}
type Universal struct{}

func (r Apple) SetVideoInfo(probeData *ffprobe.ProbeData, blogVideo *entity.BlogVideo) {
	makeName, err := probeData.Format.TagList.GetString("com.apple.quicktime.make")
	if err != nil {
		log.Println(err)
	} else {
		blogVideo.Make = makeName
	}
	model, err := probeData.Format.TagList.GetString("com.apple.quicktime.model")
	if err != nil {
		log.Println(err)
	} else {
		blogVideo.Model = model
	}
}

func (r Fujifilm) SetVideoInfo(probeData *ffprobe.ProbeData, blogVideo *entity.BlogVideo) {
	tag, err := probeData.Format.TagList.GetString("comment")
	if err != nil {
		log.Println(err.Error())
	}
	split := strings.Split(tag, " ")
	blogVideo.Make = split[0]
	blogVideo.Model = split[len(split)-1]
	blogVideo.Height = probeData.Streams[0].Height
	blogVideo.Width = probeData.Streams[0].Width
	blogVideo.CodingStandard = probeData.Streams[0].CodecName
	if probeData.Streams[0].RFrameRate != "" {
		frameRateArray := strings.Split(probeData.Streams[0].RFrameRate, "/")
		dividend, err := strconv.ParseFloat(frameRateArray[0], 64)
		divisor, err := strconv.ParseFloat(frameRateArray[1], 64)
		if err != nil {
			logger.Warn.Println("get frame rate err,", err)
		} else {
			blogVideo.FrameRate = dividend / divisor
		}

	}
}
func (r Universal) SetVideoInfo(probeData *ffprobe.ProbeData, blogVideo *entity.BlogVideo) {
	if probeData.Streams[0].RFrameRate != "" {
		frameRateArray := strings.Split(probeData.Streams[0].RFrameRate, "/")
		dividend, err := strconv.ParseFloat(frameRateArray[0], 64)
		divisor, err := strconv.ParseFloat(frameRateArray[1], 64)
		if err != nil {
			logger.Warn.Println("get frame rate err,", err)
		} else {
			blogVideo.FrameRate = dividend / divisor
		}

	}
	logger.Info.Println(probeData)
}
