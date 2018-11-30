package dbs

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"go-finepoint/configs"
	_ "strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
)

var rediscli *redis.Client

// redis 设置
type RedisVal struct {
	Key    string `json:"token"`
	Val    string `json:"val"`
	Expire int64  `json:"expire"`
}

func InitRedis() {
	logger.Info("call InitRedis()")
	opt := &redis.Options{
		Addr:     configs.Config.Db.Redisstr,
		Password: "",
		DB:       0,
	}

	rediscli = redis.NewClient(opt)
	pong, err := rediscli.Ping().Result()

	if err != nil {
		logger.Fatal("failed to get redis :", pong, err)
	}
	logger.Info("connect to redis ok :", pong)
}

// 设置数据
func (r *RedisVal) SetData() error {

	err := rediscli.Set(r.Key, r.Val, time.Duration(r.Expire)*time.Second).Err()
	if err != nil {
		logger.Error("failed to set key:", r.Key, r.Val)
	}
	return err
}

//验证数据
func (r *RedisVal) ValidData() (bool, error) {
	val, err := rediscli.Get(r.Key).Result()
	ret := false

	if err != nil {
		logger.Error("failed to get key:", r.Key)
		return ret, err
	}
	if val == "" {
		logger.Error("failed to get key:", r.Key, val)
		return ret, err
	}
	if val != r.Val {
		logger.Error("failed to get key:", r.Key, val)
		return ret, errors.New("校验码不正确")
	}

	ret = true

	return ret, err
}

//验证数据
func (r *RedisVal) CheckKey() (bool, error) {
	val, err := rediscli.Get(r.Key).Result()
	ret := false

	if err != nil {
		logger.Error("failed to get key:", r.Key)
		return ret, err
	}
	if val == "" {
		logger.Error("failed to get key:", r.Key, val)
		return ret, err
	}
	ret = true
	return ret, err
}

func (r *RedisVal) GetData() (string, error) {
	return rediscli.Get(r.Key).Result()

}

func (r *RedisVal) CreateToken(user, pass, code string) string {
	tokenstr := user + pass + code
	fmt.Println(tokenstr)
	hash := sha256.Sum256([]byte(tokenstr))
	r.Key = fmt.Sprintf("%x", hash)
	return r.Key
}
