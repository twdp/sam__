package model

import (
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business/model"
)

// 第三方授权表
type UserThirdAuthorization struct {
	model.Base

	UserId int64

	OpenId string  `orm:"size(64)"`

	UnionId string `orm:"size(64)"`

	Channel int8
}

// 多字段唯一键
func (u *UserThirdAuthorization) TableUnique() [][]string {
	return [][]string{
		{ "OpenId", "Channel", },
	}
}

func init() {
	orm.RegisterModelWithPrefix("sam_", new(UserThirdAuthorization))
}