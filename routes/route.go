package routes

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"go-finepoint/configs"
	"go-finepoint/dbs"
	"go-finepoint/utils"

	_ "github.com/gorilla/sessions"
	"github.com/labstack/echo"
	_ "github.com/labstack/echo-contrib/session"
	"github.com/wonderivan/logger"
)

func PingHandler(c echo.Context) error {

	return c.String(http.StatusOK, "pong")
}

//获取参数
func GetParams(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型

	paramType := c.Param("paramtype")
	if paramType == "" {
		logger.Error("param get err", paramType)
		resp.Errno = utils.RECODE_PARAMERR
		return nil
	}
	//查询数据库的参数
	var sql string
	if paramType == "edu" {
		//dbs.DBQuery("select * from t_params where column_name='high_edu'")
		sql = fmt.Sprintf("select * from t_params where column_name='high_edu'")
	} else if paramType == "scale" {
		sql = fmt.Sprintf("select * from t_params where column_name='et_scale'")
	} else if paramType == "age" {
		sql = fmt.Sprintf("select * from t_params where column_name='work_age'")
	}
	data, num, err := dbs.DBQuery(sql)
	if err != nil || num <= 0 {
		logger.Error("param DBQuery err:", err)
		resp.Errno = utils.RECODE_MYSQLERR
		return nil
	}
	resp.Data = data

	return nil
}

//注册仅仅是发送手机验证码
func EnterpriseRegister(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型
	var reg dbs.RegInfo
	if err := c.Bind(&reg); err != nil {
		logger.Error("Register bind err:", reg)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	var redisval dbs.RedisVal
	redisval.Key = reg.Phone
	redisval.Expire = 300
	//生成随机数字
	logger.Info(redisval, reg)
	if reg.Code == "" {
		//当前为空，需要生成数字
		reg.Code = utils.GetRandomString(configs.Config.Common.Length)
		redisval.Val = reg.Code
		//存入redis

		err := redisval.SetData()
		if err != nil {
			logger.Error("redis set data err", reg)
			resp.Errno = utils.RECODE_REDISERR
			return err
		}
		//发送短信 ????
		logger.Info("get short msg code ==", reg.Code)
		resp.Errno = utils.RECODE_NODATA
		resp.Data = redisval
		return nil

	}
	//验证code和手机是否ok
	redisval.Val = reg.Code
	ok, err := redisval.ValidData()
	if err != nil {
		logger.Error("redis Valid data err", err)
		resp.Errno = utils.RECODE_REDISERR
		return err
	}
	if !ok {
		logger.Error("redis Valid data err", err)
		resp.Errno = utils.RECODE_CODEERR
		return err
	}

	//插入到mysql--代表注册成功
	_, err = reg.Insert()
	if err != nil {
		logger.Error("reg insert err", err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}

	//生成token
	var token dbs.RedisVal
	token.CreateToken(reg.Phone, reg.Pass, reg.Code)
	token.Expire = 60 * 60
	token.Val = reg.Phone
	if err = token.SetData(); err != nil {
		logger.Error("token.SetData() err", err)
		resp.Errno = utils.RECODE_REDISERR
		return err
	}

	resp.Data = token

	return nil
}

//企业注册信息完善
func EnterpriseCommit(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型
	var reg dbs.RegInfo
	if err := c.Bind(&reg); err != nil {
		logger.Error("Register bind err:", err, reg)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	logger.Info("EnterpriseCommit:", reg)

	//更新企业信息
	err := reg.Update()
	if err != nil {
		logger.Error("reg update err", err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}

	return nil
}

//检测是否登陆
func IsLogin(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型
	var user dbs.Login
	if err := c.Bind(&user); err != nil {
		logger.Error("IsLogin bind err:", user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	//检测token在redis是否存在，如果不存在，则未登陆
	if user.Token == "" {
		logger.Error("IsLogin token is nil:", user)
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("token is nil")
	}

	var token dbs.RedisVal
	token.Key = user.Token

	ok, err := token.CheckKey()
	if err != nil {
		logger.Error("redis Valid data err", err)
		resp.Errno = utils.RECODE_REDISERR
		return err
	}
	if !ok {
		logger.Error("redis Valid data err", err)
		resp.Errno = utils.RECODE_TOKENERR
		return err
	}

	return nil
}

//企业版登陆
func EnterPriseLogin(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型
	var userlogin dbs.MultLogin
	if err := c.Bind(&userlogin); err != nil {
		logger.Error("EnterPriseLogin bind err:", userlogin)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	ok, err := userlogin.MultLogin()
	if err != nil {
		logger.Error("EnterPriseLogin MultLogin err:", userlogin)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	if !ok {
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	var token dbs.RedisVal
	token.CreateToken(userlogin.Phone, userlogin.Pass, userlogin.Code)
	token.Expire = 24 * 60 * 60 * 30
	token.Val = userlogin.Phone
	if err = token.SetData(); err != nil {
		logger.Error("token.SetData() err", err)
		resp.Errno = utils.RECODE_REDISERR
		return err
	}

	resp.Data = token

	return nil
}

//获取企业信息
func EnterPriseInfo(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型

	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	if phone == "" {
		logger.Error("login expire:", token)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	//通过手机号查询企业信息
	data, _, err := dbs.DBQuery("select * from t_enterprise where phone = ?", phone)
	if err != nil {
		logger.Error("get enterprise info err:", err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}

	resp.Data = data

	return nil
}

//上传图片
func UploadLogo(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//2. 获得token
	var user dbs.Login
	if err := c.Bind(&user); err != nil {
		logger.Error("IsLogin bind err:", user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = user.Token
	var reg dbs.RegInfo
	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	if phone == "" {
		logger.Error("login expire:", token)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	reg.Phone = phone

	//3. 解析数据

	h, err := c.FormFile("logo")
	if err != nil {
		fmt.Println("failed to FormFile ", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	src, err := h.Open()
	defer src.Close()

	//4. 计算hash
	cData := make([]byte, h.Size)
	n, err := src.Read(cData)
	if err != nil || h.Size != int64(n) {
		resp.Errno = utils.RECODE_IOERR
		return err
	}
	//这里需要hash打散
	hash, err := utils.GetFileHash(cData)
	if err != nil {
		logger.Error("failed to getfileHash", err)
		resp.Errno = utils.RECODE_HASHERR
		return err
	}

	sonPath := fmt.Sprintf("%02d", hash%100)
	allPath := configs.Config.Common.FilePath + utils.GetYearMonthDay() + sonPath
	if err := utils.CreateDir(allPath); err != nil {
		logger.Error("failed to createdir", err)
		resp.Errno = utils.RECODE_IOERR
		return err
	}

	//写文件
	reg.Logo = allPath + "/" + h.Filename
	dst, err := os.Create(reg.Logo)
	if err != nil {
		fmt.Println("failed to create file ", err, reg.Logo)
		resp.Errno = utils.RECODE_IOERR
		return err
	}
	defer dst.Close()

	dst.Write(cData)

	//5. 操作mysql-新增数据
	err = reg.UpdateLogo()
	if err != nil {
		fmt.Println("failed to update logo ", err, reg.Logo)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}

	return nil
}

//查看单个logo
func GetLogo(c echo.Context) error {
	//获得token
	var user dbs.Login
	if err := c.Bind(&user); err != nil {
		logger.Error("IsLogin bind err:", user)
		//resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = user.Token
	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		//resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	if phone == "" {
		logger.Error("loin expire:", token)
		//resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	//通过数据库获得文件路径
	m, n, err := dbs.DBQuery("select et_logo from t_enterprise where phone = ?", phone)
	if err != nil || n <= 0 {
		logger.Error("DBQuery:", token)
		//resp.Errno = utils.RECODE_MYSQLERR
		return err
	}
	allPath := m[0]["et_logo"]
	//最好查数据库获得文件路径
	http.ServeFile(c.Response(), c.Request(), allPath)
	return nil
}

//发布职位
func PublishJob(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//2. 获得token
	var info dbs.JobInfo
	if err := c.Bind(&info); err != nil {
		logger.Error("PublishJob bind err:", info)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = info.Token

	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	if phone == "" {
		logger.Error("login expire:", token)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	m, n, err := dbs.DBQuery("select et_id from t_enterprise where phone = ?", phone)
	if err != nil || n <= 0 {
		logger.Error("DBQuery:", phone, err)
		resp.Errno = utils.RECODE_NODATA
		return err
	}
	epid := m[0]["et_id"]
	//发布一个职位
	info.EPID = epid
	_, err = info.Insert()
	if err != nil || n <= 0 {
		logger.Error("Insert:", info, err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}
	return nil
}

//查看企业对应的职位
func GetJobs(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 获得token
	var info dbs.JobInfo
	if err := c.Bind(&info); err != nil {
		logger.Error("PublishJob bind err:", info)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = info.Token

	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	if phone == "" {
		logger.Error("login expire:", token)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	//3. 查询发布的职位信息
	//select a.* from t_job_info a,t_enterprise b where a.et_id = b.et_id;
	data, num, err := dbs.DBQuery("select a.* from t_job_info a,t_enterprise b where a.et_id = b.et_id and b.phone = ?", phone)
	if err != nil {
		logger.Error("DBQuery jobinfo:", info, err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}
	if num <= 0 {
		logger.Error("DBQuery jobinfo:", info, err)
		resp.Errno = utils.RECODE_NODATA
		return err
	}
	resp.Data = data
	return nil
}

//查看job类型
func GetJobTypes(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//查询请求参数类型
	data, num, err := dbs.DBQuery("select * from t_job_type ")
	if err != nil || num <= 0 {
		logger.Error("DBQuery job :", err)
		resp.Errno = utils.RECODE_NODATA
		return err
	}
	resp.Data = data

	return nil
}

//查看某一职位信息
func GetJobInfo(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 获得token

	var info dbs.JobInfo
	if err := c.Bind(&info); err != nil {
		logger.Error("PublishJob bind err:", info)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = info.Token

	_, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	jobid := c.Param("id")

	//3. 查询发布的职位信息
	//select a.* from t_job_info a,t_enterprise b where a.et_id = b.et_id;
	data, num, err := dbs.DBQuery("select a.* from t_job_info a where a.job_id = ?", jobid)
	if err != nil {
		logger.Error("DBQuery jobinfo:", info, err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}
	if num <= 0 {
		logger.Error("DBQuery jobinfo:", info, err)
		resp.Errno = utils.RECODE_NODATA
		return err
	}
	resp.Data = data
	return nil
}

//修改职位信息
func UpdateJobInfo(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 获得token

	var info dbs.JobInfo
	if err := c.Bind(&info); err != nil {
		logger.Error("PublishJob bind err:", info)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = info.Token

	_, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	err = info.UpdateAll()
	if err != nil {
		logger.Error("UpdateAll err:", err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}
	return nil
}

//查看简历投递情况
func GetJobDelivers(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 获得token

	var info dbs.JobInfo
	if err := c.Bind(&info); err != nil {
		logger.Error("PublishJob bind err:", info)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = info.Token

	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	if phone == "" {
		logger.Error("login expire:", token)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	m, n, err := dbs.DBQuery("select et_id from t_enterprise where phone = ?", phone)
	if err != nil || n <= 0 {
		logger.Error("DBQuery:", phone, err)
		resp.Errno = utils.RECODE_NODATA
		return err
	}
	epid := m[0]["et_id"]

	//select a.jd_id,d.job_name,d.job_id,b.user_name,c.age,a.target_sal,a.similar,c.tags,a.status   from t_job_deliver a,t_user b,t_user_detail c,t_job_info d  where a.job_id = d.job_id    and a.user_id = b.user_id    and a.user_id = c.user_id    and d.et_id = 1
	m, num, err := dbs.DBQuery("select a.jd_id,d.job_name,d.job_id,b.user_name,c.age,a.target_sal,a.similar,c.tags,a.status   from t_job_deliver a,t_user b,t_user_detail c,t_job_info d  where a.job_id = d.job_id    and a.user_id = b.user_id    and a.user_id = c.user_id    and d.et_id = ?", epid)
	if err != nil {
		logger.Error("DBQuery:", phone, err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}
	if n <= 0 {
		logger.Error("DBQuery nodata:", epid, err)
		resp.Errno = utils.RECODE_NODATA
		return err
	}
	total_page := int(num)/utils.PageMax + 1
	current_page := 1
	mapResp := make(map[string]interface{})
	mapResp["total_page"] = total_page
	mapResp["current_page"] = current_page
	mapResp["delivers"] = m
	resp.Data = mapResp
	return nil
}

//修改简历投递情况
//需要delier_id+token+status
func UpdateJobDeliver(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 获得token

	var info dbs.JobDeliver
	if err := c.Bind(&info); err != nil {
		logger.Error("PublishJob bind err:", info)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//检验是否登陆，并且拿到手机号
	var token dbs.RedisVal
	token.Key = info.Token

	phone, err := token.GetData()
	if err != nil {
		logger.Error("do not login:", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	if phone == "" {
		logger.Error("login expire:", token)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	//更新数据
	err = info.UpdateStatus()
	if err != nil {
		logger.Error("UpdateStatus deliver err:", err)
		resp.Errno = utils.RECODE_MYSQLERR
		return err
	}

	return nil
}
