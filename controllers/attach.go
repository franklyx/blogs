// Created by FrankLyx on 2018/2/6
// Desc :
package controllers

import (
	"net/url"
	"os"
	"io"

	"github.com/astaxie/beego"
)

type AttaController struct {
	beego.Controller
}

func (c *AttaController) Get() {
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:])
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	defer f.Close()

	_, err = io.Copy(c.Ctx.ResponseWriter, f)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

}
