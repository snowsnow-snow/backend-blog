package user

import (
	"backend-blog/internal/common"
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/logger"
	"backend-blog/internal/model"
	"backend-blog/internal/model/do"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"backend-blog/result"
	"backend-blog/utility"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
)

type (
	sUser struct{}
)

var (
	implUser = sUser{}
)

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}

func User() *sUser {
	return &implUser
}

func (s sUser) Add(c *fiber.Ctx) (*entity.User, error) {
	userRegisterInfo, parserErr := User().ParamConvertUser(c)
	if parserErr != nil {
		logger.Error.Println("register bodyParser, msg: ", parserErr)
		return nil, result.Err
	}
	// 检查注册信息是否正确
	err := checkRegisterParam(userRegisterInfo)
	if err != nil {
		return nil, err
	}
	// 系统中是否存在相同用户名
	count, err := dao.User.NumberOfUsername(c, userRegisterInfo.Username)
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		return nil, errors.New("user already exists")
	}
	// 通过盐与用户输入的用户名生成数据库储存的密码
	newPassword, salt := utility.GenerateNewPassword(userRegisterInfo.Password)
	createUser := entity.User{
		Username: userRegisterInfo.Username,
		Password: newPassword,
		Salt:     salt,
	}
	common.CreateInit(c, &createUser.BaseInfo)
	err = dao.User.Insert(createUser, c)
	if err != nil {
		return nil, err
	}
	return &createUser, nil
}
func (s sUser) UploadAvatar(c *fiber.Ctx) error {
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

func (s sUser) UpdateUserPassword(newPassword string, currUser entity.User, c *fiber.Ctx) error {
	password, salt := utility.GenerateNewPassword(newPassword)
	newUser := entity.User{Password: password, Salt: salt, BaseInfo: entity.BaseInfo{UpdatedTime: model.Now()}}
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if update := transactionDB.Model(&currUser).Updates(newUser); update.Error != nil {
		return update.Error
	}
	return nil
}

func (s sUser) GetUserByUsername(username string) (*[]entity.User, *gorm.DB, error) {
	currLoginUsers := new([]entity.User)
	users := dao.DB.Table(constant.Table.User).Where(entity.User{Username: username}).Find(&currLoginUsers)
	if users.Error != nil {
		logger.Error.Println("by username get user:", users.Error)
		return currLoginUsers, nil, users.Error
	}
	return currLoginUsers, users, nil
}

func (s sUser) VerifyPassword(password string, currUser entity.User) error {
	if utility.Encryption(password, currUser.Salt) != currUser.Password {
		return errors.New("wrong password")
	}
	return nil
}

func (s sUser) CreateToken(username string) (string, error) {
	// 声明 Claims
	claims := jwt.MapClaims{
		"username": username,
		"s":        "ss",
		"admin":    true,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(utility.PrivateKey)
	if err != nil {
		return "", err
	}
	return t, err
}

func (s sUser) CheckResetPassword(c *fiber.Ctx) (*do.UserRegister, error) {
	var userRegisterInfo do.UserRegister
	err := c.BodyParser(&userRegisterInfo)
	if err != nil {
		return nil, err
	}
	if err = checkResetPasswordParam(userRegisterInfo); err != nil {
		return nil, err
	}
	return &userRegisterInfo, nil
}

func (s sUser) Delete(c *fiber.Ctx, id string) (err error) {
	//TODO implement me
	panic("implement me")
}

// ParamConvertUser 参数转换用户
func (s sUser) ParamConvertUser(c *fiber.Ctx) (do.UserRegister, error) {
	var userRegisterInfo do.UserRegister
	parserErr := c.BodyParser(&userRegisterInfo)
	if parserErr != nil {
		return userRegisterInfo, parserErr
	}
	return userRegisterInfo, nil
}
func checkRegisterParam(userRegisterInfo do.UserRegister) error {
	if userRegisterInfo.Username == "" {
		return errors.New("please enter a username")
	}
	if userRegisterInfo.Password == "" {
		return errors.New("please enter a password")
	}
	//usernameLegal, err := utility.CheckUsername(userRegisterInfo.Username)
	//if err != nil {
	//	logger.Error.Printf("Check username, username: %+v", userRegisterInfo.Username)
	//	return result.Err
	//}
	//if !usernameLegal {
	//	return result.BuildFailResultWithMsg("用户名只包含字母、数字、下划线和连字符，且长度在3到20个字符之间")
	//}
	passwordLegal, err := utility.CheckPassword(userRegisterInfo.Password)
	if err != nil {
		return err
	}
	if !passwordLegal {
		return errors.New("密码至少包含一个小写字母、一个大写字母、一个数字和一个特殊字符，并且密码长度至少为8个字符")
	}
	return nil
}

func checkResetPasswordParam(userRegister do.UserRegister) error {
	if userRegister.OldPassword == "" {
		return errors.New("请输入原登录密码")
	}
	if userRegister.Password == "" {
		return errors.New("请输入修改后的登录密码")
	}
	passwordLegal, err := utility.CheckPassword(userRegister.Password)
	if err != nil {
		//logger.Error.Println("Check password, password:", userRegisterInfo.Password)
		return err
	}
	if !passwordLegal {
		return errors.New("修改后的密码至少包含一个小写字母、一个大写字母、一个数字和一个特殊字符，并且密码长度至少为8个字符")
	}
	return nil
}
