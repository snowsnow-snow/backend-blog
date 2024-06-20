package services

import (
	constant "backend-blog"
	"backend-blog/common"
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/result"
	"backend-blog/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	gorm "gorm.io/gorm"
	"time"
)

type userService struct{}

var UserService = &userService{}

// UserRegisterInfo 用户注册信息
type UserRegisterInfo struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}
type UserLoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func (r userService) Create(c *fiber.Ctx) (*models.User, bool, error) {
	userRegisterInfo, parserErr := r.ParamToUser(c)
	if parserErr != nil {
		logger.Error.Println("register bodyParser, msg: ", parserErr)
		return nil, false, result.Err
	}
	err := checkRegisterParam(userRegisterInfo)
	if err != nil {
		return nil, false, err
	}
	count, err := models.UserExistsByUsername(c, userRegisterInfo.Username)
	if err != nil {
		return nil, false, err
	}
	if count >= 1 {
		return nil, true, errors.New("user already exists")
	}
	newPassword, salt := util.GenerateNewPassword(userRegisterInfo.Password)
	createUser := models.User{
		Username: userRegisterInfo.Username,
		Password: newPassword,
		Salt:     salt,
	}
	common.CreateInit(c, &createUser.BaseInfo)
	err = models.InsertUser(createUser, c)
	if err != nil {
		return nil, false, err
	}
	return &createUser, false, nil
}

func (r userService) UploadAvatar(c *fiber.Ctx) error {
	//username := common.GetCurrUsername(c)
	//form, err := c.MultipartForm()
	//if err != nil {
	//	return err
	//}
	//avatar := form.File["avatar"]
	//avatarPath := config.GlobalConfig.File.Path.Public +
	//	common.Separator +
	//	"avatar" +
	//	common.Separator +
	//	avatar[0].Filename
	//c.SaveFile(avatar, avatarPath)
	return result.Success(c)
}

func (r userService) GetUserByUsername(username string) (*[]models.User, *gorm.DB, error) {
	currLoginUsers := new([]models.User)
	users := util.DB.Table(constant.Table.User).Where(models.User{Username: username}).Find(&currLoginUsers)
	if users.Error != nil {
		logger.Error.Println("by username get user:", users.Error)
		return currLoginUsers, nil, users.Error
	}
	return currLoginUsers, users, nil
}

func (r userService) VerifyPassword(password string, currUser models.User) error {
	if util.Encryption(password, currUser.Salt) != currUser.Password {
		return errors.New("wrong password")
	}
	return nil
}

func (r userService) ParamToUser(c *fiber.Ctx) (UserRegisterInfo, error) {
	var userRegisterInfo UserRegisterInfo
	parserErr := c.BodyParser(&userRegisterInfo)
	if parserErr != nil {
		return userRegisterInfo, parserErr
	}
	return userRegisterInfo, nil
}

func (r userService) CreateToken(username string) (string, error) {
	// 声明 Claims
	claims := jwt.MapClaims{
		"username": username,
		"s":        "ss",
		"admin":    true,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString(util.PrivateKey)
	return t, err
}

func (r userService) CheckResetPassword(c *fiber.Ctx) (*UserRegisterInfo, error) {
	var userRegisterInfo UserRegisterInfo
	err := c.BodyParser(&userRegisterInfo)
	if err != nil {
		return nil, err
	}
	if err = checkResetPasswordParam(userRegisterInfo); err != nil {
		return nil, err
	}
	return &userRegisterInfo, nil
}

func (r userService) UpdateUserPassword(newPassword string, currUser models.User, c *fiber.Ctx) error {
	password, salt := util.GenerateNewPassword(newPassword)
	newUser := models.User{Password: password, Salt: salt, BaseInfo: models.BaseInfo{UpdatedTime: models.DateTime{}.Now()}}
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if update := transactionDB.Model(&currUser).Updates(newUser); update.Error != nil {
		return update.Error
	}
	transactionDB.Commit()
	return nil
}

func checkRegisterParam(userRegisterInfo UserRegisterInfo) error {
	if userRegisterInfo.Username == "" {
		return errors.New("please enter a username")
	}
	if userRegisterInfo.Password == "" {
		return errors.New("please enter a password")
	}
	//usernameLegal, err := util.CheckUsername(userRegisterInfo.Username)
	//if err != nil {
	//	logger.Error.Printf("Check username, username: %+v", userRegisterInfo.Username)
	//	return result.Err
	//}
	//if !usernameLegal {
	//	return result.BuildFailResultWithMsg("用户名只包含字母、数字、下划线和连字符，且长度在3到20个字符之间")
	//}
	passwordLegal, err := util.CheckPassword(userRegisterInfo.Password)
	if err != nil {
		return err
	}
	if !passwordLegal {
		return errors.New("密码至少包含一个小写字母、一个大写字母、一个数字和一个特殊字符，并且密码长度至少为8个字符")
	}
	return nil
}

func checkResetPasswordParam(userRegisterInfo UserRegisterInfo) error {
	if userRegisterInfo.OldPassword == "" {
		return errors.New("请输入原登录密码")
	}
	if userRegisterInfo.Password == "" {
		return errors.New("请输入修改后的登录密码")
	}
	passwordLegal, err := util.CheckPassword(userRegisterInfo.Password)
	if err != nil {
		//logger.Error.Println("Check password, password:", userRegisterInfo.Password)
		return err
	}
	if !passwordLegal {
		return errors.New("修改后的密码至少包含一个小写字母、一个大写字母、一个数字和一个特殊字符，并且密码长度至少为8个字符")
	}
	return nil
}
