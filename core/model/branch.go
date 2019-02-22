package model

import (
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business/model"
)

// 组织架构表
type Branch struct {

	model.Base

	// 名称
	Name string

	// 父id
	ParentId int64
}

// 多字段唯一键
func (b *Branch) TableUnique() [][]string {
	return [][]string{
		{ "Name", "ParentId", },
	}
}

func init() {
	orm.RegisterModelWithPrefix("sam_", new(Branch))
}