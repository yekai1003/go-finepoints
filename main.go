package main

import (
	"go-finepoint/configs"
	"go-finepoint/routes"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/wonderivan/logger"
)

var EchoObj *echo.Echo //echo框架对象全局定义

func main() {
	logger.Info("get config %v ,%v", configs.Config.Common, configs.Config.Db)
	//fmt.Printf("get config %v ,%v\n", configs.Config.Common.Port, configs.Config.Db.Connstr)
	EchoObj = echo.New()             //创建echo对象
	EchoObj.Use(middleware.Logger()) //安装日志中间件
	EchoObj.Use(middleware.Recover())
	EchoObj.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	EchoObj.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	//设置路由
	EchoObj.GET("/ping", routes.PingHandler)                        //路由测试函数
	EchoObj.GET("/params/:paramtype", routes.GetParams)             //获取参数
	EchoObj.POST("/enterprise/register", routes.EnterpriseRegister) //企业注册
	EchoObj.POST("/enterprise/commit", routes.EnterpriseCommit)     //企业注册信息完善
	EchoObj.POST("/enterprise/login", routes.EnterPriseLogin)       //企业登陆
	EchoObj.GET("/enterprise/info", routes.EnterPriseInfo)          //获取企业信息
	EchoObj.GET("/islogin", routes.IsLogin)                         //检测是否登陆
	EchoObj.POST("/enterprise/logo", routes.UploadLogo)             //上传logo
	EchoObj.GET("/enterprise/logo", routes.GetLogo)                 //下载logo
	EchoObj.GET("/enterprise/jobtypes", routes.GetJobTypes)         //获取job类型
	EchoObj.GET("/enterprise/jobinfo", routes.PublishJob)           //发布职位
	EchoObj.GET("/enterprise/jobinfo/:id", routes.GetJobInfo)       //获取职位信息
	EchoObj.GET("/enterprise/jobinfos", routes.GetJobs)             //获取企业职位
	EchoObj.GET("/enterprise/udpjobinfo", routes.UpdateJobInfo)     //修改已发布职位信息
	EchoObj.GET("/enterprise/delivers", routes.GetJobDelivers)      //查看投递情况
	EchoObj.GET("/enterprise/deliver", routes.UpdateJobDeliver)     //修改投递情况
	EchoObj.Logger.Fatal(EchoObj.Start(configs.Config.Common.Port)) //启动服务
}
