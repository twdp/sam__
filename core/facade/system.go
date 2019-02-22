package facade

import (
	"context"
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/dto/res"
)

type SystemFacade interface {
	Stay(context context.Context, stay *req.SystemStay, reply *res.StayResponse) error
	ListByOwner(context context.Context, owner int64, reply *res.SystemListResponse) error
}
