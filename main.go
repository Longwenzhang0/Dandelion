package main

import (
	"Dandelion/controller"
	"Dandelion/dao/mysql"
	"Dandelion/dao/redis"
	"Dandelion/logger"
	"Dandelion/pkg/snowflake"
	"Dandelion/routes"
	"Dandelion/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Go Web开发通用脚手架

func main() {
	fmt.Println(os.Args)
	// 0. 获取命令行参数。第0个元素是可执行文件本身，传入的应该是后续的元素，这里规定第二个参数为配置文件
	if len(os.Args) < 2 {
		fmt.Println("need configuration file; eg: dandelion config.yaml")
		return
	}

	// 1. 加载配置文件
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("Init logger failed, err:%v\n", err)
		return
	}
	// 把缓冲区的日志追加到文件
	defer zap.L().Sync()

	zap.L().Debug("logger init succeed!")
	// 3. 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("Init MySQL failed, err:%v\n", err)
		return
	}
	// 初始化完成之后关闭mysql
	defer mysql.Close()

	// 4. 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("Init Redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	fmt.Println(settings.Conf.Name)
	fmt.Println(settings.Conf.StartTime)
	fmt.Println(settings.Conf.RedisConfig)
	fmt.Println(settings.Conf.MySQLConfig)
	fmt.Println(settings.Conf.LogConfig)

	// 5. 初始化snowflake算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("Init Snowflake failed, err:%v\n", err)
		return
	}
	// 6. 初始化校验器validator的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("Init validator translator failed, err:%v\n", err)
		return
	}
	// 7. 注册路由
	r := routes.SetupRouter(settings.Conf.Mode)

	// 8. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}
