package model

import (
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business/model"
)

type UserRole struct {
	model.Base

	UserId int64

	RoleId int64

	SystemId int64

	// 如果系统角色和数据权限拉平
	// 数据权限放到此字段中
	// 前端存过来的都是选中的那一级
	BranchIds string
}

func init() {
	orm.RegisterModelWithPrefix("sam_", &UserRole{})
}