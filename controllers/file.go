package controllers

import (
    "os"
    "strconv"
    "label/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type FileUploadController struct {
    beego.Controller
}

// @router /file/upload [get]
func (c *FileUploadController) UploadIndex() {
    current_user_session := c.GetSession("label_me_user")
    if current_user_session == nil {
	    logs.Info("session lost")
    	beego.Informational("session lost")
    	c.Redirect("/login", 302)
    	return
    }
    c.TplName = "file_upload.tpl"
}

// @router /file/upload [post]
func (c *FileUploadController) Upload() {
    current_user_session := c.GetSession("label_me_user")
    if current_user_session == nil {
	    logs.Info("session lost")
    	beego.Informational("session lost")
    	c.Redirect("/login", 302)
    	return
    }
    current_user := current_user_session.(models.User)
    upload_folder := beego.AppConfig.String("upload_folder") + "/" + strconv.FormatInt(current_user.Id, 10)
	logs.Info(upload_folder)
    if _, err := os.Stat(upload_folder); os.IsNotExist(err) {
        os.Mkdir(upload_folder, 0700)
    }
	f, h, _ := c.GetFile("myfile")
	path := upload_folder + "/" + h.Filename
	logs.Info(path)
	f.Close()
	c.SaveToFile("myfile", path)
    c.TplName = "file_upload_sucess.tpl"
}
