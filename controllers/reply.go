// Created by FrankLyx on 2018/2/6
// Desc :
package controllers

import (
	"github.com/astaxie/beego"
	"beeblogs/models"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	tid := c.Input().Get("tid")
	nickname := c.Input().Get("nickname")
	content := c.Input().Get("content")

	err := models.AddReply(tid, nickname, content)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 302)
	return
}

func (c *ReplyController) Delete() {
	if !checkoutAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	tid := c.Input().Get("tid")
	rid := c.Input().Get("rid")
	err := models.DeleteReply(rid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 302)
	return
}
