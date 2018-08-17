package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"label/models"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Get() {
	c.TplName = "login.tpl"
}

// @router /user/info [get]
func (c *UserController) GetInfo() {
	beego.Informational("get Info")
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	c.Data["json"] = current_user
	c.ServeJSON()
}

// @router /organizer/:organizer_id/user [get]
func (c *UserController) GetUserByOrganizer() {
	organizer_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":organizer_id"), 10, 64)
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	users := models.GetUsersByOrganizer(organizer_id)
	c.Data["users"] = users
	c.Data["organizer_id"] = organizer_id
	c.TplName = "organizer_users.tpl"
}

// @router /batch/:batch_id/users [get]
func (c *UserController) GetUserByBatch() {
	batch_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":batch_id"), 10, 64)
	logs.Info(batch_id)
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	users := models.GetUsersByBatchId(batch_id)
	all_users := models.GetAllUsers()
	c.Data["users"] = users
	c.Data["all_users"] = all_users
	c.Data["batch_id"] = batch_id
	c.TplName = "admin_batch_users.tpl"
}

// @router /batch/:batch_id/checkers [get]
func (c *UserController) GetCheckerByBatch() {
	batch_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":batch_id"), 10, 64)
	logs.Info(batch_id)
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	checkers := models.GetCheckerByBatchId(batch_id)
	all_checkers := models.GetAllCheckers()
	c.Data["checkers"] = checkers
	c.Data["all_checkers"] = all_checkers
	c.Data["batch_id"] = batch_id
	c.TplName = "admin_batch_checkers.tpl"
}

// @router /user/save [post]
func (c *UserController) SaveUser() {
	name := c.Input().Get("name")
	phone := c.Input().Get("phone")
	username := c.Input().Get("username")
	email := c.Input().Get("email")
	id, _ := strconv.ParseInt(c.Input().Get("user_id"), 10, 64)
	organizer_id, _ := strconv.ParseInt(c.Input().Get("organizer_id"), 10, 64)
	user_type, _ := strconv.Atoi(c.Input().Get("type"))
	output := make(map[string]interface{})
	flag := int64(0)
	if id != 0 {
		flag = models.SaveUser(id, name, phone, username, email, user_type)
	} else {
		password := c.Input().Get("password")
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(password))
		cipherStr := md5Ctx.Sum(nil)
		md5_pass := hex.EncodeToString(cipherStr)
		flag = models.AddUser(name, phone, username, email, user_type, organizer_id, md5_pass)
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
