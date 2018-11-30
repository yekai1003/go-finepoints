package dbs

import (
	"database/sql"
	"errors"
	"go-finepoint/configs"
	_ "strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wonderivan/logger"
)

//数据库连接的全局变量
var DBConn *sql.DB

type TParams struct {
	TableName   string `json:"table_name"`
	ColumnName  string `json:"column_name"`
	ColumnValue string `json:"column_value"`
	ParamValue  string `json:"param_value"`
}

type RegInfo struct {
	Phone      string    `json:"phone"`
	Code       string    `json:"code"`
	Token      string    `json:"token"`
	ETID       int       `json:"et_id"`
	Name       string    `json:"name"`
	CreditID   string    `json:"credit_id"`
	BussLic    string    `json:"lic"`
	LoginID    string    `json:"login_id"`
	Pass       string    `json:"passwd"`
	Logo       string    `json:"logo"`
	Url        string    `json:"url"`
	Email      string    `json:"email"`
	Info       string    `json:"info"`
	Scale      string    `json:"scale"`
	CreateTime time.Time `json:"create_time"`
}

type Login struct {
	Phone string `json:"phone"`
	Pass  string `json:"passwd"`
	Token string `json:"token"`
	EPID  int    `json:"id"`
}

type MultLogin struct {
	CreditID  string `json:"credit"`
	Pass      string `json:"pass"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Code      string `json:"code"`
	LoginType int    `json:"logintype"` // 0 - phone+pass,1 = phone+code,2 = email + pass ,3 = CreditID+pass
}

type JobInfo struct {
	Token        string `json:"token"`
	JobID        string `json:"job_id"`
	EPID         string `json:"et_id"`
	JobName      string `json:"job_name"`
	SalScope     string `json:"sal_scope"`
	HighEdu      string `json:"high_edu"`
	Experience   string `json:"experience"`
	JobTypeID    string `json:"job_type_id"`
	WorkCity     string `json:"work_city"`
	WorkAddr     string `json:"work_addr"`
	JobProperty  string `json:"job_property"`
	IsCheck      string `json:"is_check"`
	IsHeadHunter string `json:"is_headhunter"`
	JobDetail    string `json:"job_detail"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
}

type JobDeliver struct {
	Token     string `json:"token"`
	EPID      string `json:"et_id"`
	JdID      string `json:"jd_id"`
	JobID     string `json:"job_id"`
	UserID    string `json:"user_id"`
	TargetSal string `json:"target_sal"`
	Similar   string `json:"similar"`
	Comment   string `json:"comment"`
	Status    string `json:"status"`
}

//init函数是本包被其他文件引用时自动执行，并且整个工程只会执行一次
func init() {
	//fmt.Println("call dbs.Init", configs.Config)
	logger.Info("call dbs.Init", configs.Config)
	DBConn = InitDB(configs.Config.Db.Connstr, configs.Config.Db.Driver)
	InitRedis()
}

//初始化数据库连接
func InitDB(connstr, Driver string) *sql.DB {
	db, err := sql.Open(Driver, connstr)
	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}

//func (p *TParams) QueryParams() ([]TParams, int, error) {
//	rows, err := DBConn.Query("select param_type,param_name,param_value from t_params where param_type =?", p.ParamType)
//	if err != nil {
//		//log.Println("failed to query Prams,type=", QueryType)
//		logger.Error("failed to query Prams,type=", p.ParamType, err)
//		return nil, -1, err
//	}
//	var param TParams
//	var paramDatas []TParams
//	count := 0
//	for rows.Next() {
//		err = rows.Scan(&param.ParamType, &param.ParamName, &param.ParamValue)
//		if err != nil {
//			logger.Error("failed to Scan data,type=", p.ParamType, err)
//			return nil, -1, err
//		}
//		paramDatas = append(paramDatas, param)
//		count++
//	}
//	return paramDatas, count, nil
//}

func (r *RegInfo) Insert() (int64, error) {
	res, err := DBConn.Exec("insert into t_enterprise(phone) values(?)", r.Phone)
	if err != nil {
		logger.Error("RegInfo.Insert() err ", r.Phone, err)
		return 0, err
	}
	return res.LastInsertId()
}

func (r *RegInfo) Update() error {
	_, err := DBConn.Exec("update t_enterprise set et_name=?,credit_id=?,business_lic=?,login_id=?,et_email=?,passwd=?,et_logo=?,et_url=?,et_info=?,et_scale=? where phone = ?",
		r.Name, r.CreditID, r.BussLic, r.LoginID, r.Email, r.Pass, r.Logo, r.Url, r.Info, r.Scale, r.Phone)
	if err != nil {
		logger.Error("RegInfo.Update() err ", r.Phone, err)
		return err
	}
	return err
}

//添加logo
func (r *RegInfo) UpdateLogo() error {
	_, err := DBConn.Exec("update t_enterprise set et_logo=? where phone = ?", r.Logo, r.Phone)
	if err != nil {
		logger.Error("RegInfo.UpdateLogo() err ", r.Phone, err)
		return err
	}
	return err
}

//通用查询，返回map嵌套map
func DBQuery(sql string, args ...interface{}) ([]map[string]string, int, error) {
	logger.Info("query is called:", sql)
	rows, err := DBConn.Query(sql, args...)
	if err != nil {
		logger.Error("query data err", err)
		return nil, 0, err
	}
	//得到列名数组
	cols, err := rows.Columns()
	//获取列的个数
	colCount := len(cols)
	values := make([]string, colCount)
	oneRows := make([]interface{}, colCount)
	for k, _ := range values {
		oneRows[k] = &values[k] //将查询结果的返回地址绑定，这样才能变参获取数据
	}
	//存储最终结果
	results := []map[string]string{}
	idx := 0
	//循环处理结果集
	for rows.Next() {
		rows.Scan(oneRows...)
		rowmap := make(map[string]string)
		for k, v := range values {
			rowmap[cols[k]] = v

		}
		results = append(results, rowmap)
		idx++
		//fmt.Println(values)
	}
	//fmt.Println("---------------------------------------")
	logger.Info("query..idx===", idx)
	return results, idx, nil

}

func (m *MultLogin) MultLogin() (bool, error) {
	if m.LoginType == 0 {
		//phone + pass
		row, err := DBConn.Query("select et_name from t_enterprise where phone=? and passwd = ?", m.Phone, m.Pass)
		if err != nil {
			logger.Error("failed to select t_enterprise,", err)
			return false, err
		}
		return row.Next(), err
	} else if m.LoginType == 2 {
		//phone + pass
		row, err := DBConn.Query("select et_name from t_enterprise where et_email=? and passwd = ?", m.Email, m.Pass)
		if err != nil {
			logger.Error("failed to select t_enterprise,", err)
			return false, err
		}
		return row.Next(), err
	} else if m.LoginType == 3 {
		//phone + pass
		row, err := DBConn.Query("select et_name from t_enterprise where credit_id=? and passwd = ?", m.CreditID, m.Pass)
		if err != nil {
			logger.Error("failed to select t_enterprise,", err)
			return false, err
		}
		return row.Next(), err
	}
	return false, errors.New("login type err")
}

func (j *JobInfo) Insert() (int64, error) {
	res, err := DBConn.Exec("insert into t_job_info(et_id,job_name,sal_scope,high_edu,experience,job_type_id,work_city,work_addr,job_property,is_check,is_headhunter,job_detail) values(?,?,?,?,?,?,?,?,?,?,?,?)",
		j.EPID, j.JobName, j.SalScope, j.HighEdu, j.Experience, j.JobTypeID, j.WorkCity, j.WorkAddr, j.JobProperty, j.IsCheck, j.IsHeadHunter, j.JobDetail)
	if err != nil {
		logger.Error("JobInfo.Insert() err ", j.EPID, j.JobName, err)
		return 0, err
	}
	return res.LastInsertId()
}

func (j *JobInfo) UpdateStatus() error {
	_, err := DBConn.Exec("update  t_job_info set status = ? where job_id= ?", j.Status, j.JobID)
	if err != nil {
		logger.Error("UpdateStatus err ", j.JobID, j.Status, err)
		return err
	}
	return err
}
func (j *JobInfo) UpdateAll() error {
	_, err := DBConn.Exec("update  t_job_info set job_name = ?,sal_scope = ?,high_edu = ?,experience = ?,job_type_id =?,work_city = ?,work_addr = ?,is_check = ?, is_headhunter = ?,job_detail = ? where job_id= ?",
		j.JobName, j.SalScope, j.HighEdu, j.Experience, j.JobTypeID, j.WorkCity, j.WorkAddr, j.IsCheck, j.IsHeadHunter, j.JobDetail, j.JobID)
	if err != nil {
		logger.Error("UpdateAll err ", j.JobID, j.Status, err)
		return err
	}
	return err
}

func (j *JobDeliver) UpdateStatus() error {
	_, err := DBConn.Exec("update  t_job_deliver set status = ? where jd_id= ? and job_id=?", j.Status, j.JdID, j.JobID)
	if err != nil {
		logger.Error("UpdateStatus err ", j.JobID, j.Status, err)
		return err
	}
	return err
}
