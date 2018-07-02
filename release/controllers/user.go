package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	models "release/models"
	"fmt"
	"github.com/astaxie/beego/orm"
)

func init() {
	CheckUsers()
}

func CheckUsers() {
	var Check = func(ctx *context.Context) {
		username := ctx.Input.Session("username")
		if username == nil {
			ctx.Redirect(302, "/user/login")
			return
		}
	}
	beego.InsertFilter("/release/*", beego.BeforeRouter, Check)
}

type Login struct {
	beego.Controller
	username		string
	password		string
	postUsername	string
	postPassword	string

}

func (c *Login) Get() {
	username := c.GetSession("username")
	if username != nil {
		c.Redirect("/release", 302)
	}
	c.TplName = "login.html"
}

func (c *Login) Post() {
	c.postUsername = c.GetString("username")
	c.username = c.postUsername
	c.postPassword = c.GetString("password")
	qs, user := models.ReleaseUserConn()
	if err := qs.Filter("Username", c.username).One(user); err != nil {
		beego.Debug("User Login: ", err)
		c.Redirect("/user/login", 302)
	}

	if user.Password == c.postPassword {
		c.SetSession("username", c.username)
		beego.Debug("User Login: ", c.username)
		c.Redirect("/release", 302)
	}

	c.Redirect("/user/login", 302)
}

type Logout struct {
	beego.Controller
}

func (c *Logout) Get() {
	username := c.GetSession("username")
	if username != nil {
		c.DelSession("username")
	}

	c.Redirect("/user/login", 302)
}

// admin user manange
type UserList struct {
	beego.Controller
}

func (c *UserList) Prepare() {
	username := c.GetSession("username")
	if username != models.AdminUser {
		c.Ctx.WriteString("NO AUTH")
	}
}

func (c *UserList) Get() {
	qs, _ := models.ReleaseUserConn()
	var users []*models.ReleaseUser
	if _, err := qs.All(&users); err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.Data["users"] = users
	c.TplName = "user_list.html"
}

type UserAdd struct {
	beego.Controller
}

func (c *UserAdd) Prepare() {
	username := c.GetSession("username")
	if username != models.AdminUser {
		c.Ctx.WriteString("NO AUTH")
	}
}

func (c *UserAdd) Get() {
	c.TplName = "user_add.html"
}

func (c *UserAdd) Post() {
	username := c.GetString("username")
	password := c.GetString("password")
	qs, user := models.ReleaseUserConn()
	if qs.Filter("Username", username).Exist() == true {
		c.Ctx.WriteString("user is exists")
		return
	}

	user.Username = username
	user.Password = password
	o := orm.NewOrm()
	if _, err := o.Insert(user); err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	beego.Info("User Add, username: ", username)
	c.Redirect("/user/list", 302)
}


type UserDel struct {
	beego.Controller
}

func (c *UserDel) Prepare() {
	username := c.GetSession("username")
	if username != models.AdminUser {
		c.Ctx.WriteString("NO AUTH")
	}
}

func (c *UserDel) Get() {
	username := c.Ctx.Input.Param(":username")
	if username == models.AdminUser {
		c.Ctx.WriteString("ADMIN user cannot del")
		return
	}

	qs, _ := models.ReleaseUserConn()
	if _, err := qs.Filter("Username", username).Delete(); err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	beego.Info("User Del, username: ", username)
	c.Redirect("/user/list", 302)
}


