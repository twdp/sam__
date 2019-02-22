package res

import (
	"tianwei.pro/business"
	"time"
)

// 系统入驻
type StayResponse struct {
	business.Response

	Id int64 `json:"id"`

	Name string `json:"name"`

	Status int8

	Description string `json:"description"`

	Avatar string `json:"avatar"`

	AppKey string `json:"app_key"`

	Secret string `json:"secret"`

	Strategy int8 `json:"strategy"`
}


type SystemListResponse struct {
	business.Response

	List []*SystemList `json:"list"`
}

type SystemList struct {
	
	Id int64 `json:"id"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	Status int8 `json:"status"`

	Name string `json:"name"`

	Description string `json:"description"`

	// 系统头像
	Avatar string `json:"avatar"`

	AppKey string `json:"app_key"`

	Secret string `json:"secret"`

	// 是否需要保持登录
	KeepSign bool `json:"keep_sign"`

	// 0 不需要数据权限
	// 1 数据权限与操作 权限拉平
	// 2 层级数据权限和操作权限
	Strategy int8 `json:"strategy"`

	// 使用了模板角色，是否对模板角色可见
	// 管理员手工给员工改角色时可见
	// 1. 电商系统：创建店铺时，会自动创建店铺管理员，此时所有的店铺管理员角色均是模板角色，这种模板角色在用户申请角色时是否可见
	// 2. 三方系统接入权限系统时，owner角色为所有系统的模板角色，是否可见
	// 3. 如果本系统自己会创建模板角色给系统内部用
	TemplateRoleVisible bool `json:"template_role_visible"`
}