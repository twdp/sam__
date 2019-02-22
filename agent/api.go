package agent

import (
	"tianwei.pro/business"
	"tianwei.pro/micro/di"
	"tianwei.pro/micro/di/single"
)

// 系统信息
type SystemInfo struct {
	business.Response

	Id int64

	// 权限类型
	PermissionType int8

	// 是否使用token保持登录
	KeepSign bool

	// 配置在sam中的url列表
	Routers []*Router
}

// url信息
type Router struct {
	// 系统url id
	Id int64

	// url
	Url string

	// method
	Method string

	// url类型
	Type int8
}

// 请求系统信息的参数
type SystemInfoParam struct {
	// app key
	AppKey string

	// secret
	Secret string
}

// 校验token的参数
type VerifyTokenParam struct {
	SystemInfoParam

	Token string
}

type UserInfo struct {
	business.Response

	// 用户id
	Id int64

	// 用户名
	UserName string

	// 用户头像
	Avatar string

	// 邮箱  可用于发送验证信息之类的
	Email string

	// 手机号  可用于发送验证信息之类的
	Phone string

	Permissions []*Permission
}

func init() {
	single.Provide(di.NewRpcConsumerName("samAgentFacade"), &SamAgentFacade{})
}
type SamAgentFacade struct {

	// 根据appKey 和secret获取系统信息
	LoadSystemInfo func(param *SystemInfoParam) (*SystemInfo, error)

	// 验证token
	VerifyToken func(param *VerifyTokenParam) (*UserInfo, error)
}