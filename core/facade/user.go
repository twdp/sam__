package facade

import (
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/dto/res"
)

type RpcUserFacade struct {
	Login func(loginParam *req.EmailLoginDto) (*res.LoginDto, error)
}

type UserFacade interface {
	// 登录验证
	Login(loginParam *req.EmailLoginDto) *res.LoginDto
}
