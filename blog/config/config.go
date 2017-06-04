package config

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/toolkits/file"
	"qiniupkg.com/api.v7/conf"
)

type LogConfig struct {
	Access string `json:"access"`
	Error  string `json:"error"`
}

type HttpConfig struct {
	Listen string `json:"listen"`
	Secret string `json:"secret"`
}

type MysqlConfig struct {
	Addr    string `json:"addr"`
	Idle    int    `json:"idle"`
	Max     int    `json:"max"`
	ShowSql bool   `json:"show_sql"`
}

type TimeoutConfig struct {
	Conn  int64 `json:"conn"`
	Read  int64 `json:"read"`
	Write int64 `json:"write"`
}

type RedisConfig struct {
	Addr    string         `json:"addr"`
	Idle    int            `json:"idle"`
	Max     int            `json:"max"`
	Timeout *TimeoutConfig `json:"timeout"`
}

type CacheConfig struct {
	Provider    string        `json:"provider"`
	ExpireInt64 int64         `json:"expire"`
	Expire      time.Duration `json:"-"`
}

type QiniuConfig struct {
	Enabled   bool   `json:"enabled"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Domain    string `json:"domain"`
}

type ContactConfig struct {
	Weibo   string `json:"weibo"`
	Twitter string `json:"twitter"`
	Github  string `json:"github"`
}

type GlobalConfig struct {
	Debug    bool           `json:"debug"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Portrait string         `json:"portrait"`
	Motto    string         `json:"motto"`
	Log      *LogConfig     `json:"log"`
	Http     *HttpConfig    `json:"http"`
	Mysql    *MysqlConfig   `json:"mysql"`
	Redis    *RedisConfig   `json:"redis"`
	Cache    *CacheConfig   `json:"cache"`
	Upload   string         `json:"upload"`
	Qiniu    *QiniuConfig   `json:"qiniu"`
	Contact  *ContactConfig `json:"contact"`
}

var (
	File string
	G    *GlobalConfig
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	File = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	c.Cache.Expire = time.Duration(c.Cache.ExpireInt64) * time.Second

	if c.Qiniu != nil && c.Qiniu.Enabled {
		conf.ACCESS_KEY = c.Qiniu.AccessKey
		conf.SECRET_KEY = c.Qiniu.SecretKey
	}

	G = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
