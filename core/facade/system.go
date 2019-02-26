package facade

import (
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/dto/res"
)

type RpcSystemFacade struct {
	// 系统入驻
	Stay func(stay *req.SystemStay) (*res.StayResponse, error)

	// 获取入驻的系统列表
	ListByOwner func(owner int64) (*res.SystemListResponse, error)
}


type SystemFacade interface {
	Stay(stay *req.SystemStay) (reply *res.StayResponse)
	ListByOwner(owner int64) (reply *res.SystemListResponse)
}
