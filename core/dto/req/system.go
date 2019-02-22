package req

// 系统入住
type SystemStay struct {
	// 系统名称
	Name string `json:"name"`

	Description string `json:"description"`

	// 系统头像
	Avatar string `json:"avatar"`

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

	// 操作人
	Operator int64
}
