package main

import (
	_ "beeblogs/routers"
	"beeblogs/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	orm.RunSyncdb("default", false, true)
	// 创建静态上传文件上传路径
	os.Mkdir("attachment", os.ModePerm)
	beego.Run()
}
