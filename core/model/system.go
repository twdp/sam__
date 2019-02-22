package model

import (
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business/model"
)

// 对接的系统
type System struct {
	model.Base

	Status int8

	Name string `orm:"size(100);unique"`

	Description string

	// 系统头像
	Avatar string

	AppKey string `orm:"size(64);unique"`

	Secret string

	// 是否需要保持登录
	KeepSign bool

	// 0 不需要数据权限
	// 1 数据权限与操作 权限拉平
	// 2 层级数据权限和操作权限
	Strategy int8

	// 使用了模板角色，是否对模板角色可见
	// 管理员手工给员工改角色时可见
	// 1. 电商系统：创建店铺时，会自动创建店铺管理员，此时所有的店铺管理员角色均是模板角色，这种模板角色在用户申请角色时是否可见
	// 2. 三方系统接入权限系统时，owner角色为所有系统的模板角色，是否可见
	// 3. 如果本系统自己会创建模板角色给系统内部用
	TemplateRoleVisible bool

	// todo:: 使用模板角色，拷贝一份出来，然后让其可修改
	// 当前：要么完全自己设置角色权限，要么使用模板角色，但不允许修改
	// 系统对接设置以后，尽量不能修后面这3个字段，如果需要修改，需要清除所有设置
}

func init() {
	orm.RegisterModelWithPrefix("sam_", new(System))
}
