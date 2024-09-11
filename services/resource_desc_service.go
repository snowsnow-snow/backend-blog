package services

import (
	constant "backend-blog"
	"backend-blog/common"
	"backend-blog/dto"
	"backend-blog/handle"
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/result"
	"backend-blog/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type resourceDescService struct {
}

var (
	ResourceDescService         = &resourceDescService{}
	deleteResourceDescErr       = errors.New("delete resource desc err")
	selectDeleteResourceDescErr = errors.New("select resource desc err")
)

func (r resourceDescService) AddImg(c *fiber.Ctx) error {
	saveFiles, err := util.SaveFiles(c)
	if err != nil {
		return err
	}
	rdsInfos := make([]*models.ResourceDesc, len(saveFiles))
	rdsInfoInitHandle := handle.Img{}
	for index, file := range saveFiles {
		info, err := rdsInfoInitHandle.InitFileInfo(file, c)
		if err != nil {
			return err
		}
		if err = ImgService.Add(c, &info.ImgInfo); err != nil {
			return err
		}
		rdsInfos[index], err = initRdsInfo(info, c)
		if rdsInfos[index].Cover == 1 {
			err := setContentCover(*rdsInfos[index], c)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	err = models.CreateInBatchesResourceDesc(rdsInfos, c)
	if err != nil {
		return err
	}
	return nil
}

func (r resourceDescService) AddVideo(c *fiber.Ctx) error {
	saveFiles, err := util.SaveFiles(c)
	if err != nil {
		return err
	}
	rdsInfos := make([]*models.ResourceDesc, len(saveFiles))
	rdsInfoInitHandle := handle.Video{}
	for index, file := range saveFiles {
		info, err := rdsInfoInitHandle.InitFileInfo(file, c)
		if err != nil {
			return err
		}
		if err = ImgService.Add(c, &info.ImgInfo); err != nil {
			return err
		}
		rdsInfos[index], err = initRdsInfo(info, c)
		if err != nil {
			return err
		}
	}
	err = models.CreateInBatchesResourceDesc(rdsInfos, c)
	if err != nil {
		return err
	}
	return nil
}

func (r resourceDescService) AddLivePhotos(c *fiber.Ctx) error {
	saveFiles, err := util.SaveFiles(c)
	if err != nil {
		return err
	}
	rdsInfos := make([]*models.ResourceDesc, len(saveFiles))
	videoInfos, imgInfos := make([]*models.SaveFileInfo, len(saveFiles)/2), make([]*models.SaveFileInfo, len(saveFiles)/2)
	videoIndex, imgIndex := 0, 0
	rdsInfoInitHandle := handle.LivePhotos{}
	for index, file := range saveFiles {
		info, err := rdsInfoInitHandle.InitFileInfo(file, c)
		if err != nil {
			return err
		}
		rdsInfos[index], err = initRdsInfo(info, c)
		if err != nil {
			return err
		}
		if util.IsInArrayNoCaseSensitive(util.ImgType, file.Type) {
			imgInfos[imgIndex] = info
			if rdsInfos[index].Cover == 1 {
				err := setContentCover(*rdsInfos[index], c)
				if err != nil {
					return err
				}
			}
			imgIndex++
			continue
		}
		if util.IsInArrayNoCaseSensitive(util.VideoType, file.Type) {
			videoInfos[videoIndex] = info
			videoIndex++
		}
	}
	for _, video := range videoInfos {
		for _, img := range imgInfos {
			if strings.Split(img.FileInfo.RawFileName, ".")[0] == strings.Split(video.FileInfo.RawFileName, ".")[0] {
				img.LivePhotosId = video.VideoInfo.ID
				video.VideoInfo.LivePhoto = 1
			}
		}
	}
	for _, img := range imgInfos {
		err := ImgService.Add(c, &img.ImgInfo)
		if err != nil {
			return err
		}
	}
	for _, video := range videoInfos {
		err := VideoService.Add(c, &video.VideoInfo)
		if err != nil {
			return err
		}
	}
	err = models.CreateInBatchesResourceDesc(rdsInfos, c)
	if err != nil {
		return err
	}
	return nil
}

func (r resourceDescService) AddMarkDown(c *fiber.Ctx) error {
	saveFiles, err := util.SaveFiles(c)
	if err != nil {
		return err
	}
	rdsInfos := make([]*models.ResourceDesc, len(saveFiles))
	rdsInfoInitHandle := handle.Markdown{}
	for index, file := range saveFiles {
		info, err := rdsInfoInitHandle.InitFileInfo(file, c)
		if err != nil {
			return err
		}
		if err = MarkdownService.Add(c, &info.MarkdownInfo); err != nil {
			return err
		}
		rdsInfos[index], err = initRdsInfo(info, c)
		if err != nil {
			return err
		}
	}
	err = models.CreateInBatchesResourceDesc(rdsInfos, c)
	if err != nil {
		return err
	}
	return nil
}
func initRdsInfo(saveFile *models.SaveFileInfo, c *fiber.Ctx) (*models.ResourceDesc, error) {
	rds := new(models.ResourceDesc)
	err := c.BodyParser(rds)
	if err != nil {
		return nil, err
	}
	rds.ResourceType = saveFile.Type
	if saveFile.VideoInfo.ID != "" {
		rds.FileId = saveFile.VideoInfo.ID
	} else if saveFile.ImgInfo.ID != "" {
		rds.FileId = saveFile.ImgInfo.ID
	} else if saveFile.MarkdownInfo.ID != "" {
		rds.FileId = saveFile.MarkdownInfo.ID
	}
	rds.FileInfo = saveFile.FileInfo
	common.CreateInit(c, &rds.BaseInfo)
	return rds, nil
}

func setContentCover(rd models.ResourceDesc, c *fiber.Ctx) error {
	err := models.UpdateContentCover(rd.FileId, rd.ContentId, c)
	if err != nil {
		return err
	}
	return nil
}
func (r resourceDescService) Delete(id string, c *fiber.Ctx) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	db = db.Select("resource_type")
	resourceDesc, err := models.SelectResourceDescById(id, db)
	if err != nil {
		logger.Error.Println(err)
		return selectDeleteResourceDescErr
	}
	if resourceDesc == nil {
		return nil
	}
	err = models.DeleteResourceDesc(id, c)
	if err != nil {
		logger.Error.Println("delete resource desc err,", err)
		return deleteResourceDescErr
	}
	if err = deleteResourceDescFile(id, resourceDesc.ResourceType, c); err != nil {
		return err
	}
	return nil
}

func (r resourceDescService) DeleteByResIds(ids string, c *fiber.Ctx) error {
	split := strings.Split(ids, ",")
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	list := new([]models.ResourceDesc)
	db = db.Table(constant.Table.ResourceDesc).
		Select("id", "resource_type", "file_path", "file_name", "file_id", "type").
		Where(util.OrByParams("id", len(split)), getIds(ids)...).
		Scan(&list)
	return r.DeleteByRes(*list, c)
}

func (r resourceDescService) DeleteByRes(list []models.ResourceDesc, c *fiber.Ctx) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	resIds, videoIds, imgIds := getDeleteIds(list)
	tx := db.Where(util.OrByParams("id", len(resIds)), resIds...).
		Delete(&models.ResourceDesc{})
	if len(videoIds) > 0 {
		tx = db.Where(util.OrByParams("id", len(videoIds)), videoIds...).
			Delete(&models.VideoInfo{})
	}
	if len(imgIds) > 0 {
		tx = db.Where(util.OrByParams("id", len(imgIds)), imgIds...).
			Delete(&models.ImgInfo{})
	}
	if tx.Error != nil {
		logger.Error.Println("Delete resources error", tx.Error)
		return result.Err
	}
	deleteFiles(list)
	return nil
}

func deleteFiles(list []models.ResourceDesc) {
	for _, desc := range list {
		for _, ratio := range util.CompressionRatio {
			path := desc.FilePath + util.Separator +
				desc.FileName + util.Delimiter + strconv.Itoa(ratio) + util.Point + desc.FileInfo.Type
			err := util.DeleteFile(path)
			if err != nil {
				logger.Error.Println("Delete file error, path: ", path, err)
			}
		}
		path := desc.FilePath + util.Separator + desc.FileName + util.Point + desc.FileInfo.Type
		err := util.DeleteFile(path)
		if err != nil {
			logger.Error.Println("Delete file error, path: ", path, err)
		}
	}
}
func getDeleteIds(list []models.ResourceDesc) ([]interface{}, []interface{}, []interface{}) {
	var resIds, videoIds, imgIds []interface{}
	for _, desc := range list {
		resIds = append(resIds, desc.ID)
		if util.IsInArray(util.VideoType, strings.ToUpper(desc.Type)) {
			videoIds = append(videoIds, desc.FileId)
		} else {
			imgIds = append(imgIds, desc.FileId)
		}
	}
	return resIds, videoIds, imgIds
}

func getIds(ids string) []interface{} {
	var arr []interface{}
	split := strings.Split(ids, ",")
	for _, id := range split {
		arr = append(arr, id)
	}
	return arr
}

func deleteResourceDescFile(resourceId string, resourceType string, c *fiber.Ctx) error {
	if util.IsInArray(util.VideoType, resourceType) {
		if err := VideoService.Delete(resourceId, c); err != nil {
			return err
		}
	}
	if util.IsInArray(util.ImgType, resourceType) {
		if err := ImgService.Delete(resourceId, c); err != nil {
			return err
		}
	}
	return nil
}
func insertResourceDescFile(videoInfos []*models.SaveFileInfo, imgInfos []*models.SaveFileInfo, rd models.ResourceDesc, c *fiber.Ctx) error {
	if len(videoInfos) > 0 {
		video := videoInfos[0]
		rd.ResourceType = video.Type
		if err := VideoService.Add(c, &video.VideoInfo); err != nil {
			logger.Error.Println("insert video err", err)
			return result.SaveVideoErr
		}
	}
	if len(imgInfos) > 0 {
		img := imgInfos[0]
		rd.ResourceType = img.Type
		if err := ImgService.Add(c, &img.ImgInfo); err != nil {
			logger.Error.Println("insert img err", err)
			return result.SaveImgErr
		}
	}
	return nil
}

func (r resourceDescService) Update(c *fiber.Ctx) error {
	rd := new(models.ResourceDesc)
	err := c.BodyParser(rd)
	if err != nil {
		return result.WrongParameter
	}
	common.UpdateInit(c, &rd.BaseInfo)
	if err = deleteResourceDescFile(rd.ID, rd.ResourceType, c); err != nil {
		return err
	}
	videoInfos, imgInfos, err := FileService.SaveFile(c)
	if err != nil {
		logger.Error.Println("save file", err)
		return err
	}
	if err = insertResourceDescFile(videoInfos, imgInfos, *rd, c); err != nil {
		return err
	}
	err = models.UpdateResourceDesc(*rd, c)
	if err != nil {
		return createContentErr
	}
	return nil
}

func (r resourceDescService) List(c *fiber.Ctx) ([]*models.PublicResourceDesc, error) {
	db := util.DB
	contentId := c.Query("contentId")
	images, err := dto.SelectResourcesDescImgList(db, contentId)
	videos, err := dto.SelectResourcesDescVideoList(db, contentId)
	println(images, videos)
	//listResourceDesc, err := models.ListResourceDesc(db)
	if err != nil {
		logger.Error.Println("select resource desc list err", err)
		return nil, selectContentErr
	}
	return nil, nil
}

func (r resourceDescService) PublicMarkdownList(c *fiber.Ctx) ([]*models.PublicMarkdownResourceDesc, error) {
	db := util.DB
	contentId := c.Query("contentId")
	list, err := dto.SelectResourcesDescMarkdownList(db, contentId)
	if err != nil {
		logger.Error.Println("select resource desc list err", err)
		return nil, selectContentErr
	}
	return list, nil
}

func (r resourceDescService) PublicList(c *fiber.Ctx) ([]*models.PublicResourceDesc, error) {
	db := util.DB
	contentId := c.Query("contentId")
	images, err := dto.SelectResourcesDescImgList(db, contentId)
	if err != nil {
		logger.Error.Println("select img resource desc list err", err)
		return nil, selectContentErr
	}
	if len(images) > 0 {
		formatFilmParams(images)
	}
	videos, err := dto.SelectResourcesDescVideoList(db, contentId)
	if err != nil {
		logger.Error.Println("select video resource desc list err", err)
		return nil, selectContentErr
	}

	return append(images, videos...), nil
}

func formatFilmParams(list []*models.PublicResourceDesc) {
	for _, item := range list {
		setFilmMode(item)
	}
}

func setFilmMode(desc *models.PublicResourceDesc) {
	//img.FilmModeFormat = util.GetChineseFilmMode(img.FilmMode)
	desc.FilmModeFormat = desc.ImgInfo.FilmMode
	desc.DynamicRangeFormat = util.GetChineseDynamicRange(desc.ImgInfo.DynamicRange)
	desc.WhiteBalanceFormat = util.GetChineseWhiteBalance(desc.ImgInfo.DynamicRange)
	desc.WhiteBalanceFineTuneFormat = util.GetWhiteBalanceFineTuneFormat(desc.ImgInfo.WhiteBalanceFineTune)
	desc.SharpnessFormat = util.GetChineseGenericDescriptionMap(desc.ImgInfo.Sharpness)
	desc.GrainEffectRoughnessFormat = util.GetChineseGenericDescriptionMap(desc.ImgInfo.GrainEffectRoughness)
	desc.ColorChromeEffectFormat = util.GetChineseGenericDescriptionMap(desc.ImgInfo.ColorChromeEffect)
	desc.ShadowToneFormat = util.GetNumeric(desc.ImgInfo.ShadowTone)
	desc.HighlightToneFormat = util.GetNumeric(desc.ImgInfo.HighlightTone)
	desc.SaturationFormat = util.GetNumericAndCharParam(desc.ImgInfo.Saturation)
	desc.NoiseReductionFormat = util.GetNumeric(desc.ImgInfo.NoiseReduction)
	desc.ColorChromeFXBlueFormat = util.GetChineseGenericDescriptionMap(desc.ImgInfo.ColorChromeFXBlue)

}
