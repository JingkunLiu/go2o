/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package product

import (
	"go2o/core/infrastructure/domain"
	"sort"
)

var (
	ErrReadonlyCategory *domain.DomainError = domain.NewDomainError(
		"err_readonly_category", "无权修改系统分类",
	)
	ErrNoSuchCategory *domain.DomainError = domain.NewDomainError(
		"err_category_not_exist", "分类不存在",
	)

	ErrCategoryCycleReference *domain.DomainError = domain.NewDomainError(
		"err_category_cycle_reference", "分类上级循环引用")

	ErrHasChildCategories *domain.DomainError = domain.NewDomainError(
		"err_has_child_categories", "分类包含子分类,无法删除",
	)
	ErrCategoryContainGoods *domain.DomainError = domain.NewDomainError(
		"err_category_contain_goods", "分类包含商品,无法删除",
	)
)

type (
	ICategory interface {
		// 获取领域编号
		GetDomainId() int32

		// 获取值
		GetValue() *Category

		// 设置值
		SetValue(*Category) error

		//todo: 做成界面,同时可后台管理项
		// 获取扩展数据
		GetOption() domain.IOptionStore

		// 保存
		Save() (int32, error)

		// 获取子栏目的编号
		GetChildes() []int32
	}
	//分类
	Category struct {
		Id int32 `db:"id" auto:"yes" pk:"yes"`
		//父分类
		ParentId int32 `db:"parent_id"`
		// 商品规格模型
		ProModel int32 `db:"pro_model"`
		//名称
		Name string `db:"name"`
		//层级,用于判断2个分类是否为同一级
		Level int `db:"level"`
		// 图标
		Icon string `db:"icon"`
		// 地址
		Url string `db:"url"`
		// 排序序号
		SortNum int `db:"sort_num"`
		// 是否启用,默认为不启用
		Enabled int `db:"enabled"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
	}
	ICategoryRepo interface {
		// 获取系统的栏目服务
		GlobCatService() IGlobCatService

		// 保存分类
		SaveCategory(*Category) (int32, error)

		// 检查分类是否关联商品
		CheckGoodsContain(mchId, id int32) bool

		// 删除分类
		DeleteCategory(mchId, id int32) error

		// 获取分类
		GetCategory(mchId, id int32) *Category

		// 获取所有分类
		GetCategories(mchId int32) []*Category
	}

	// 公共分类服务
	IGlobCatService interface {
		// 是否只读,当商户共享系统的分类时,
		// 没有修改的权限,即只读!
		ReadOnly() bool

		// 创建分类
		CreateCategory(*Category) ICategory

		// 获取分类
		GetCategory(id int32) ICategory

		// 获取所有分类
		GetCategories() []ICategory

		// 删除分类
		DeleteCategory(id int32) error
	}
)

var (
	C_OptionViewName string = "viewName" //显示的视图名称
	C_OptionDescribe string = "describe" //描述
)

var _ sort.Interface = new(CategoryList)

type CategoryList []*Category

func (c CategoryList) Len() int {
	return len(c)
}

func (c CategoryList) Less(i, j int) bool {
	return c[i].SortNum < c[j].SortNum ||
		// 如果序号相同,则判断ID
		(c[i].SortNum == c[j].SortNum && c[i].Id < c[j].Id)
}

func (c CategoryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
