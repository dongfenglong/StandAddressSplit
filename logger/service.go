package logger

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"log"
	"os"
	"time"
)

var Log *log.Logger

func NewLog() (*log.Logger, error) {
	_, err := os.Stat("./log")
	if err != nil {
		if err := os.Mkdir("./log", os.ModePerm); err != nil {
			return nil, err
		}
	}

	//f, err := os.OpenFile("log/run.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	//if err != nil {
	//	return nil, err
	//}

	logfile := "log/run"
	writer, err := rotatelogs.New(
		logfile+".%Y-%m-%d.log",
		//rotatelogs.WithLinkName(logfile),          //为最新的日志建立软连接
		rotatelogs.WithMaxAge(10*24*time.Hour),    //设置文件清理前的最长保存时间
		rotatelogs.WithRotationTime(24*time.Hour), //设置日志分割的时间，隔多久分割一次
	)

	if err != nil {
		panic(err.Error())
	}

	l := log.New(writer, "", log.Llongfile|log.Ldate|log.Ltime)
	return l, nil
}
