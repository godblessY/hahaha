package controllers

import (
    "crypto/md5"
    "encoding/hex"
    "label/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "math/rand"
//    "strconv"
)

type LoginController struct {
	beego.Controller
}

// @router /login [get]
func (c *LoginController) Get() {
    c.TplName = "login.tpl"
}

func randInt(min int , max int) int {
    return min + rand.Intn(max-min)
}

// @router /login [post]
func (c *LoginController) Post() {
    username:= c.Input().Get("username")
    pwd:= c.Input().Get("pwd")
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(pwd))
    cipherStr := md5Ctx.Sum(nil)
    md5_pass := hex.EncodeToString(cipherStr)
    output := make(map[string]interface{})
    data := models.GetUserInfo(username, md5_pass)
    logs.Info(data)
    if data.Id != 0  {
        output["error_code"] = "0"
        output["user_type"] = data.Type
	    c.SetSession("label_me_user", data)
    } else {
        output["error_code"] = "1"
        output["message"] = "用户不存在，请联系管理员"
    }

    c.Data["json"] = &output
    logs.Info(c.Data)
    c.ServeJSON()
}

// @router /logout [get]
func (c *LoginController) Logout() {
	c.DelSession("label_me_user")
	c.Redirect("/login", 302)
}
