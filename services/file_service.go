package services

import (
	"backend-blog/dto"
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"runtime/debug"
	"strings"
	"sync"
)

type fileService struct {
}

var (
	FileService      = &fileService{}
	saveFileError    = errors.New("save file error")
	notFoundResource = errors.New("not found resource")
	getImgError      = errors.New("get img error")
)

func (r fileService) SaveFile(c *fiber.Ctx) ([]*models.SaveFileInfo, []*models.SaveFileInfo, error) {
	var videos []*models.FileInfo
	var images []*models.FileInfo
	saveFiles, err := util.SaveFiles(c)
	if err != nil {
		logger.Error.Println(err)
		return nil, nil, saveFileError
	}
	var wg sync.WaitGroup
	for i := range saveFiles {
		currFileInfo := saveFiles[i]
		if util.IsInArray(util.VideoType, strings.ToUpper(currFileInfo.Type)) {
			if len(videos) == 0 {
				wg.Add(1)
			}
			videos = append(videos, saveFiles[i])
			continue
		}
		if util.IsInArray(util.ImgType, strings.ToUpper(currFileInfo.Type)) {
			if len(images) == 0 {
				wg.Add(1)
			}
			images = append(images, saveFiles[i])
		}
	}
	var videoInfos []*models.SaveFileInfo
	var imageInfos []*models.SaveFileInfo
	if len(videos) > 0 {
		go func() {
			defer wg.Done()
			for i := range videos {
				videoInfo := videos[i]
				video, err := util.ReadInfo(videoInfo.FilePath + util.Separator + videoInfo.FileName)
				if err != nil {
					logger.Error.Println(err)
					continue
				}
				videoInfos = append(videoInfos, &models.SaveFileInfo{
					FileInfo:  *videoInfo,
					VideoInfo: *video,
				})
			}
		}()
	}
	if len(images) > 0 {
		go func() {
			defer wg.Done()
			for i := range images {
				imgInfo := images[i]
				img := util.ReadExif(imgInfo.FilePath, imgInfo.FileName, imgInfo.Type)
				if img == nil {
					img = &models.ImgInfo{}
				}
				imageInfos = append(imageInfos, &models.SaveFileInfo{
					FileInfo: *imgInfo,
					ImgInfo:  *img,
				})
			}
		}()
	}
	wg.Wait()
	return videoInfos, imageInfos, nil
}

func (r fileService) ByIdGetImgInfo(imgId string, args ...string) (*models.ImgInfo, error) {
	db := util.DB.Select(args)
	imgInfo, err := models.SelectImgInfoById(imgId, db)
	if err != nil {
		logger.Error.Println(err, string(debug.Stack()))
		return nil, getImgError
	}
	if imgInfo == nil {
		return nil, notFoundResource
	}
	return imgInfo, nil
}

func (r fileService) ByIdGetVideoInfo(imgId string, args ...string) (*models.VideoInfo, error) {
	db := util.DB.Select(args)
	videoInfo, err := models.SelectVideoInfoById(imgId, db)
	if err != nil {
		logger.Error.Println(err, string(debug.Stack()))
		return nil, getImgError
	}
	if videoInfo == nil {
		return nil, notFoundResource
	}
	return videoInfo, nil
}

func (r fileService) ByIdCompressionRatioGetImgPath(imgId string, compressionRatio string) (string, error) {
	imgInfo, err := dto.SelectResourcesDescImgById(" WHERE RD.file_id = ?", imgId)
	if err != nil {
		return "", err
	}
	var fileName string
	if compressionRatio == "100" {
		fileName = imgInfo.FilePath + util.Separator + imgInfo.FileName + util.Point + imgInfo.Type
	} else {
		fileName = imgInfo.FilePath + util.Separator + imgInfo.FileName + util.Delimiter + compressionRatio + util.Point + imgInfo.Type
	}
	return fileName, nil
}

func (r fileService) ByIdGetVideoPath(videoId string) (string, error) {
	videoInfo, err := r.ByIdGetVideoInfo(videoId, "file_path", "file_name")
	if err != nil {
		return "", err
	}
	println(videoInfo)
	//return videoInfo.FilePath + common.Separator + videoInfo.FileName, nil
	return "", nil
}
