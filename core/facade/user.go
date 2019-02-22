package facade

import (
	"context"
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/dto/res"
)

type UserFacade interface {
	// 登录验证
	Login(ctx context.Context, loginParam *req.EmailLoginDto, reply *res.LoginDto) error
}
