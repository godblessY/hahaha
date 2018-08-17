package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"label/models"
	"strconv"
)

type CategoryController struct {
	beego.Controller
}

// @router /category/list [get]
func (c *CategoryController) GetCategoryList() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user.Id)
	categories := models.GetCategoryList()
	c.Data["categories"] = categories
	c.TplName = "category_list.tpl"
}

// @router /category/save [post]
func (c *OrganizerController) SaveCategory() {
	name := c.Input().Get("name")
	configjson := c.Input().Get("configjson")
	id, _ := strconv.ParseInt(c.Input().Get("category_id"), 10, 64)
	logs.Info(name)
	logs.Info(configjson)
	output := make(map[string]interface{})
	flag := int64(0)
	if id != 0 {
		flag = models.SaveCategory(id, name, configjson)
	} else {
		flag = models.AddCategory(name, configjson)
	}
	if flag != 0 {
		output["error_code"] = "0"
	} else {
		output["error_code"] = "1"
		output["message"] = "Failed"
	}

	c.Data["json"] = &output
	logs.Info(c.Data)
	c.ServeJSON()
}
