package services

import (
	constant "backend-blog"
	"backend-blog/common"
	"backend-blog/dto"
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

func (r resourceDescService) Add(c *fiber.Ctx) error {
	rd := new(models.ResourceDesc)
	err := c.BodyParser(rd)
	if err != nil {
		return result.WrongParameter
	}
	common.CreateInit(c, &rd.BaseInfo)
	videoInfos, imgInfos, err := FileService.SaveFile(c)
	if err != nil {
		logger.Error.Println("Save file", err)
		return err
	}
	if len(videoInfos) > 0 {
		video := videoInfos[0]
		rd.ResourceType = video.Type
		if err = VideoService.Add(c, &video.VideoInfo); err != nil {
			logger.Error.Println("Insert video err", err)
			return result.SaveVideoErr
		}
		rd.FileId = video.VideoInfo.ID
	}
	fileId := ""
	if len(imgInfos) > 0 {
		img := imgInfos[0]
		rd.ResourceType = img.Type
		if err = ImgService.Add(c, &img.ImgInfo); err != nil {
			logger.Error.Println("Insert img err", err)
			return result.SaveImgErr
		}
		rd.FileId = img.ImgInfo.ID
		fileId = img.ImgInfo.ID
		rd.FileInfo = img.FileInfo
	}
	err = models.CreateResourceDesc(*rd, c)
	if err != nil {
		return createContentErr
	}
	if rd.Cover == 1 {
		err = models.UpdateContentCover(fileId, rd.ContentId, c)
		if err != nil {
			return err
		}
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
func getImgIdByResourceDescList(list []models.ResourceDesc) []interface{} {
	var arr []interface{}
	for _, desc := range list {
		arr = append(arr, desc.ID)
	}
	return arr
}
func getIds(ids string) []interface{} {
	var arr []interface{}
	split := strings.Split(ids, ",")
	for _, id := range split {
		arr = append(arr, id)
	}
	//var s []interface{} = strings.Split(ids, ",")
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
	common.UpdateInit(&rd.BaseInfo)
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

func (r resourceDescService) List(c *fiber.Ctx) ([]*dto.ResourcesDescImg, error) {
	db := util.DB
	contentId := c.Query("contentId")
	//if contentId := c.Query("contentId"); contentId != "" {
	//	db = db.Where("RD.content_id = ?", contentId)
	//}
	listResourceDesc, err := dto.SelectResourcesDescImgList(db, contentId)
	//listResourceDesc, err := models.ListResourceDesc(db)
	if err != nil {
		logger.Error.Println("select resource desc list err", err)
		return nil, selectContentErr
	}
	return listResourceDesc, nil
}

func (r resourceDescService) PublicList(c *fiber.Ctx) ([]*dto.PublicResourcesDescImg, error) {
	db := util.DB
	contentId := c.Query("contentId")
	//if contentId := c.Query("contentId"); contentId != "" {
	//	db = db.Where("RD.content_id = ?", contentId)
	//}
	listResourceDesc, err := dto.SelectPublicResourcesDescImgList(db, contentId)
	//listResourceDesc, err := models.ListResourceDesc(db)
	if err != nil {
		logger.Error.Println("select resource desc list err", err)
		return nil, selectContentErr
	}
	return listResourceDesc, nil
}
