package agent


// 系统数据权限类型
const (
	OnlyOperationPermission = iota
	Equality
	Child
)

// url检查策略
const (
	Anonymous = iota
	OnlyNeedLogin
	CheckRolePermission
)