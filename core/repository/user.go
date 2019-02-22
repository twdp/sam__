package repository

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"tianwei.pro/business"
	"tianwei.pro/micro/di/single"
	"tianwei.pro/sam/core/model"
)

var (
	UserNotExistErr = errors.New("邮箱不存在")
	FindUserErr     = errors.New("查找用户失败")
)

type UserRepository struct {
}

func init() {
	single.Provide("userRepository", &UserRepository{})
}

// 根据email查找用户信息
func (u *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{
		Email: email,
	}

	if err := orm.NewOrm().Read(user, "Email"); err != nil {
		if business.IsNoRowsError(err) {
			return nil, UserNotExistErr
		} else {
			logs.Error("find user by email failed. email: %s, err: %v", email, err)
			return nil, FindApiErr
		}
	} else {
		return user, nil
	}
}
