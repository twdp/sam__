package dto

type UserDto struct {
	BaseDto

	UserName string `orm:"size(64);unique" json:"user_name"`

	DisplayName string `orm:"size(64)" json:"display_name"`

	Avatar string `json:"avatar"`

	Email string `orm:"size(64);unique" json:"email"`

	Phone string `orm:"size(64);unique" json:"phone"`

	Sex int8 `json:"sex"`

	Password string `json:"password"`

	Type int8 `json:"type"`

	Status int8 `json:"status"`

	// 退出时记录哪些端需要重新登录
	NeedLoginTerminus string `orm:"size(120)" json:"need_login_terminus"`
}
