package _const

// 状态
const (
	Init = iota
	Active
	Freeze
	Delete
)

// url类型
const (
	Menu = iota
	Page
	Api
)

const (
	Official = iota // 默认
	WxMp            // 微信小程序
)
