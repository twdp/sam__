package repository

import (
	"errors"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business"
	cache2 "tianwei.pro/sam/core/cache"
	"tianwei.pro/sam/core/model"
)

var SystemRepositoryInstance = &SystemRepository{
	appKeyCache: cache2.NewCache(),
}

var (
	AppKeyNotExistErr = errors.New("appKey不存在")
	LoadSystemInfoErr = errors.New("查询系统信息失败")
	SystemErr         = errors.New("系统错误")
	CreateSystemErr   = errors.New("系统入驻失败")
)

type SystemRepository struct {
	appKeyCache cache.Cache
}

func (s *SystemRepository) FindByAppKey(appKey string) (*model.System, error) {
	info := s.appKeyCache.Get(appKey)
	if info == nil {
		query := &model.System{
			AppKey: appKey,
		}
		if err := orm.NewOrm().Read(query, "AppKey"); err != nil {
			if business.IsNoRowsError(err) {
				return nil, AppKeyNotExistErr
			} else {
				logs.Error("query system info failed. appKey: %s, err: %v", appKey, err)
				return nil, LoadSystemInfoErr
			}

		} else {
			return query, nil
		}
	} else {
		return info.(*model.System), nil
	}
}

func (s *SystemRepository) Stay(system *model.System, userId int64) error {
	orm, err := business.TransactionStart()
	if err != nil {
		logs.Error("system stay err. system: %v, err: %v", system, err)

		return err
	}
	_, err = orm.Insert(system)
	if business.TransactionProcess(orm, err) {
		logs.Error("system stay err. system: %v, err: %v", system, err)
		return CreateSystemErr
	}

	mapping := &model.UserRole{
		UserId:   userId,
		RoleId:   1,
		SystemId: system.Id,
	}
	_, err = orm.Insert(mapping)
	if business.TransactionEnd(orm, err) {
		logs.Error(" stay err. system: %v, err: %v", system, err)
		return CreateSystemErr
	}
	return nil
}
