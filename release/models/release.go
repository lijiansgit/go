package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	AdminUser     = beego.AppConfig.String("adminUser")
	AdminPassword = beego.AppConfig.String("adminPassword")
)

type ReleaseUser struct {
	Id       int
	Username string `orm:"unique"`
	Password string
}

type ReleaseOrder struct {
	Id        int
	Env       string `orm:"index"` //发布环境，预发布：pre,线上：pro
	Title     string //工单主题
	OpType    string //工单类型:发布或者回滚
	OpName    string //发布工单用户名
	Status    bool   //工单状态
	Timestamp int    `orm:"index;unique"` //创建工单时间戳
	FileLog   string `orm:"type(text)"`   //发布文件日志
}

func init() {
	// set default database
	//orm.RegisterDataBase("default", "mysql", "test:123456@tcp(10.10.10.11:3306)/test?charset=utf8", 30)
	orm.RegisterDataBase("default", "sqlite3", "models/data.db")

	// register model
	orm.RegisterModel(new(ReleaseUser), new(ReleaseOrder))

	// create table
	orm.RunSyncdb("default", false, true)

	// init admin user
	admin := &ReleaseUser{Username: AdminUser, Password: AdminPassword}
	o := orm.NewOrm()
	qs := o.QueryTable(admin)
	if qs.Filter("Username", AdminUser).Exist() == false {
		o.Insert(admin)
	}
}

func ReleaseUserConn() (qs orm.QuerySeter, user *ReleaseUser) {
	o := orm.NewOrm()
	user = new(ReleaseUser)
	qs = o.QueryTable(user)
	return qs, user
}

func ReleaseOrderConn() (qs orm.QuerySeter, order *ReleaseOrder) {
	o := orm.NewOrm()
	order = new(ReleaseOrder)
	qs = o.QueryTable(order)
	return qs, order
}
