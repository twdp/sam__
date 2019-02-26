package impl

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strings"
	"tianwei.pro/business"
	"tianwei.pro/micro/di"
	"tianwei.pro/micro/di/single"
	"tianwei.pro/sam/agent"
	"tianwei.pro/sam/core/model"
	"tianwei.pro/sam/core/repository"
)

var (
	AppKeyOrSecretError = errors.New("请检查appKey或secret")
	SystemError         = errors.New("权限系统错误")
)

type SamCoreAgentImpl struct {
	SystemRepository *repository.SystemRepository `inject:"systemRepository"`

	ApiRepository *repository.ApiRepository `inject:"apiRepository"`

	BranchRepository *repository.BranchRepository `inject:"branchRepository"`
}

func init() {
	agent := &SamCoreAgentImpl{
		//systemRepository: repository.SystemRepositoryInstance,
		//apiRepository:    repository.ApiRepositoryInstance,
		//branchRepository: repository.BranchRepositoryInstance,
	}
	single.Provide(di.NewRpcProviderName("samAgentFacade"), agent)
}

func (s *SamCoreAgentImpl) LoadSystemInfo(param *agent.SystemInfoParam) (reply *agent.SystemInfo) {
	reply = &agent.SystemInfo{
		Response: business.Response{
			Success: true,
		},
	}
	if systemInfo, err := s.verifySecret(param.AppKey, param.Secret); err != nil {
		reply.Error(err)
		return reply
	} else {
		routes, err := s.ApiRepository.FindUrlBySystemId(systemInfo.Id)
		if business.IsError(err) {
			reply.Error(err)
			return reply
		}
		var apis []*agent.Router

		for _, api := range routes {
			apis = append(apis, &agent.Router{
				Id:     api.Id,
				Url:    api.Path,
				Method: api.Method,
				Type:   api.VerificationType,
			})
		}

		reply.Id = systemInfo.Id
		reply.PermissionType = systemInfo.Strategy
		reply.KeepSign = systemInfo.KeepSign
		reply.Routers = apis

		return reply
	}
}

func (s *SamCoreAgentImpl) verifySecret(appKey, secret string) (*model.System, error) {
	if system, err := s.SystemRepository.FindByAppKey(appKey); err != nil {
		logs.Warn("app key: %v not found", appKey)
		return nil, AppKeyOrSecretError
	} else {
		if system.Secret != secret {
			logs.Warn("param: %v, system info: %v", secret, system)
			return nil, AppKeyOrSecretError
		} else {

		}
		return system, nil
	}
}

func (s *SamCoreAgentImpl) VerifyToken(param *agent.VerifyTokenParam) (reply *agent.UserInfo) {
	reply = &agent.UserInfo{}
	system, err := s.verifySecret(param.AppKey, param.Secret)
	if err != nil {
		reply.Error(err)
		return reply
	}
	if user, err := tokenFacadeImpl.DecodeToken(param.Token); err != nil {
		reply.Error(err)
		return reply
	} else if err := s.loadUserInfo(user.Id, reply); err != nil {
		reply.Error(err)
		return reply
	}

	// todo:: add cache

	var userRoles []*model.UserRole
	if _, err := orm.NewOrm().QueryTable(&model.UserRole{}).Filter("SystemId", system.Id).Filter("UserId", reply.Id).All(&userRoles); err != nil {
		logs.Error("read user role mapping failed. userInfo: %v, system: %v, err: %v", reply, system, err)
		reply.Error(SystemError)
		return reply
	}
	if len(userRoles) == 0 {
		return reply
	}
	var roleIds []int64
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole.Id)
	}

	var roles []*model.Role
	if _, err := orm.NewOrm().QueryTable(&model.Role{}).Filter("id__in", roleIds).All(&roles); err != nil {
		logs.Error("load roles failed. ids: %v, err: %v", roleIds, err)
		reply.Error(SystemError)
		return reply
	}

	roleIdBranchIds := make(map[int64]string)
	for _, userRole := range userRoles {
		if userRole.BranchIds == "" {
			continue
		}
		roleIdBranchIds[userRole.RoleId] = userRole.BranchIds
	}

	var permissions []*agent.Permission
	for _, role := range roles {
		ps := strings.Split(role.PermissionSet, ",")
		var psCastIds []int64
		for _, psId := range ps {
			if psId == "" {
				continue
			}
			psCastIds = append(psCastIds, business.CastStringToInt64(psId))
		}
		permission := &agent.Permission{
			RoleId:        role.Id,
			RoleName:      role.Name,
			PermissionSet: psCastIds,
		}

		if system.Strategy == agent.OnlyOperationPermission {
			// 仅操作权限
		} else if system.Strategy == agent.Equality {
			if branchIds, exist := roleIdBranchIds[role.Id]; exist {
				ps := strings.Split(branchIds, ",")
				var psCastIds []int64
				for _, psId := range ps {
					if psId == "" {
						continue
					}
					psCastIds = append(psCastIds, business.CastStringToInt64(psId))
				}
				permission.BranchIds = psCastIds
			}
		} else {
			// 树形权限模型
			branchIds := s.BranchRepository.RecursionBranchIds(role.BranchId)
			if role.BranchId != 0 {
				branchIds = append(branchIds, role.BranchId)
			}
			permission.BranchIds = branchIds
		}

		permissions = append(permissions, permission)
	}

	reply.Permissions = permissions
	return reply
}

func (s *SamCoreAgentImpl) loadUserInfo(userId int64, reply *agent.UserInfo) error {
	u := &model.User{}
	u.Id = userId
	if err := orm.NewOrm().Read(u); err != nil {
		logs.Error("find user by token failed. err: %v", err)
		reply.Error(SystemError)
		return SystemError
	}
	reply.Id = u.Id
	reply.UserName = u.UserName
	reply.Avatar = u.Avatar
	reply.Email = u.Email
	reply.Phone = u.Phone

	return nil
}
