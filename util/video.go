package util

import (
	"backend-blog/logger"
	"backend-blog/models"
	"context"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var makes = [...]string{"Apple", "iPhone", "FUJIFILM", "qt"}

func ReadInfo(path string) (*models.VideoInfo, error) {
	ffprobe.SetFFProbeBinPath("/Users/snowsnowsnow/Library/CloudStorage/OneDrive-个人/附件/ffprobe/ffprobe")
	// 打开要解析的音视频文件
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	fileReader, err := os.Open(path)
	//fileReader, err := os.Open("/Users/yangjian/Library/CloudStorage/OneDrive-个人/图片/本机照片/2023/03/AE129.MOV")
	//fileReader, err := os.Open("/Users/yangjian/Library/CloudStorage/OneDrive-个人/FUJIFILM/相册/2023-04-28@南京/AE362.MOV")
	if err != nil {
		return nil, err
	}

	data, err := ffprobe.ProbeReader(ctx, fileReader)
	if err != nil {
		return nil, err
	}
	videoInfo := models.VideoInfo{}
	setBrand(data, &videoInfo)
	return &videoInfo, nil
}
func setFrameRate() {

}
func setBrand(probeData *ffprobe.ProbeData, videoInfo *models.VideoInfo) {
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
					handle.SetVideoInfo(probeData, videoInfo)
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
	SetVideoInfo(probeData *ffprobe.ProbeData, videoInfo *models.VideoInfo)
}

type Apple struct{}
type Fujifilm struct{}
type Universal struct{}

func (r Apple) SetVideoInfo(probeData *ffprobe.ProbeData, videoInfo *models.VideoInfo) {
	makeName, err := probeData.Format.TagList.GetString("com.apple.quicktime.make")
	if err != nil {
		log.Println(err)
	} else {
		videoInfo.Make = makeName
	}
	model, err := probeData.Format.TagList.GetString("com.apple.quicktime.model")
	if err != nil {
		log.Println(err)
	} else {
		videoInfo.Model = model
	}
}

func (r Fujifilm) SetVideoInfo(probeData *ffprobe.ProbeData, videoInfo *models.VideoInfo) {
	tag, err := probeData.Format.TagList.GetString("comment")
	if err != nil {
		log.Println(err.Error())
	}
	split := strings.Split(tag, " ")
	videoInfo.Make = split[0]
	videoInfo.Model = split[len(split)-1]
	videoInfo.Height = probeData.Streams[0].Height
	videoInfo.Width = probeData.Streams[0].Width
	videoInfo.CodingStandard = probeData.Streams[0].CodecName
	if probeData.Streams[0].RFrameRate != "" {
		frameRateArray := strings.Split(probeData.Streams[0].RFrameRate, "/")
		dividend, err := strconv.ParseFloat(frameRateArray[0], 64)
		divisor, err := strconv.ParseFloat(frameRateArray[1], 64)
		if err != nil {
			logger.Warn.Println("get frame rate err,", err)
		} else {
			videoInfo.FrameRate = dividend / divisor
		}

	}
}
func (r Universal) SetVideoInfo(probeData *ffprobe.ProbeData, videoInfo *models.VideoInfo) {
	if probeData.Streams[0].RFrameRate != "" {
		frameRateArray := strings.Split(probeData.Streams[0].RFrameRate, "/")
		dividend, err := strconv.ParseFloat(frameRateArray[0], 64)
		divisor, err := strconv.ParseFloat(frameRateArray[1], 64)
		if err != nil {
			logger.Warn.Println("get frame rate err,", err)
		} else {
			videoInfo.FrameRate = dividend / divisor
		}

	}
	logger.Info.Println(probeData)
}
