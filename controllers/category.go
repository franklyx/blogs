// Created by FrankLyx on 2018/2/1
// Desc :
package controllers

import (
	"github.com/astaxie/beego"
	"beeblogs/models"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.Data["IsCategory"] = true

	cates, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Categories"] = cates

	c.TplName = "category.html"

	op := c.Input().Get("op")

	switch op {
	case "add":
		name := c.Input().Get("name")
		if len(name) == 0 {
			break
		}

		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}

		c.Redirect("/category", 302)

	case "del":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}

		err := models.DelCategory(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 302)
		return

	}
}
