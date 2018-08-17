package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"label/models"
	"strconv"
)

type OrganizerController struct {
	beego.Controller
}

// @router /organizer/list [get]
func (c *OrganizerController) GetOrganizerList() {
	logs.Info("get organizer list")
	beego.Informational("get organizer list")
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	organizers := models.GetAllOrganizer()
	logs.Info(organizers)
	c.Data["organizers"] = organizers
	c.TplName = "organizer_list.tpl"
}

// @router /organizer/save [post]
func (c *OrganizerController) SaveOrganizer() {
	name := c.Input().Get("name")
	aliase := c.Input().Get("aliase")
	id, _ := strconv.ParseInt(c.Input().Get("organizer_id"), 10, 64)
	logs.Info(name)
	logs.Info(aliase)
	logs.Info(id)
	output := make(map[string]interface{})
	flag := int64(0)
	if id != 0 {
		flag = models.SaveOrganizer(id, name, aliase)
	} else {
		flag = models.AddOrganizer(name, aliase)
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
