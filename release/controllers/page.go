package controllers

import (
	"release/models"
	"github.com/astaxie/beego"
)

type Page struct {
	NowPage    int //当前页
	UpPage     int //上一页
	NextPage   int //下一页
	Next2Page  int //下两页
	PageSize   int //每页多少数据
	TotalPage  int //总共多少页
	TotalCount int //总共多少条数据
}

func NewPage(env string, nowPage int) (page *Page) {
	qs, _ := models.ReleaseOrderConn()
	count, err := qs.Filter("Env", env).Count()
	if err != nil {
		beego.Error(err)
		return page
	}

	pageSize, _ := beego.AppConfig.Int("pageSize")
	totalCount := int(count)
	totalPage := totalCount / pageSize
	if totalCount % pageSize > 0 {
		totalPage = totalCount / pageSize + 1
	}
	page =  &Page{
		NowPage:    nowPage,
		UpPage:     nowPage - 1,
		NextPage:   nowPage + 1,
		Next2Page:  nowPage + 2,
		PageSize:   pageSize,
		TotalPage:  totalPage,
		TotalCount: totalCount,
	}
	return page
}
