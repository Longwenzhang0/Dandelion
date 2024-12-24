package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// new返回的是一个指针，全局都可以使用
var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filename string) (err error) {
	// 设置配置文件的参数
	//viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config")
	//viper.AddConfigPath(".")

	// 设置配置文件
	viper.SetConfigFile(filename)

	// 读取配置文件，错误处理
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err: %v\n", err)
		return
	}

	// 将配置反序列化到结构体中 conf 中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper unmarshal failed, err:%v\n", err)
	}
	// 监控配置文件的变化，变化之后，也要反序列化一遍
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生了修改...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper unmarshal failed, err:%v\n", err)
		}
	})
	return
}
