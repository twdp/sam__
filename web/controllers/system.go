package controllers

import (
	"github.com/pkg/errors"
	"net/http"
	"tianwei.pro/business"
	"tianwei.pro/business/controller"
	"tianwei.pro/sam/agent"
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/dto/res"
	"tianwei.pro/sam/core/facade"
)

type SystemController struct {
	controller.RestfulController
}

var (
	NoPermission = errors.New("无权限")
)

const systemListPermissionSession = "systemListPermissionSession_"

// 系统入住
// @router / [post]
func (s *SystemController) Stay() {
	if s.GetSession(agent.SamUserInfoSessionKey) == nil {
		s.Abort("401")
		return
	}
	u := s.GetSession(agent.SamUserInfoSessionKey).(*agent.UserInfo)

	stay := &req.SystemStay{}
	if business.IsError(s.ReadBody(stay)) {
		s.E500("参数错误")
	}

	stay.Operator = u.Id

	reply, err := facade.RpcSystem.Stay(stay)
	business.CheckError(&reply.Response, err)

	s.Return(reply)
}

// @router / [get]
func (s *SystemController) List() {
	reply := s.list()
	s.Return(reply.List)
}

func (s *SystemController) list() *res.SystemListResponse {
	u := s.GetSession(agent.SamUserInfoSessionKey).(*agent.UserInfo)

	reply, err := facade.RpcSystem.ListByOwner(u.Id)
	business.CheckError(&reply.Response, err)

	var sids []int64
	for _, s := range reply.List {
		sids = append(sids, s.Id)
	}

	s.SetSession(systemListPermissionSession, sids)
	return reply
}

// @router /api/:id [get]
func (s *SystemController) ListApis() {
	id := business.CastStringToInt64(s.Ctx.Input.Param(":id"))
	if business.IsError(s.checkPermission(id)) {
		s.Code(http.StatusUnauthorized, NoPermission.Error())
	}
	s.Return([]int64{})
}

func (s *SystemController) checkPermission(id int64) error {
	var systemIds []int64
	if ss, ok := s.GetSession(systemListPermissionSession).([]int64); !ok {
		s.list()
		if ss, ok = s.GetSession(systemListPermissionSession).([]int64); !ok {
			return nil
		}
		systemIds = ss
	} else {
		systemIds = ss
	}

	if len(systemIds) == 0 {
		return NoPermission
	}

	for _, sid := range systemIds {
		if id == sid {
			return nil
		}
	}

	return NoPermission
}
