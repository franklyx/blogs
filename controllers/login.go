// Created by FrankLyx on 2018/2/1
// Desc :
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	isExit := c.Input().Get("exit") == "true"
	if isExit {
		c.Ctx.SetCookie("uname", "", "-1", "/")
		c.Ctx.SetCookie("upwd", "", "-1", "/")
		c.Redirect("/", 302)
		return
	}
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	uname := c.Input().Get("uname")
	upwd := c.Input().Get("upwd")
	remember := c.Input().Get("remember") == "on"

	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("upwd") == upwd {
		maxAge := 0
		if remember {
			maxAge = 1<<31 - 1
		}
		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("upwd", upwd, maxAge, "/")

	}
	c.Redirect("/", 302)
	return
}

func checkoutAccount(ctx *context.Context) bool {
	uname := ctx.GetCookie("uname")
	upwd := ctx.GetCookie("upwd")

	return beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("upwd") == upwd

}
