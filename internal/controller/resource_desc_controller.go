package controller

//type fileController struct {
//}
//
//var FileController = new(fileController)

//func (r fileController) ViewImage(imageId string, compressionRatio string) (string, error) {
//	path, err := service.Image().View(imageId, compressionRatio)
//	if err != nil {
//		return "", err
//	}
//	return path, nil
//}
//
//
//func (r fileController) PublicMarkdownList(c *fiber.Ctx) error {
//	list, err := service.ResourceDescService.PublicMarkdownList(c)
//	if err != nil {
//		return result.ErrorWithMsg(c, err.Error())
//	}
//	return result.SuccessData(c, list)
//}
//
//func (r fileController) PublicList(c *fiber.Ctx) error {
//	list, err := service.ResourceDescService.PublicList(c)
//	if err != nil {
//		return result.ErrorWithMsg(c, err.Error())
//	}
//	return result.SuccessData(c, list)
//}
//
//func (r fileController) AddVideo(c *fiber.Ctx) error {
//	err := service.ResourceDescService.AddVideo(c)
//	if err != nil {
//		logger.Error.Println("add resource desc error", err)
//		return result.ErrorWithMsg(c, result.WrongParameter.Error())
//	}
//	return result.Success(c)
//}
//
//func (r fileController) AddLivePhotos(c *fiber.Ctx) error {
//	err := service.ResourceDescService.AddLivePhotos(c)
//	if err != nil {
//		logger.Error.Println("add resource desc error", err)
//		return result.ErrorWithMsg(c, result.WrongParameter.Error())
//	}
//	return result.Success(c)
//}
//
//func (r fileController) AddMarkDown(c *fiber.Ctx) error {
//	err := service.ResourceDescService.AddMarkDown(c)
//	if err != nil {
//		logger.Error.Println("add resource desc error", err)
//		return result.ErrorWithMsg(c, result.WrongParameter.Error())
//	}
//	return result.Success(c)
//}
