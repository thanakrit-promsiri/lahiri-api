package main

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger สร้างและตั้งค่า zap.Logger พร้อม Daily Log Rotation (หมุนไฟล์ตามวัน)
func InitLogger() *zap.Logger {
	isProd := os.Getenv("APP_ENV") == "production"

	// กำหนดโฟลเดอร์สำหรับเก็บไฟล์ Log (logs/)
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	// รูปแบบชื่อไฟล์: logs/app-%Y-%m-%d.log
	logPattern := filepath.Join(logDir, "app-%Y-%m-%d.log")
	currentSymlink := filepath.Join(logDir, "app.log")

	// ตั้งค่า file-rotatelogs ให้หมุนไฟล์ทุก 24 ชั่วโมง และเก็บย้อนหลัง 30 วัน
	rotator, err := rotatelogs.New(
		logPattern,
		rotatelogs.WithLinkName(currentSymlink),       // สร้าง symlink ไปยังไฟล์ปัจจุบัน (logs/app.log)
		rotatelogs.WithMaxAge(30*24*time.Hour),       // เก็บไฟล์ย้อนหลังสูงสุด 30 วัน
		rotatelogs.WithRotationTime(24*time.Hour),    // ตัดรอบไฟล์ใหม่ทุก 24 ชั่วโมง (เที่ยงคืน)
	)
	if err != nil {
		panic("failed to initialize log rotator: " + err.Error())
	}

	// กำหนด Encoder สำหรับ Console และ File
	var encoderConfig zapcore.EncoderConfig
	if isProd {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // รูปแบบเวลาอ่านง่าย เช่น 2026-11-06T18:35:00.000+0700

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// กำหนด Log Level
	level := zapcore.DebugLevel
	if isProd {
		level = zapcore.InfoLevel
	}

	// รวม Core: เขียนออกทั้ง Terminal (Console) และไฟล์หมุนตามวัน (Rotator)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(rotator), level),
	)

	logger := zap.New(core, zap.AddCaller())
	Log = logger
	return logger
}

// SyncLogger ล้าง Buffer ของ Logger ก่อนปิดโปรแกรม
func SyncLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
