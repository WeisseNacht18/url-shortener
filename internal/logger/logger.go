package logger

import (
	"go.uber.org/zap"
)

var (
	Logger zap.SugaredLogger
)

func Init() {
	temp, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer temp.Sync()

	Logger = *temp.Sugar()
}
