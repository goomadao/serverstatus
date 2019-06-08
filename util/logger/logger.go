package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LogFile, LogLevel              string
	debugCore, infoCore, errorCore zapcore.Core
	Logger                         *zap.Logger
)

func InitLogger() {
	if LogLevel == "debug" {
		debugFile, err := os.Create(LogFile + ".debug")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Create log file failed, create on /tmp/statusserver.log.debug")
			var err2 error
			debugFile, err2 = os.Create("/tmp/serverstatus.log.debug")
			if err2 != nil {
				fmt.Println(err)
				fmt.Println("Create log file on /tmp failed")
				os.Exit(1)
			}
		}
		errorFile, err := os.Create(LogFile + ".error")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Create log file failed, create on /tmp/statusserver.log.error")
			var err2 error
			errorFile, err2 = os.Create("/tmp/serverstatus.log.error")
			if err2 != nil {
				fmt.Println(err)
				fmt.Println("Create log file on /tmp failed")
				os.Exit(1)
			}
		}
		debugCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(debugFile)),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.DebugLevel && lvl <= zapcore.WarnLevel
			}),
		)
		errorCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(errorFile)),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)
		core := zapcore.NewTee(debugCore, errorCore)
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel), zap.Development())
	} else if LogLevel == "info" {
		infoFile, err := os.Create(LogFile + ".info")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Create log file failed, create on /tmp/statusserver.log.info")
			var err2 error
			infoFile, err2 = os.Create("/tmp/serverstatus.log.info")
			if err2 != nil {
				fmt.Println(err)
				fmt.Println("Create log file on /tmp failed")
				os.Exit(1)
			}
		}
		errorFile, err := os.Create(LogFile + ".error")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Create log file failed, create on /tmp/statusserver.log.error")
			var err2 error
			errorFile, err2 = os.Create("/tmp/serverstatus.log.error")
			if err2 != nil {
				fmt.Println(err)
				fmt.Println("Create log file on /tmp failed")
				os.Exit(1)
			}
		}
		infoCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(infoFile)),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl == zapcore.InfoLevel || lvl == zapcore.WarnLevel
			}),
		)
		errorCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(errorFile)),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)
		core := zapcore.NewTee(infoCore, errorCore)
		Logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		errorFile, err := os.Create(LogFile + ".error")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Create log file failed, create on /tmp/statusserver.log.error")
			var err2 error
			errorFile, err2 = os.Create("/tmp/serverstatus.log.error")
			if err2 != nil {
				fmt.Println(err)
				fmt.Println("Create log file on /tmp failed")
				os.Exit(1)
			}
		}
		errorCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(errorFile)),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)
		core := zapcore.NewTee(errorCore)
		Logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	}
}
