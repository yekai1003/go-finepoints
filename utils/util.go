package utils

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

const (
	RECODE_OK         = "0"
	RECODE_NODATA     = "4001"
	RECODE_DATAEXISTS = "4002"
	RECODE_NOKEY      = "4003"
	RECODE_REDISERR   = "4004"
	RECODE_MYSQLERR   = "4005"
	RECODE_IOERR      = "4006"
	RECODE_CODEERR    = "4007"
	RECODE_SERVERERR  = "4008"
	RECODE_LOGINERR   = "4009"
	RECODE_PARAMERR   = "4010"
	RECODE_USERERR    = "4011"
	RECODE_HASHERR    = "4012"
	RECODE_PWDERR     = "4013"
	RECODE_TOKENERR   = "4014"
	RECODE_EXISTSERR  = "4015"
	RECODE_IPCERR     = "4016"
	RECODE_THIRDERR   = "4017"
	RECODE_UNKNOWERR  = "4018"
)

var recodeText = map[string]string{
	RECODE_OK:         "成功",
	RECODE_NODATA:     "无数据",
	RECODE_DATAEXISTS: "数据已存在",
	RECODE_NOKEY:      "KEY值未找到",
	RECODE_REDISERR:   "redis错误",
	RECODE_MYSQLERR:   "数据库访问错误",
	RECODE_IOERR:      "文件读写错误",
	RECODE_CODEERR:    "编码错误",
	RECODE_SERVERERR:  "内部错误",
	RECODE_LOGINERR:   "登陆错误",
	RECODE_PARAMERR:   "参数错误",
	RECODE_USERERR:    "用户错误",
	RECODE_HASHERR:    "计算hash错误",
	RECODE_PWDERR:     "密码错误",
	RECODE_TOKENERR:   "令牌错误",
	RECODE_EXISTSERR:  "重复上传错误",
	RECODE_IPCERR:     "IPC错误",
	RECODE_THIRDERR:   "与以太坊交互失败",
	RECODE_UNKNOWERR:  "未知错误",
}

const (
	PageMax = 5
)

func RecodeText(code string) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[RECODE_UNKNOWERR]
}

type Resp struct {
	Errno  string      `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

//resp数据响应
func ResponseData(c echo.Context, resp *Resp) {
	resp.ErrMsg = RecodeText(resp.Errno)
	c.JSON(http.StatusOK, resp)
}

//读取dir目录下文件名带address的文件
func GetFileName(address, dirname string) (string, error) {

	data, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println("read dir err", err)
		return "", err
	}
	for _, v := range data {
		if strings.Index(v.Name(), address) > 0 {
			//代表找到文件
			return v.Name(), nil
		}
	}

	return "", nil
}

//获取hash前8位得到对应的数字
func GetFileHash(data []byte) (uint64, error) {

	yy := sha256.Sum256(data)
	return strconv.ParseUint(fmt.Sprintf("%x", yy[:8]), 16, 64)
}

//获取 yyyy/mm/dd 格式字符串
func GetYearMonthDay() string {
	//time.Now().Year()
	return fmt.Sprintf("%04d/%02d/%02d/", time.Now().Year(), time.Now().Month(), time.Now().Day())

}

//创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
