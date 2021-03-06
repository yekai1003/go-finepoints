package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	Version   = "1.0.0"
	Author    = "yekai"
	BuildTime = "2018-11-27"
)

type ServerConfig struct {
	Common *CommonConfig
	Db     *DbConfig
}

type CommonConfig struct {
	Port     string //启动侦听的端口
	Length   int    // 验证码长度
	FilePath string //文件路径
}

type DbConfig struct {
	Driver    string
	Connstr   string
	Redisstr  string
	RedisPass string
	RedisDB   int
}

func usage() {
	fmt.Printf("Usage: %s -c config_file [-v] [-h]\n", os.Args[0])
}

var Config *ServerConfig //引用配置文件结构

func init() {
	fmt.Println("call config.init")
	Config = GetConfig()
}

func GetConfig() (config *ServerConfig) {
	var configFile = flag.String("c", "", "Config file")

	var ver = flag.Bool("v", false, "version")
	var help = flag.Bool("h", false, "Help")

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		return nil
	}

	if *ver {
		fmt.Println("Version: ", Version)
		fmt.Println("Commit: ", Author)
		fmt.Println("BuildTime: ", BuildTime)
		return nil
	}
	if *configFile == "" {
		*configFile = "etc/finepoints.dev.toml"
	}
	// get server config
	if *configFile != "" {
		config = &ServerConfig{}
		if _, err := toml.DecodeFile(*configFile, &config); err != nil {
			panic(err)
		}
	} else {
		usage()
		return nil
	}

	return config
}
