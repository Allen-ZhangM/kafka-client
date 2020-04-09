package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initLogger()

	if err := ConfigInit(); err != nil {
		logs.Error("failed to load kafka conf file  conf/kafka_configure.json: %s", err.Error())
		os.Exit(1)
	}

	if err := InitKafkaConsumers(); err != nil {
		logs.Error("failed to init kafka : %s", err.Error())
		os.Exit(1)
	}
	setSignal()
	select {}
}

func initLogger() {
	logs.SetLevel(7)
	// default 4 for calling with beego.Debug()/Warn()...
	// while 3 is correct for calling with logs.Debug()/Warn()...
	logs.SetLogFuncCallDepth(3)

	logConfig := fmt.Sprintf(`{"filename":"%s", "perm":"0644"}`, "logs/log-srv.log")
	logs.SetLogger("file", logConfig)
	beego.BeeLogger.DelLogger("console")
}

func setSignal() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		sig := <-signals
		logs.Info("got signal : %v ", sig)
		MarkOffsetConsumers(sig.String())
		os.Exit(3)
	}()
}
