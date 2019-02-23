package repository

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"tianwei.pro/business"
	"tianwei.pro/micro/di/single"
	cache2 "tianwei.pro/sam/core/cache"
	"tianwei.pro/sam/core/dto/res"
	"tianwei.pro/sam/core/model"
	"time"
)

var (
	FindBranchError = errors.New("获取数据权限失败")
)

var BranchRepositoryInstance = &BranchRepository{
	branchIdCache: cache2.NewCache(),
	branchCache:   cache2.NewCache(),
}

func init() {
	single.Provide("branchRepository", BranchRepositoryInstance)
}
type BranchRepository struct {
	branchIdCache cache.Cache
	branchCache   cache.Cache
}

func (b *BranchRepository) RecursionBranchIds(pid int64) []int64 {
	var childrenIds []int64
	if b.branchIdCache.IsExist(business.CastInt64ToString(pid)) {
		branchIds := b.branchIdCache.Get(business.CastInt64ToString(pid)).([]int64)
		return branchIds
	}
	var childrens []*model.Branch
	orm.NewOrm().QueryTable(&model.Branch{}).Filter("ParentId", pid).All(&childrens)
	for _, children := range childrens {
		childrenIds = append(childrenIds, children.Id)
		childrenIds = append(childrenIds, b.RecursionBranchIds(children.Id)...)
	}
	b.branchIdCache.Put(business.CastInt64ToString(pid), childrenIds, time.Minute*30)
	return childrenIds
}

func (b *BranchRepository) RecursionBranchTree(pid int64, bt *res.BranchTree) {

	if b.branchCache.IsExist(business.CastInt64ToString(pid)) {
		btree := b.branchCache.Get(business.CastInt64ToString(pid)).(*res.BranchTree)
		bt.Id = btree.Id
		bt.Name = btree.Name
		bt.CreatedAt = btree.CreatedAt
		bt.UpdatedAt = btree.UpdatedAt
		bt.Childrens = btree.Childrens
		return
	}

	var childrens []*model.Branch
	orm.NewOrm().QueryTable(&model.Branch{}).Filter("ParentId", pid).All(&childrens)

	for _, children := range childrens {
		branchTree := &res.BranchTree{
			Id:        children.Id,
			Name:      children.Name,
			CreatedAt: children.CreatedAt.Value(),
			UpdatedAt: children.UpdatedAt.Value(),
		}

		b.RecursionBranchTree(children.Id, branchTree)

		bt.Childrens = append(bt.Childrens, branchTree)
	}

	b.branchCache.Put(business.CastInt64ToString(pid), bt, time.Minute*30)
}
