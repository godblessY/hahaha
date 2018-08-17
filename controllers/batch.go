package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"label/models"
	"strconv"
)

type BatchController struct {
	beego.Controller
}

type BatchView struct {
	Id          int64
	BatchName   string
	Created     string
	Sum         int64
	Finished    int64
	Rest        int64
	Passed      int64
	Rejected    int64
	CategoryId  int64
	AssignCount int64
	CheckCount  int64
	BatchType   int64
}

// @router /batch [get]
func (c *BatchController) Get() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user.Id)
	var batches = models.GetBatchByUserId(current_user.Id)
	logs.Info(batches)
	batch_view_list := []BatchView{}
	for i := 0; i < len(batches); i++ {
		batch_view := BatchView{
			Id:          batches[i].Id,
			BatchName:   batches[i].BatchName,
			AssignCount: batches[i].AssignCount,
			CheckCount:  batches[i].CheckerAssignCount,
			BatchType:   batches[i].BatchType,
			Created:     batches[i].Created,
		}
		batch_view.Finished = models.GetTaskFinishedCount(current_user.Id, batches[i].Id)
		batch_view.Rest = models.GetTaskRestCount(current_user.Id, batches[i].Id)
		batch_view_list = append(batch_view_list, batch_view)
	}
	c.Data["batches"] = batch_view_list
	c.TplName = "batch.tpl"
}

// @router /batch/status/:page [get]
func (c *BatchController) GetBatchList() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user.Id)
	if current_user.Type != 1 {

		c.Redirect("/login", 302)
	}
	page, _ := strconv.ParseInt(c.Ctx.Input.Param(":page"), 10, 64)
	batch_view_list := []BatchView{}
	pagesize, _ := strconv.ParseInt(beego.AppConfig.String("pagesize"), 10, 64)
	begin := (page - int64(1)) * pagesize
	var batches, batch_count = models.GetBatchForAdminId(begin, pagesize)
	logs.Info(batches)
	for i := 0; i < len(batches); i++ {
		batch_view := BatchView{
			Id:        batches[i].Id,
			BatchName: batches[i].BatchName,
			Created:   batches[i].Created,
		}
		batch_view.Finished = models.GetTaskCount(0, batches[i].Id, 0)
		batch_view.Sum = models.GetTaskCount(0, batches[i].Id, 100)
		batch_view.Rest = models.GetTaskCount(0, batches[i].Id, 0)
		batch_view.Passed = models.GetTaskCount(0, batches[i].Id, 2)
		batch_view.Rejected = models.GetTaskCount(0, batches[i].Id, 3)
		batch_view.CategoryId = batches[i].Category.Id
		batch_view.AssignCount = batches[i].AssignCount
		batch_view.CheckCount = batches[i].CheckerAssignCount
		batch_view.BatchType = batches[i].BatchType
		batch_view_list = append(batch_view_list, batch_view)
	}

	categories := models.GetCategoryList()
	c.Data["categories"] = categories
	c.Data["tagleaders"] = models.GetAllUsers()
	c.Data["batches"] = batch_view_list
	c.Data["Page"] = PageUtil(batch_count, page, pagesize, batch_view_list)
	c.TplName = "admin_batch_list.tpl"
}

// @router /batch/tagleader/review [get]
func (c *BatchController) GetBatchTagLeaderReviewList() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user.Id)
	if current_user.Type != 4 {
		c.Redirect("/login", 302)
	}
	var batches = models.GetBatchByTagleader(current_user.Id)
	logs.Info(batches)
	batch_view_list := []BatchView{}
	for i := 0; i < len(batches); i++ {
		batch_view := BatchView{
			Id:        batches[i].Id,
			BatchName: batches[i].BatchName,
			Created:   batches[i].Created,
		}
		batch_view.Finished = models.GetTaskCount(0, batches[i].Id, 0)
		batch_view.Sum = models.GetTaskCount(0, batches[i].Id, 100)
		batch_view.Rest = models.GetTaskCount(0, batches[i].Id, 0)
		batch_view.Passed = models.GetTaskCount(0, batches[i].Id, 2)
		batch_view.Rejected = models.GetTaskCount(0, batches[i].Id, 3)
		batch_view.CategoryId = batches[i].Category.Id
		batch_view.AssignCount = batches[i].AssignCount
		batch_view.CheckCount = batches[i].CheckerAssignCount
		batch_view.BatchType = batches[i].BatchType
		batch_view_list = append(batch_view_list, batch_view)
	}

	categories := models.GetCategoryList()
	c.Data["categories"] = categories
	c.Data["batches"] = batch_view_list
	c.TplName = "tagleader_batch_list.tpl"
}

// @router /check [get]
func (c *BatchController) GetCheckBatch() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user.Id)
	var batches = models.GetCheckBatch(current_user.Id)
	logs.Info(batches)
	batch_view_list := []BatchView{}
	for i := 0; i < len(batches); i++ {
		batch_view := BatchView{
			Id:        batches[i].Id,
			BatchName: batches[i].BatchName,
			Created:   batches[i].Created,
		}
		batch_view.Finished = models.GetCheckTaskCount(current_user.Id, batches[i].Id, []int64{2, 3})
		batch_view.Rest = models.GetCheckTaskCount(current_user.Id, batches[i].Id, []int64{1})
		batch_view_list = append(batch_view_list, batch_view)
	}
	c.Data["batches"] = batch_view_list
	c.TplName = "check_batch.tpl"
}

// @router /batch/apply [post]
func (c *BatchController) Apply() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	batch_id, _ := strconv.ParseInt(c.Input().Get("batch_id"), 10, 64)
	current_user := current_user_session.(models.User)
	batch := models.GetBatchById(batch_id)
	assign_count := batch.AssignCount
	logs.Info(assign_count)
	is_success := true
	if assign_count != -1 {
		is_success = models.AssignTask(batch_id, current_user.Id, assign_count)
	} else {
		is_success = models.AssignTaskByPerson(batch_id, current_user.Id)
	}
	output := make(map[string]interface{})
	if is_success {
		output["error_code"] = "0"
		output["message"] = "successful"
	} else {
		output["error_code"] = "1"
		output["message"] = "error"
	}
	c.Data["json"] = &output
	logs.Info(c.Data)
	c.ServeJSON()
}

// @router /batch/check/apply [post]
func (c *BatchController) CheckApply() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	batch_id, _ := strconv.ParseInt(c.Input().Get("batch_id"), 10, 64)
	current_user := current_user_session.(models.User)
	batch := models.GetBatchById(batch_id)
	assign_count := batch.CheckerAssignCount
	logs.Info(assign_count)
	is_success := true
	if assign_count != -1 {
		is_success = models.AssignCheckTask(batch_id, current_user.Id, assign_count)
	} else {
		is_success = models.AssignCheckTaskByPerson(batch_id, current_user.Id)
	}
	output := make(map[string]interface{})
	if is_success {
		output["error_code"] = "0"
		output["message"] = "successful"
	} else {
		output["error_code"] = "1"
		output["message"] = "error"
	}
	c.Data["json"] = &output
	logs.Info(c.Data)
	c.ServeJSON()
}

// @router /batch/save [post]
func (c *BatchController) SaveBatch() {
	name := c.Input().Get("name")
	id, _ := strconv.ParseInt(c.Input().Get("id"), 10, 64)
	category_id, _ := strconv.ParseInt(c.Input().Get("category"), 10, 64)
	batch_type, _ := strconv.ParseInt(c.Input().Get("batch_type"), 10, 64)
	assign_count, _ := strconv.ParseInt(c.Input().Get("assign_count"), 10, 64)
	check_count, _ := strconv.ParseInt(c.Input().Get("check_count"), 10, 64)
	tag_leader_id, _ := strconv.ParseInt(c.Input().Get("tag_leader"), 10, 64)
	logs.Info(id)
	logs.Info(name)
	logs.Info(category_id)
	output := make(map[string]interface{})
	flag := int64(0)
	if id != 0 {
		flag = models.SaveBatch(id, name, category_id, assign_count, check_count, batch_type, tag_leader_id)
	} else {
		flag = models.AddBatch(name, category_id, assign_count, check_count, batch_type, tag_leader_id)
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

// @router /batch/delete/user [post]
func (c *BatchController) BatchDeleteUser() {
	user_id, _ := strconv.ParseInt(c.Input().Get("user_id"), 10, 64)
	batch_id, _ := strconv.ParseInt(c.Input().Get("batch_id"), 10, 64)
	output := make(map[string]interface{})
	flag := models.DeleteBatchUserLink(user_id, batch_id)
	if flag {
		output["error_code"] = "0"
	} else {
		output["error_code"] = "1"
		output["message"] = "Failed"
	}

	c.Data["json"] = &output
	logs.Info(c.Data)
	c.ServeJSON()
}

// @router /batch/delete/checker [post]
func (c *BatchController) BatchDeleteChecker() {
	checker_id, _ := strconv.ParseInt(c.Input().Get("checker_id"), 10, 64)
	batch_id, _ := strconv.ParseInt(c.Input().Get("batch_id"), 10, 64)
	output := make(map[string]interface{})
	flag := models.DeleteBatchCheckerLink(checker_id, batch_id)
	if flag {
		output["error_code"] = "0"
	} else {
		output["error_code"] = "1"
		output["message"] = "Failed"
	}

	c.Data["json"] = &output
	logs.Info(c.Data)
	c.ServeJSON()
}

// @router /batch/add/user [post]
func (c *BatchController) BatchAddUser() {
	user_id, _ := strconv.ParseInt(c.Input().Get("user_id"), 10, 64)
	batch_id, _ := strconv.ParseInt(c.Input().Get("batch_id"), 10, 64)
	output := make(map[string]interface{})
	flag := models.AddBatchChecker(user_id, batch_id)
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
