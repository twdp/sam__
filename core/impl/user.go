package impl

import (
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"tianwei.pro/business"
	"tianwei.pro/micro/di"
	"tianwei.pro/micro/di/single"
	"tianwei.pro/sam/agent"
	"tianwei.pro/sam/core/const"
	"tianwei.pro/sam/core/dto"
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/dto/res"
	"tianwei.pro/sam/core/repository"
)

var (
	EmailOrPassErr = errors.New("账号或密码错误")
	UserNotActive  = errors.New("用户未激活或被冻结")
)

type UserFacadeImpl struct {
	UserRepository *repository.UserRepository `inject:"userRepository"`
	Agent          *SamCoreAgentImpl          `inject:"rpc_s_samAgentFacade"`
}

func init() {
	single.Provide(di.NewRpcProviderName("userFacade"), &UserFacadeImpl{})
}

func (u *UserFacadeImpl) Login(loginParam *req.EmailLoginDto) (reply *res.LoginDto) {
	reply = &res.LoginDto{}
	if user, err := u.UserRepository.FindByEmail(loginParam.Email); err != nil {
		reply.Error(err)
		return reply
	} else if _, err := business.ValidateCrypto(loginParam.Password, user.Password); err != nil {
		reply.Error(EmailOrPassErr)
		return reply
	} else {
		if user.Status != _const.Active {
			logs.Warn("user status not active, user: %s, orm user: %v", loginParam.Email, user)
			reply.Error(UserNotActive)
			return reply
		}
		userDto := &dto.UserDto{
			BaseDto: dto.BaseDto{
				Id: user.Id,
			},
			UserName: user.UserName,
		}
		if token, err := tokenFacadeImpl.EncodeToken(userDto, loginParam.Terminal); err != nil {
			reply.Error(err)
			return reply
		} else {
			reply.Token = token

			r := u.Agent.VerifyToken(&agent.VerifyTokenParam{
				SystemInfoParam: agent.SystemInfoParam{
					AppKey: loginParam.AppKey,
					Secret: loginParam.Secret,
				},
				Token: token,
			})
			reply.Response = r.Response

			reply.UserInfo = agent.UserInfo{
				Id:          r.Id,
				UserName:    r.UserName,
				Avatar:      r.Avatar,
				Email:       r.Email,
				Phone:       r.Phone,
				Permissions: r.Permissions,
			}

			return reply
		}
	}
}
