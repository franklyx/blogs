package controllers

import (
	"github.com/astaxie/beego"
	"beeblogs/models"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Data["IsHome"] = true
	c.Data["IsLogin"] = checkoutAccount(c.Ctx)

	cate := c.Input().Get("cate")
	label := c.Input().Get("label")
	topics, err := models.GetAllTopics(cate, label, true)

	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"] = topics

	// 添加分类
	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Categories"] = categories

	c.TplName = "home.html"

}
