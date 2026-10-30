package main

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

// InitLogger สร้างและตั้งค่า zap.Logger โดยเลือกโหมดตาม APP_ENV (production หรือ development)
func InitLogger() *zap.Logger {
	var config zap.Config
	if os.Getenv("APP_ENV") == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// กำหนดให้เก็บบันทึกลงทั้งหน้าจอ Console และไฟล์ app.log พร้อมกัน
	config.OutputPaths = []string{"stdout", "app.log"}
	config.ErrorOutputPaths = []string{"stderr", "app.log"}

	logger, err := config.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	Log = logger
	return logger
}

// SyncLogger ล้าง Buffer ของ Logger ก่อนปิดโปรแกรม
func SyncLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
