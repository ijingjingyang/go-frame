package entry

import (
	"fmt"
	"strings"
)

// 通用entry结构封装

// PageReq 分页通用请求参数
type PageReq struct {
	// 通用参数--第几页
	Page int `json:"page" form:"page" binding:"omitempty"`

	// 通用参数--每页数量
	Limit int `json:"limit" form:"limit" binding:"omitempty,max=1000"`
}

// GetOffset 获取偏移量
func (s PageReq) GetOffset() int {
	return (s.GetPage() - 1) * s.GetLimit()
}

// GetPage 获取页码，不能小于0
func (s PageReq) GetPage() int {
	if s.Page <= 0 {
		return 1
	}
	return s.Page
}

// GetLimit 获取显示条数 不能小于0 不能大于1000
func (s PageReq) GetLimit() int {
	if s.Limit <= 0 {
		return 10
	} else if s.Limit > 1000 {
		return 1000
	} else {
		return s.Limit
	}
}

// SortReq 通用排序请求参数
type SortReq struct {
	// 通用参数--排序字段名称
	OrderBy string `form:"order_by" json:"order_by"`
	// 通用参数--排序类型<asc--升序 desc--降序>
	// enum:desc,asc
	Sort string `form:"sort" json:"sort" binding:"omitempty,oneof=desc asc"`
}

// GetOrderBy 组装排序参数
func (s SortReq) GetOrderBy() string {
	s.OrderBy = strings.Replace(s.OrderBy, "`", "", -1) // 去除恶意的字段反引号
	if s.OrderBy != "" && s.Sort != "" {
		return fmt.Sprintf("`%s` %s", s.OrderBy, s.Sort)
	}
	return "`id` desc"
}

// PageRes 分页响应数据
type PageRes struct {
	// 总数 全量返回或者不需要数量得情况返回0
	Total int64 `json:"total"`
	// 列表数据
	List interface{} `json:"list"`
}

// NumRes 返回一个数字
type NumRes struct {
	Num int64 `json:"num"`
}

// BaseRes 基本响应
type BaseRes struct {
	// 错误码 0成功 非0失败
	Code int64 `json:"code"`

	// 错误码对应的错误信息
	Msg string `json:"msg"`
}
