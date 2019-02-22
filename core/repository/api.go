package repository

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"tianwei.pro/business"
	cache2 "tianwei.pro/sam/core/cache"
	"tianwei.pro/sam/core/const"
	"tianwei.pro/sam/core/model"
	"time"
)

var ApiRepositoryInstance = &ApiRepository{
	systemApiCache: cache2.NewCache(),
}

var (
	FindApiErr = errors.New("查询api列表失败")
)

type ApiRepository struct {
	systemApiCache cache.Cache
}

// 获取系统的url列表
func (a *ApiRepository) FindUrlBySystemId(systemId int64) ([]*model.Api, error) {
	allRoutes, err := a.FindAllApiBySystemId(systemId)
	if err != nil {
		return nil, err
	}
	var apis []*model.Api

	for _, route := range allRoutes {
		if route.Type == _const.Api {
			apis = append(apis, route)
		}
	}
	return apis, nil
}

// 获取系统的按钮\页面\url列表
func (a *ApiRepository) FindAllApiBySystemId(systemId int64) ([]*model.Api, error) {
	apis := a.systemApiCache.Get(business.CastInt64ToString(systemId))
	var allApis []*model.Api
	if apis == nil {
		if _, err := orm.NewOrm().QueryTable(&model.Api{}).Filter("SystemId", systemId).Filter("Status", _const.Active).All(&allApis); err != nil {
			logs.Error("query all api failed. systemId: %d, err: %v", systemId, err)
			return nil, FindApiErr
		} else {
			a.systemApiCache.Put(business.CastInt64ToString(systemId), allApis, 10*time.Minute)
			return allApis, nil
		}
	} else {
		allApis = apis.([]*model.Api)
	}
	return allApis, nil
}
