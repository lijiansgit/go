package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"strings"
)

type BranchList struct {
	beego.Controller
}

func (c *BranchList) Get() {
	r := NewRelease()
	if err := r.branchMaster(); err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.TplName = "branch_list.html"
}

type BranchStatus struct {
	beego.Controller
}

func (c *BranchStatus) Get() {
	r := NewRelease()
	status, err := r.branchStatus()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}
	c.Data["status"] = status
	c.TplName = "branch_status.html"
}

type BranchLog struct {
	beego.Controller
}

func (c *BranchLog) Get() {
	r := NewRelease()
	log, err := r.branchLog()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.Data["log"] = log
	c.TplName = "branch_log.html"
}

type BranchPull struct {
	beego.Controller
}

func (c *BranchPull) Get() {
	r := NewRelease()
	res, err := r.branchPull()
	if err != nil {
		c.Data["status"] = 1
		res = strings.Split(fmt.Sprint(err), "\n")
		c.Data["res"] = res
	} else {
		c.Data["status"] = 0
		c.Data["res"] = res
	}
	c.TplName = "branch_pull.html"
}