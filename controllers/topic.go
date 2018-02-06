// Created by FrankLyx on 2018/2/2
// Desc :
package controllers

import (
	"github.com/astaxie/beego"
	"beeblogs/models"
	"strings"
	"path"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkoutAccount(c.Ctx)
	c.TplName = "topic.html"

	topics, err := models.GetAllTopics("", "", false)

	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"] = topics

}

func (c *TopicController) Post() {
	c.Data["IsTopic"] = true
	if !checkoutAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	tid := c.Input().Get("tid")
	title := c.Input().Get("title")
	label := c.Input().Get("label")
	content := c.Input().Get("content")
	category := c.Input().Get("category")

	// 获取附件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}

	var attachment string
	// 判断是否上传附件
	if fh != nil {
		attachment = fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	println(tid)
	if len(tid) == 0 {
		err = models.AddTopic(title, category, label, content, attachment)
	} else {
		err = models.ModifyTopic(tid, title, category, label, content, attachment)
	}

	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic", 302)
	return

}

func (c *TopicController) Add() {
	c.Data["IsTopic"] = true
	c.TplName = "topic_add.html"
}

func (c *TopicController) View() {
	tid := c.Ctx.Input.Param("0")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid

	replies, err := models.GetAllReplies(tid)
	if err != nil {
		beego.Error(err)
		return
	}

	c.Data["Labels"] = strings.Split(topic.Labels, " ")
	c.Data["Replies"] = replies
	c.TplName = "topic_view.html"
}

func (c *TopicController) Modify() {
	tid := c.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
	c.TplName = "topic_modify.html"
}

func (c *TopicController) Delete() {
	if !checkoutAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	err := models.DelTopic(c.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/", 302)
	return

}
