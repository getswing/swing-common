package sw

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/google/uuid"
	"gopkg.in/natefinch/lumberjack.v2"

	gormLogger "gorm.io/gorm/logger"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogEntry struct {
	Timestamp string   `json:"timestamp"`
	Service   string   `json:"service"`
	Level     LogLevel `json:"level"`
	Message   string   `json:"message"`
	Function  string   `json:"function,omitempty"`
	Duration  string   `json:"duration,omitempty"`
	File      string   `json:"file,omitempty"`
	RequestID string   `json:"request_id,omitempty"`
}

type contextKey string

const requestIDKey = contextKey("request_id")

var (
	GlobalLogger *log.Logger
	ServiceName  string
	Skip         int = 1
)

func InitLumberJack() io.Writer {
	logFile := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./logs/%s.log", ServiceName),
		MaxSize:    10,   // max is 10 mb
		MaxBackups: 5,    // max file rotate
		MaxAge:     30,   // 30 days
		Compress:   true, // (.gz) format for rotate
	}

	return io.MultiWriter(os.Stdout, logFile)
}

// init log
func LoggerInit(serviceName string, skip *int) {
	ServiceName = serviceName

	if skip != nil {
		Skip = *skip
	}

	if err := os.MkdirAll("./logs", 0o755); err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	mw := InitLumberJack()
	GlobalLogger = log.New(mw, "", 0)
	log.Printf("[logger] initialized for service: %s (rotating logs, max 10MB x5, 30d)\n", serviceName)
}

func WithRequestID(ctx context.Context) context.Context {
	id := uuid.New().String()
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

func WithRequestIDFromHeader(ctx context.Context, header string) context.Context {
	if header == "" {
		return WithRequestID(ctx)
	}
	return context.WithValue(ctx, requestIDKey, header)
}

func Print(ctx context.Context, level LogLevel, msg string, data ...interface{}) {
	file, function := getCallerInfo()
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   ServiceName,
		Level:     level,
		Message:   fmt.Sprintf(msg, data...),
		File:      file,
		Function:  function,
		RequestID: GetRequestID(ctx),
	}

	dataEntry, err := json.Marshal(entry)
	if err != nil {
		log.Printf("failed to marshal log: %v", err)
		return
	}

	fmt.Fprintln(GlobalLogger.Writer(), string(dataEntry))
}

func LoggerInfo(ctx context.Context, msg string, fields ...interface{}) {
	Print(ctx, LevelInfo, msg, fields...)
}

func LoggerWarn(ctx context.Context, msg string, fields ...interface{}) {
	Print(ctx, LevelWarn, msg, fields...)
}

func LoggerError(ctx context.Context, msg string, fields ...interface{}) {
	Print(ctx, LevelError, msg, fields...)
}

func LoggerInfof(ctx context.Context, msg string, fields ...interface{}) {
	Print(ctx, LevelInfo, msg, fields...)
}

func LoggerWarnf(ctx context.Context, msg string, fields ...interface{}) {
	Print(ctx, LevelWarn, msg, fields...)
}

func LoggerErrorf(ctx context.Context, fields ...interface{}) {
	logWithLevel(ctx, LevelError, fields...)
}

func logWithLevel(ctx context.Context, level LogLevel, fields ...interface{}) {
	msg := formatMessage(fields...)
	Print(ctx, level, msg)
}

func formatMessage(v ...interface{}) string {
	if len(v) == 1 {
		val := v[0]
		switch reflect.TypeOf(val).Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice:
			b, _ := json.MarshalIndent(val, "", "  ")
			return string(b)
		}
	}
	return fmt.Sprint(v...)
}

// gorm logger
type GormJSONLogger struct {
	serviceName string
	writer      io.Writer
	level       gormLogger.LogLevel
}

func NewGormLogger(serviceName string) gormLogger.Interface {
	mw := InitLumberJack()
	return &GormJSONLogger{
		serviceName: serviceName,
		writer:      mw,
		level:       gormLogger.Info,
	}
}

func getCallerInfo() (file string, function string) {
	pc, filename, line, ok := runtime.Caller(Skip)
	if !ok {
		return "unknown", "unknown"
	}
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
	}
	shortFile := fmt.Sprintf("%s:%d", filepath.Base(filename), line)
	return shortFile, funcName
}

func (l *GormJSONLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *GormJSONLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.writeLog(ctx, "info", msg, data...)
}

func (l *GormJSONLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.writeLog(ctx, "warn", msg, data...)
}

func (l *GormJSONLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.writeLog(ctx, "error", msg, data...)
}

func (l *GormJSONLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)

	message := fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)
	var level LogLevel = "info"
	if err != nil {
		level = "error"
		message = fmt.Sprintf("%s | error: %v", message, err)
	}

	l.writeLog(ctx, level, message, nil)
}

// writeLog writes a JSON-formatted log entry.
func (l *GormJSONLogger) writeLog(ctx context.Context, level LogLevel, msg string, data ...interface{}) {
	file, function := getCallerInfo()

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   l.serviceName,
		Level:     level,
		Message:   fmt.Sprintf(msg, data...),
		File:      file,
		Function:  function,
		RequestID: GetRequestID(ctx),
	}

	_ = json.NewEncoder(l.writer).Encode(entry)
}
