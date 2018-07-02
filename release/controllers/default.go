package controllers

import (
	"fmt"
	"release/models"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["webTitle"] = beego.AppConfig.String("webTitle")
	c.Data["username"] = c.GetSession("username")
	c.Data["adminuser"] = models.AdminUser
	c.TplName = "index.html"
}

type Welcome struct {
	beego.Controller
}

func (c *Welcome) Get() {
	c.TplName = "welcome.html"
}

type PreFile struct {
	beego.Controller
}

func (c *PreFile) Get() {
	c.TplName = "pre_file.html"
}

func (c *PreFile) Post() {
	f, header, err := c.GetFile("myfile")
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	f.Close()
	release := NewRelease()
	proFile := c.GetString("myfilePath")
	if proFile == "/" {
		proFile = release.prePath + "/"
	} else {
		proFile = release.prePath + proFile
	}
	preFileTmp := release.preFileTmp + "/" + header.Filename
	localFileTmp := "/tmp/" + header.Filename
	c.SaveToFile("myfile", localFileTmp)
	err = release.File(localFileTmp, preFileTmp, proFile)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprint(err))
		return
	}

	c.Ctx.WriteString("OK")
}
