package controllers

import (
	"fmt"
	"release/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OrderList struct {
	beego.Controller
}

func (c *OrderList) Get() {
	env := c.Ctx.Input.Param(":env")
	page := c.Ctx.Input.Param(":page")
	nowPage, err := strconv.Atoi(page)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	Page := NewPage(env, nowPage)
	qs, _ := models.ReleaseOrderConn()
	var orders []*models.ReleaseOrder
	if nowPage == 1 {
		_, err = qs.Filter("Env", env).OrderBy("-Id").Limit(Page.PageSize).All(&orders)
	} else {
		_, err = qs.Filter("Env", env).OrderBy("-Id").Limit(Page.PageSize,
			Page.PageSize*(nowPage-1)).All(&orders)
	}
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.Data["env"] = env
	c.Data["Page"] = Page
	c.Data["orders"] = orders
	c.TplName = "order_list.html"
}

type OrderLog struct {
	beego.Controller
}

func (c *OrderLog) Get() {
	timestamp := c.Ctx.Input.Param(":timestamp")
	qs, order := models.ReleaseOrderConn()
	err := qs.Filter("Timestamp", timestamp).One(order)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	log := strings.Split(order.FileLog, "\n")
	c.Data["log"] = log
	c.TplName = "log.html"
}

type OrderSearch struct {
	beego.Controller
}

func (c *OrderSearch) Post() {
	env := c.GetString("env")
	field := c.GetString("field")
	keywords := c.GetString("keywords")
	if keywords == "" {
		c.Ctx.WriteString("搜索关键字不能为空")
		return
	}

	var orders []*models.ReleaseOrder
	qs, _ := models.ReleaseOrderConn()
	qs.Filter("Env", env).Filter(fmt.Sprintf("%s__icontains", field), keywords).All(&orders)
	c.Data["orders"] = orders
	c.TplName = "order_list_search.html"
}

type OrderListFalse struct {
	beego.Controller
}

func (c *OrderListFalse) Get() {
	qs, _ := models.ReleaseOrderConn()
	var orders []*models.ReleaseOrder
	_, err := qs.Filter("Status", false).All(&orders)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.Data["orders"] = orders
	c.TplName = "order_list_false.html"
}

type OrderModify struct {
	beego.Controller
}

func (c *OrderModify) Get() {
	orderId := c.Ctx.Input.Param(":id")
	qs, _ := models.ReleaseOrderConn()
	_, err := qs.Filter("Id", orderId).Update(orm.Params{"Status": true})
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.Redirect("/release/order/list/false", 302)
}

type OrderAdd struct {
	beego.Controller
}

func (c *OrderAdd) Get() {
	c.TplName = "order_add.html"
}

func (c *OrderAdd) Post() {
	qs, order := models.ReleaseOrderConn()
	if qs.Filter("Status", false).Exist() == true {
		c.Ctx.WriteString("有未发布完成的工单，无法新增工单!")
		return
	}

	postEnv := c.GetString("env")
	if postEnv != "pre" && postEnv != "pro" {
		c.Ctx.WriteString("pre or pro")
		return
	}

	postTitle := c.GetString("title")
	if len(postTitle) < 3 {
		c.Ctx.WriteString("主题不能小于三个字符,请重新输入")
		return
	}

	opName := fmt.Sprint(c.GetSession("username"))
	order.Env = postEnv
	order.Title = postTitle
	order.OpType = "发布"
	order.OpName = opName
	order.Status = false
	order.Timestamp = int(time.Now().Unix())
	o := orm.NewOrm()
	if _, err := o.Insert(order); err != nil {
		beego.Error(err)
		return
	}

	release := NewRelease()
	if postEnv == "pre" {
		go release.Pre(qs, order)
	}
	if postEnv == "pro" {
		go release.Pro(qs, order)
	}
	redirect := "/release/order/list/" + postEnv + "/1"
	c.Redirect(redirect, 302)
}

type OrderBack struct {
	beego.Controller
}

func (c *OrderBack) Get() {
	backTimestamp := c.Ctx.Input.Param(":timestamp")
	qs, order := models.ReleaseOrderConn()

	err := qs.Filter("Timestamp", backTimestamp).One(order)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}
	c.Data["order"] = order
	c.TplName = "order_back.html"
}

func (c *OrderBack) Post() {
	qs, order := models.ReleaseOrderConn()
	if qs.Filter("Status", false).Exist() == true {
		c.Ctx.WriteString("有未发布完成的工单，无法新增工单!")
		return
	}

	backTimestamp := c.Ctx.Input.Param(":timestamp")
	backTimestampInt, err := strconv.Atoi(backTimestamp)
	if err != nil {
		beego.Error(err)
		return
	}

	postTitle := c.GetString("title")
	if len(postTitle) < 3 {
		c.Ctx.WriteString("主题不能小于三个字符,请重新输入")
		return
	}

	if err := qs.Filter("Timestamp", backTimestampInt).Filter("Env", "pre").One(order); err != nil {
		c.Ctx.WriteString("此工单不存在!")
		return
	}
	backOrderId := order.Id
	opName := fmt.Sprint(c.GetSession("username"))
	order = new(models.ReleaseOrder)
	order.Env = "pre"
	order.Title = fmt.Sprintf("回滚目标工单号: %d\n回滚主题: \n%s", backOrderId, postTitle)
	order.OpType = "回滚"
	order.OpName = opName
	order.Status = false
	order.Timestamp = int(time.Now().Unix())
	o := orm.NewOrm()
	if _, err := o.Insert(order); err != nil {
		beego.Error(err)
		return
	}

	release := NewRelease()
	go release.PreBack(backTimestamp, qs, order.Timestamp)

	c.Redirect("/release/order/list/pre/1", 302)
}
