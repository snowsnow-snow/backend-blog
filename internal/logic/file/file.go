package file

import (
	"backend-blog/config"
	"backend-blog/internal/common"
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/bo"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"backend-blog/internal/service"
	"backend-blog/result"
	"backend-blog/utility"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"sync"
)

type (
	sFile struct{}
)

var (
	implFile = sFile{}
)

func init() {
	service.RegisterFile(New())
}

func New() *sFile {
	return &sFile{}
}

func File() *sFile {
	return &implFile
}

func (s sFile) SaveBatch(c *fiber.Ctx, resourceDesc []*entity.File) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if err := dao.FileDao.BatchInsert(db, resourceDesc); err != nil {
		return err
	}
	return nil
}
func (s sFile) InitRdsInfo(saveFile *bo.ResourceDescBo, c *fiber.Ctx) error {
	file := new(entity.File)
	err := c.BodyParser(file)
	if err != nil {
		return err
	}
	file.ResourceType = saveFile.File.Type
	file = &saveFile.File
	common.CreateInit(c, &file.BaseInfo)
	return nil
}

func (s sFile) Remove(id string, c *fiber.Ctx) (err error) {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	// 查询文件资源路径
	resourceDesc, err := dao.FileDao.SelectPathById(db, id)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	if err = dao.FileDao.DeleteById(db, id); err != nil {
		logger.Error.Printf("delete resource image err %s\n", err)
		return err
	}
	// 删除文件
	if err := utility.DeleteFile(resourceDesc.FilePath); err != nil {
		logger.Error.Printf("delete file err %s path %s\n", err, resourceDesc.FilePath)
		return result.DeleteFileErr
	}
	return nil
}

func (s sFile) RemoveByContentId(c *fiber.Ctx, contentId string) error {
	if contentId == "" {
		return result.WrongParameter
	}
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	// 通过内容 ID 获取所有文件路径信息，等待删除使用
	files, err := dao.FileDao.SelectFilesByContentId(db, contentId)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	videoIds, imagIds, markdownIds := byTypeGetIdGroups(files)
	var wg sync.WaitGroup
	if len(videoIds) > 0 {
		wg.Add(1)
		err = utility.RunWithRecover(func() error {
			defer wg.Done()
			return service.Video().RemoveVideoByFileIds(c, videoIds...)
		})
		if err != nil {
			return err
		}
	}
	if len(imagIds) > 0 {
		wg.Add(1)
		err = utility.RunWithRecover(func() error {
			defer wg.Done()
			return service.Image().RemoveImageByFileIds(c, imagIds...)
		})
		if err != nil {
			return err
		}
	}
	if len(markdownIds) > 0 {
		wg.Add(1)
		err = utility.RunWithRecover(func() error {
			defer wg.Done()
			return service.Markdown().RemoveMarkdownByFileIds(c, markdownIds...)
		})
		if err != nil {
			return err
		}
	}
	wg.Wait()
	paths := make([]string, len(files))
	for index, item := range files {
		paths[index] = config.GlobalConfig.File.Path.System + item.FilePath
	}
	// 删除文件表中的数据
	err = dao.FileDao.DeleteByContentId(db, contentId)
	if err != nil {
		return err
	}
	// 删除此条 blog 中所有的文件
	err = utility.DeleteFolders(paths...)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}
func byTypeGetIdGroups(files []*entity.File) ([]string, []string, []string) {
	var (
		video    []string
		image    []string
		markdown []string
	)
	for _, item := range files {
		if utility.IsInArrayNoCaseSensitive(utility.VideoType, item.Type) {
			video = append(video, item.ID)
			continue
		}
		if utility.IsInArrayNoCaseSensitive(utility.ImgType, item.Type) {
			image = append(image, item.ID)
			continue
		}
		if utility.IsInArrayNoCaseSensitive(utility.MarkDownType, item.Type) {
			markdown = append(markdown, item.ID)
		}
	}
	return video, image, markdown
}

func (s sFile) Update(c *fiber.Ctx) error {
	//db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	//rd := new(entity.File)
	//if err := c.BodyParser(rd); err != nil {
	//	return result.WrongParameter
	//}
	//// 查询文件资源路径
	//rawResourceDesc, err := dao.FileDao.SelectPathById(db, rd.ID)
	//// 删除文件
	//filePathName := rawResourceDesc.FilePath + utility.Separator + rawResourceDesc.FileName + utility.Delimiter + rawResourceDesc.Type
	//if err := utility.DeleteFile(filePathName); err != nil {
	//	logger.Error.Printf("delete file err %s path %s\n", err, filePathName)
	//	return result.DeleteFileErr
	//}
	//common.UpdateInit(c, &rd.BaseInfo)
	//videoInfos, imgInfos, err := FileService.SaveFile(c)
	//if err != nil {
	//	logger.Error.Println("save file", err)
	//	return err
	//}
	//if err = insertResourceDescFile(videoInfos, imgInfos, *rd, c); err != nil {
	//	return err
	//}
	//err = entity.UpdateResourceDesc(*rd, c)
	//if err != nil {
	//	return createContentErr
	//}
	return nil
}

func (s sFile) PublicList(c *fiber.Ctx) (vo.FileVo, error) {
	contentId := c.Query("contentId")
	var fileVo vo.FileVo
	var wg sync.WaitGroup
	wg.Add(2)
	err := utility.RunWithRecover(func() error {
		defer wg.Done()
		imageVos, err := service.Image().GetImagesByContentId(contentId)
		if err != nil {
			return err
		}
		formatFilmParams(imageVos)
		fileVo.ImageVos = imageVos
		return nil
	})
	if err != nil {
		return fileVo, err
	}
	err = utility.RunWithRecover(func() error {
		defer wg.Done()
		videoVos, err := service.Video().GetVideoListByContentId(contentId)
		if err != nil {
			return err
		}
		fileVo.VideoVos = videoVos
		return nil
	})
	if err != nil {
		return fileVo, err
	}
	wg.Wait()
	return fileVo, nil
}

func (s sFile) ManageList(c *fiber.Ctx) ([]entity.File, error) {
	contentId := c.Query("contentId")
	list, err := dao.FileDao.SelectManageList(dao.DB, contentId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s sFile) PublicMarkdownList(c *fiber.Ctx) (*vo.FileVo, error) {
	//service.Markdown().Remove()
	//TODO implement me
	panic("implement me")
}

func (s sFile) GetFilePath(c *fiber.Ctx, id string, compressionRatio, extension string) (string, error) {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	file, err := dao.FileDao.SelectPathById(db, id)
	if err != nil {
		return "", err
	}
	var filePath string
	if extension == "" {
		extension = file.Extension
	} else if utility.IsInArrayNoCaseSensitive(constant.ImageExtensions, extension) &&
		utility.IsInArrayNoCaseSensitive(constant.VideoExtensions, extension) {
		extension = file.Extension
	}
	if compressionRatio == "100" || compressionRatio == "" || compressionRatio == "html" {
		filePath = file.FilePath + file.FileName + utility.Point + extension
	} else {
		filePath = file.FilePath + file.FileName + utility.Delimiter + compressionRatio + utility.Point + extension
	}
	if file.Type == "mov" {
		c.Set("Content-Type", "video/quicktime")
	}
	return filePath, nil
}
func formatFilmParams(list []*vo.ImageVo) {
	for _, item := range list {
		setFilmMode(item)
	}
}

func setFilmMode(image *vo.ImageVo) {
	//img.FilmModeFormat = utility.GetChineseFilmMode(img.FilmMode)
	image.FilmModeFormat = image.FilmMode
	image.DynamicRangeFormat = utility.GetChineseDynamicRange(image.DynamicRange)
	image.WhiteBalanceFormat = utility.GetChineseWhiteBalance(image.DynamicRange)
	image.WhiteBalanceFineTuneFormat = utility.GetWhiteBalanceFineTuneFormat(image.WhiteBalanceFineTune)
	image.SharpnessFormat = utility.GetChineseGenericDescriptionMap(image.Sharpness)
	image.GrainEffectRoughnessFormat = utility.GetChineseGenericDescriptionMap(image.GrainEffectRoughness)
	image.ColorChromeEffectFormat = utility.GetChineseGenericDescriptionMap(image.ColorChromeEffect)
	image.ShadowToneFormat = utility.GetNumeric(image.ShadowTone)
	image.HighlightToneFormat = utility.GetNumeric(image.HighlightTone)
	image.SaturationFormat = utility.GetNumericAndCharParam(image.Saturation)
	image.NoiseReductionFormat = utility.GetNumeric(image.NoiseReduction)
	image.ColorChromeFXBlueFormat = utility.GetChineseGenericDescriptionMap(image.ColorChromeFXBlue)

}
