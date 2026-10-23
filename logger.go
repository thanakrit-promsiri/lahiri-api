package main

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	// กำหนดให้เก็บบันทึกลงทั้งหน้าจอ Console และไฟล์ app.log พร้อมกัน
	config := zap.NewDevelopmentConfig()

	// เพิ่มไฟล์ที่จะให้บันทึก (เช่น app.log)
	config.OutputPaths = []string{"stdout", "app.log"}
	config.ErrorOutputPaths = []string{"stderr", "app.log"}

	var err error
	Log, err = config.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
}

func SyncLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
