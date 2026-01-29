package dao

import (
	"backend-blog/internal/model/entity"
	"context"
	"time"
)

type UserDao struct {
	// 这里可以为空，或者存放其他配置
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// CheckExists 检查用户是否存在
func (d *UserDao) CheckExists(ctx context.Context, username string) (bool, error) {
	var count int64
	// 使用你提供的 GetDB(ctx) 获取实例
	err := GetDB(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// GetByUsername 查询用户
func (d *UserDao) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	//go func() {
	//	time.Sleep(3 * time.Second)
	//	pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
	//}()
	err := GetDB(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (d *UserDao) Create(ctx context.Context, user *entity.User) error {
	return GetDB(ctx).Create(user).Error
}

// UpdatePassword 更新密码
func (d *UserDao) UpdatePassword(ctx context.Context, username, pwd, salt string) error {
	return GetDB(ctx).Model(&entity.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"password":     pwd,
			"salt":         salt,
			"updated_time": time.Now().Unix(),
		}).Error
}
