package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"label/models"
	"strconv"
	"time"
)

type TaskController struct {
	beego.Controller
}

type taskView struct {
	Id         int64
	ImagePath  string
	Updated    string
	Result     string
	Comment    string
	TaskStatus string
}

type AdminTaskView struct {
	UserName string
	Sum      int64
	Finished int64
	Rest     int64
	Passed   int64
	Rejected int64
}

// @router /batch/:batch_id/tag [get]
func (c *TaskController) Get() {
	logs.Info("get batch id " + c.Ctx.Input.Param(":batch_id"))
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	batch_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":batch_id"), 10, 64)

	rest := models.GetTaskRestCount(current_user.Id, batch_id)
	if rest == 0 {
		beego.Informational("no task")
		c.Redirect("/batch", 302)
		return
	}

	var task = models.GetOneTaskByUserIdAndBatchId(current_user.Id, batch_id)
	logs.Info(task)
	c.Data["task"] = task
	c.TplName = "task.tpl"
}

// @router /batch/:batch_id/status [get]
func (c *TaskController) GetTaskStatusForEveryTagger() {
	logs.Info("get batch id " + c.Ctx.Input.Param(":batch_id"))
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	if current_user.Type != 1 {
		c.Redirect("/login", 302)
	}
	batch_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":batch_id"), 10, 64)

	users := models.GetUsersByBatchId(batch_id)
	item_view_list := []AdminTaskView{}
	for i := 0; i < len(users); i++ {
		item_view := AdminTaskView{
			UserName: users[i].Name,
		}
		tag_count := models.GetTaskCount(users[i].Id, batch_id, 1)
		item_view.Sum = models.GetTaskCount(users[i].Id, batch_id, 100)
		item_view.Rest = models.GetTaskCount(users[i].Id, batch_id, 0)
		item_view.Passed = models.GetTaskCount(users[i].Id, batch_id, 2)
		item_view.Rejected = models.GetTaskCount(users[i].Id, batch_id, 3)
		item_view.Finished = tag_count + item_view.Rejected + item_view.Passed
		item_view_list = append(item_view_list, item_view)
	}
	c.Data["users"] = item_view_list
	c.TplName = "admin_user_task_list.tpl"
}

// @router /batch/:batch_id/check/:page [get]
func (c *TaskController) Get_Check_Task_List() {
	logs.Info("get batch id " + c.Ctx.Input.Param(":batch_id"))
	current_user_session := c.GetSession("label_me_user")
	page, _ := strconv.ParseInt(c.Ctx.Input.Param(":page"), 10, 64)
	pagesize, _ := strconv.ParseInt(beego.AppConfig.String("pagesize"), 10, 64)
	begin := (page - int64(1)) * pagesize
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	batch_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":batch_id"), 10, 64)

	taskes, task_count := models.GetTasksForChecker(batch_id, current_user.Id, begin, pagesize)
	logs.Info(taskes)
	task_view_list := []taskView{}
	for i := 0; i < len(taskes); i++ {
		task_view := taskView{
			Id:        taskes[i].Id,
			ImagePath: taskes[i].Image.Path,
			Updated:   taskes[i].Updated,
			Comment:   taskes[i].Comment,
		}
		switch taskes[i].TaskStatus {
		case 0:
			task_view.TaskStatus = "Created"
		case 1:
			task_view.TaskStatus = "Tagged"
		case 2:
			task_view.TaskStatus = "Pass"
		case 3:
			task_view.TaskStatus = "Reject"
		}
		result_json := models.Label_Result_Json{}
		if err := json.Unmarshal([]byte(taskes[i].ResultJson), &result_json); err != nil {
			logs.Info(err)
		}
		options := result_json.Classification.Options
		for j := 0; j < len(options); j++ {
			task_view.Result += options[j] + ","
		}
		task_view_list = append(task_view_list, task_view)
	}
	c.Data["taskes"] = task_view_list
	c.Data["batch_id"] = batch_id
	c.Data["Page"] = PageUtil(task_count, page, pagesize, taskes)
	c.TplName = "task_check_list.tpl"
}

// @router /batch/:batch_id/tagleader/review/:page [get]
func (c *TaskController) TagleaderReviewTask() {
	logs.Info("get batch id " + c.Ctx.Input.Param(":batch_id"))
	current_user_session := c.GetSession("label_me_user")
	page, _ := strconv.ParseInt(c.Ctx.Input.Param(":page"), 10, 64)
	pagesize, _ := strconv.ParseInt(beego.AppConfig.String("pagesize"), 10, 64)
	begin := (page - int64(1)) * pagesize
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	//current_user := current_user_session.(models.User)
	batch_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":batch_id"), 10, 64)

	taskes, task_count := models.GetTasksForTagleader(batch_id, begin, pagesize)
	logs.Info(taskes)
	c.Data["taskes"] = taskes
	c.Data["batch_id"] = batch_id
	c.Data["Page"] = PageUtil(task_count, page, pagesize, taskes)
	c.TplName = "task_tagleader_review_list.tpl"
}

// @router /task/post [post]
func (c *TaskController) PostTask() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	task_id, _ := strconv.ParseInt(c.Input().Get("task_id"), 10, 64)
	comment := c.Input().Get("comment")
	region_data := c.Input().Get("region_data")
	var options []string
	c.Ctx.Input.Bind(&options, "options")
	var num = models.TagTask(task_id, options, comment, region_data)
	logs.Info(num)
	output := make(map[string]interface{})
	if num == 1 {
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

// @router /task/pass [post]
func (c *TaskController) PassTask() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	task_id, _ := strconv.ParseInt(c.Input().Get("task_id"), 10, 64)
	comment := c.Input().Get("comment")
	var num = models.CheckTask(task_id, comment, true)
	logs.Info(num)
	output := make(map[string]interface{})
	if num == 1 {
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

// @router /task/reject [post]
func (c *TaskController) RejectTask() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	task_id, _ := strconv.ParseInt(c.Input().Get("task_id"), 10, 64)
	comment := c.Input().Get("comment")
	var num = models.CheckTask(task_id, comment, false)
	logs.Info(num)
	output := make(map[string]interface{})
	if num == 1 {
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

// @router /task/multicheck [post]
func (c *TaskController) MultiTasksCheck() {
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user)
	comment := c.Input().Get("comment")
	logs.Info(comment)
	var pass_tasks []int64
	c.Ctx.Input.Bind(&pass_tasks, "pass")
	logs.Info(pass_tasks)
	var reject_tasks []int64
	c.Ctx.Input.Bind(&reject_tasks, "reject")
	logs.Info(reject_tasks)
	flag := models.CheckTasks(pass_tasks, reject_tasks, comment)
	output := make(map[string]interface{})
	if flag == true {
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

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// @router /task/list/:time_span/:page [get]
func (c *TaskController) GetList() {
	page, _ := strconv.ParseInt(c.Ctx.Input.Param(":page"), 10, 64)
	pagesize, _ := strconv.ParseInt(beego.AppConfig.String("pagesize"), 10, 64)
	begin := (page - int64(1)) * pagesize
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	time_span := c.Ctx.Input.Param(":time_span")
	starttime := ""
	endtime := ""
	if time_span == "today" {
		starttime = time.Now().Format("2006-01-02 00:00:00")
		endtime = time.Now().Format("2006-01-02 15:04:05")
	}

	taskes, task_count := models.GetTaskList(current_user.Id, starttime, endtime, begin, pagesize)
	logs.Info(taskes)
	c.Data["taskes"] = taskes
	c.Data["time_span"] = time_span
	c.Data["Page"] = PageUtil(task_count, page, pagesize, taskes)
	c.TplName = "task_list.tpl"
}

// @router /check/review/today/:page [get]
func (c *BatchController) GetCheckToday() {
	current_user_session := c.GetSession("label_me_user")
	page, _ := strconv.ParseInt(c.Ctx.Input.Param(":page"), 10, 64)
	pagesize, _ := strconv.ParseInt(beego.AppConfig.String("pagesize"), 10, 64)
	begin := (page - int64(1)) * pagesize
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	current_user := current_user_session.(models.User)
	logs.Info(current_user.Id)
	starttime := time.Now().Format("2006-01-02 00:00:00")
	endtime := time.Now().Format("2006-01-02 15:04:05")

	taskes, task_count := models.GetCheckTaskList(current_user.Id, starttime, endtime, begin, pagesize)
	logs.Info(taskes)
	task_view_list := []taskView{}
	for i := 0; i < len(taskes); i++ {
		task_view := taskView{
			Id:        taskes[i].Id,
			ImagePath: taskes[i].Image.Path,
			Updated:   taskes[i].Updated,
			Comment:   taskes[i].Comment,
		}
		switch taskes[i].TaskStatus {
		case 0:
			task_view.TaskStatus = "Created"
		case 1:
			task_view.TaskStatus = "Tagged"
		case 2:
			task_view.TaskStatus = "Pass"
		case 3:
			task_view.TaskStatus = "Reject"
		}
		result_json := models.Label_Result_Json{}
		if err := json.Unmarshal([]byte(taskes[i].ResultJson), &result_json); err != nil {
			logs.Info(err)
		}
		options := result_json.Classification.Options
		for j := 0; j < len(options); j++ {
			task_view.Result += options[j] + ","
		}
		task_view_list = append(task_view_list, task_view)
	}
	c.Data["taskes"] = task_view_list
	c.Data["Page"] = PageUtil(task_count, page, pagesize, taskes)
	c.TplName = "task_check_review_list.tpl"
}

// @router /task/update/:task_id [get]
func (c *TaskController) ShowTaskResult() {
	logs.Info("get task id " + c.Ctx.Input.Param(":task_id"))
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	//current_user := current_user_session.(models.User)
	task_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":task_id"), 10, 64)

	result := models.GetTaskResult(task_id)
	c.Data["task"] = result
	c.TplName = "task_detail.tpl"
}

// @router /task/check/:task_id [get]
func (c *TaskController) ShowCheckTaskResult() {
	logs.Info("get task id " + c.Ctx.Input.Param(":task_id"))
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	//current_user := current_user_session.(models.User)
	task_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":task_id"), 10, 64)

	task := models.GetTaskResult(task_id)
	if task.Batch.BatchType == 1 {
		c.Data["task"] = task
		c.TplName = "task_check_detail.tpl"
	} else if task.Batch.BatchType == 2 {
		tasks := models.GetTasksResultWithSameImage(task.Batch.Id, task.Image.Id)
		c.Data["extend_field"] = task.Image.ExtendField
		c.Data["image_path"] = task.Image.Path
		c.Data["batch_id"] = task.Batch.Id
		c.Data["tasks"] = tasks
		c.TplName = "task_check_detail2.tpl"
	} else {
		c.Redirect("/login", 302)
		return
	}
}

// @router /task/check/update/:task_id [get]
func (c *TaskController) ShowCheckUpdateTaskResult() {
	logs.Info("get task id " + c.Ctx.Input.Param(":task_id"))
	current_user_session := c.GetSession("label_me_user")
	if current_user_session == nil {
		logs.Info("session lost")
		beego.Informational("session lost")
		c.Redirect("/login", 302)
		return
	}
	//current_user := current_user_session.(models.User)
	task_id, _ := strconv.ParseInt(c.Ctx.Input.Param(":task_id"), 10, 64)

	task := models.GetTaskResult(task_id)
	if task.Batch.BatchType == 1 {
		c.Data["task"] = task
		c.TplName = "task_check_update_detail.tpl"
	} else if task.Batch.BatchType == 2 {
		tasks := models.GetTasksResultWithSameImage(task.Batch.Id, task.Image.Id)
		c.Data["extend_field"] = task.Image.ExtendField
		c.Data["image_path"] = task.Image.Path
		c.Data["batch_id"] = task.Batch.Id
		c.Data["tasks"] = tasks
		c.TplName = "task_check_detail2.tpl"
	} else {
		c.Redirect("/login", 302)
		return
	}
}
