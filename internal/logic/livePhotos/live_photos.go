package livePhotos

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/handle"
	"backend-blog/internal/model/bo"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"backend-blog/utility"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type (
	sLivePhotos struct{}
)

var (
	implLivePhotos = sLivePhotos{}
)

func init() {
	service.RegisterLivePhotos(New())
}

func New() *sLivePhotos {
	return &sLivePhotos{}
}

func LivePhotos() *sLivePhotos {
	return &implLivePhotos
}
func (s sLivePhotos) Add(c *fiber.Ctx, img *entity.BlogImage) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s sLivePhotos) SaveLivePhotosAndResourceDesc(c *fiber.Ctx) (err error) {
	fileBos, err := utility.ParamToFileBo(c)
	if err != nil {
		return err
	}
	// 保存文件
	resourceDescDos, err := utility.SaveFiles(c, fileBos)
	if err != nil {
		return err
	}
	rdsInfos := make([]*entity.File, len(resourceDescDos))
	videoInfos, imgInfos := make([]*bo.ResourceDescBo, len(resourceDescDos)/2), make([]*bo.ResourceDescBo, len(resourceDescDos)/2)
	videoIndex, imgIndex := 0, 0
	rdsInfoInitHandle := handle.LivePhotos{}
	for index, resourceDescDo := range resourceDescDos {
		rdsInfoInitHandle.InitFileInfo(resourceDescDo, c)
		// 根据文件信息生成文件通用数据
		err = service.File().InitRdsInfo(resourceDescDo, c)
		if err != nil {
			return err
		}
		rdsInfos[index] = &resourceDescDo.File
		if utility.IsInArrayNoCaseSensitive(utility.ImgType, resourceDescDo.File.Extension) {
			imgInfos[imgIndex] = resourceDescDo
			imgInfos[imgIndex].BlogImage.FileId = resourceDescDo.File.ID
			if rdsInfos[index].Cover == constant.YesOrNo.Yes {
				err := service.Content().SetTheCoverContent(c, &entity.BlogContent{
					BaseInfo: entity.BaseInfo{
						ID: rdsInfos[index].ContentId,
					},
					TheCover: rdsInfos[index].ID,
				})
				if err != nil {
					return err
				}
			}
			imgIndex++
			continue
		}
		if utility.IsInArrayNoCaseSensitive(utility.VideoType, resourceDescDo.File.Extension) {
			videoInfos[videoIndex] = resourceDescDo
			videoInfos[videoIndex].BlogVideo.FileId = resourceDescDo.File.ID
			videoIndex++
		}
	}
	// 照片与视频关联
	for _, video := range videoInfos {
		for _, img := range imgInfos {
			if strings.Split(img.File.RawFileName, ".")[0] ==
				strings.Split(video.File.RawFileName, ".")[0] {
				img.BlogImage.LivePhotosId = video.BlogVideo.FileId
				video.BlogVideo.LivePhoto = 1
			}
		}
	}
	for _, img := range imgInfos {
		err := service.Image().Add(c, &img.BlogImage)
		if err != nil {
			return err
		}
	}
	for _, video := range videoInfos {
		err := service.Video().Add(c, &video.BlogVideo)
		if err != nil {
			return err
		}
	}
	err = service.File().SaveBatch(c, rdsInfos)
	if err != nil {
		return err
	}
	return nil
}

func (s sLivePhotos) Delete(c *fiber.Ctx, id string) (err error) {
	//TODO implement me
	panic("implement me")
}
