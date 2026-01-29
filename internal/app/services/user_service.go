package services

import (
	"backend-blog/internal/app/dao"
	"backend-blog/internal/middleware"
	"backend-blog/internal/model"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"backend-blog/utility"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	// 关键改动：注入 UserDao 而不是 gorm.DB
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

// Login 登录逻辑
func (s *UserService) Login(ctx context.Context, req dto.LoginReq) (*vo.LoginVo, error) {
	// ... existing logic ...
	user, err := s.userDao.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户不存在或查询失败")
	}

	if utility.Encryption(req.Password, user.Salt) != user.Password {
		return nil, errors.New("密码错误")
	}

	token, err := s.CreateToken(*user)
	if err != nil {
		slog.ErrorContext(ctx, "生成 Token 失败", "err", err)
		return nil, errors.New("系统错误")
	}

	return &vo.LoginVo{
		Token:    token,
		Username: user.Username,
	}, nil
}

// Add 注册用户
func (s *UserService) Add(ctx context.Context, req dto.RegisterReq) error {
	// 1. 业务校验：用户是否已存在
	exists, _ := s.userDao.CheckExists(ctx, req.Username)
	if exists {
		return errors.New("用户已存在")
	}

	// 2. 准备数据
	newPassword, salt := utility.GenerateNewPassword(req.Password)
	now := model.Now()
	user := &entity.User{
		BaseInfo: entity.BaseInfo{
			ID:          utility.GenID(),
			CreatedTime: now,
			UpdatedTime: now,
		},
		Username: req.Username,
		Password: newPassword,
		Salt:     salt,
	}

	// 3. 调用 DAO 写入
	return s.userDao.Create(ctx, user)
}

// ResetPassword 更新密码
func (s *UserService) ResetPassword(ctx context.Context, username, oldPassword, newPassword string) error {
	// 1. 验证旧密码
	user, err := s.userDao.GetByUsername(ctx, username)
	if err != nil {
		return errors.New("用户不存在")
	}

	if utility.Encryption(oldPassword, user.Salt) != user.Password {
		return errors.New("旧密码错误")
	}

	// 2. 生成新密码
	password, salt := utility.GenerateNewPassword(newPassword)

	// 3. 调用 DAO 的更新方法
	return s.userDao.UpdatePassword(ctx, username, password, salt)
}

// CreateToken 生成 JWT（内部逻辑保持不变）
func (s *UserService) CreateToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"uid":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(middleware.PrivateKey)
}
