package logging

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    "runtime"
    "time"
)

// LogLevel định nghĩa các cấp độ log
type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

// String trả về tên của cấp độ log
func (l LogLevel) String() string {
    return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// Logger cung cấp chức năng ghi log có cấu trúc
type Logger struct {
    Level  LogLevel
    Writer io.Writer
}

// LogEntry đại diện cho một bản ghi log
type LogEntry struct {
    Timestamp string      `json:"timestamp"`
    Level     string      `json:"level"`
    Message   string      `json:"message"`
    File      string      `json:"file,omitempty"`
    Line      int         `json:"line,omitempty"`
    Data      interface{} `json:"data,omitempty"`
}

// NewLogger tạo một logger mới
func NewLogger(level LogLevel) *Logger {
    return &Logger{
        Level:  level,
        Writer: os.Stdout,
    }
}

// SetOutput thiết lập writer cho logger
func (l *Logger) SetOutput(w io.Writer) {
    l.Writer = w
}

// log ghi một bản ghi log
func (l *Logger) log(level LogLevel, msg string, data interface{}) {
    if level < l.Level {
        return
    }

    // Lấy thông tin file và line number
    _, file, line, ok := runtime.Caller(2)
    if !ok {
        file = "unknown"
        line = 0
    }

    entry := LogEntry{
        Timestamp: time.Now().Format(time.RFC3339),
        Level:     level.String(),
        Message:   msg,
        File:      file,
        Line:      line,
        Data:      data,
    }

    jsonData, err := json.Marshal(entry)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
        return
    }

    fmt.Fprintln(l.Writer, string(jsonData))
    
    // Thoát chương trình nếu là log FATAL
    if level == FATAL {
        os.Exit(1)
    }
}

// Debug ghi log ở cấp độ DEBUG
func (l *Logger) Debug(msg string, data ...interface{}) {
    var logData interface{}
    if len(data) > 0 {
        logData = data[0]
    }
    l.log(DEBUG, msg, logData)
}

// Info ghi log ở cấp độ INFO
func (l *Logger) Info(msg string, data ...interface{}) {
    var logData interface{}
    if len(data) > 0 {
        logData = data[0]
    }
    l.log(INFO, msg, logData)
}

// Warn ghi log ở cấp độ WARN
func (l *Logger) Warn(msg string, data ...interface{}) {
    var logData interface{}
    if len(data) > 0 {
        logData = data[0]
    }
    l.log(WARN, msg, logData)
}

// Error ghi log ở cấp độ ERROR
func (l *Logger) Error(msg string, data ...interface{}) {
    var logData interface{}
    if len(data) > 0 {
        logData = data[0]
    }
    l.log(ERROR, msg, logData)
}

// Fatal ghi log ở cấp độ FATAL và thoát chương trình
func (l *Logger) Fatal(msg string, data ...interface{}) {
    var logData interface{}
    if len(data) > 0 {
        logData = data[0]
    }
    l.log(FATAL, msg, logData)
}

// Biến logger toàn cục
var DefaultLogger = NewLogger(INFO)

// Các hàm log toàn cục
func Debug(msg string, data ...interface{}) {
    DefaultLogger.Debug(msg, data...)
}

func Info(msg string, data ...interface{}) {
    DefaultLogger.Info(msg, data...)
}

func Warn(msg string, data ...interface{}) {
    DefaultLogger.Warn(msg, data...)
}

func Error(msg string, data ...interface{}) {
    DefaultLogger.Error(msg, data...)
}

func Fatal(msg string, data ...interface{}) {
    DefaultLogger.Fatal(msg, data...)
}

// Khởi tạo logger
func init() {
    // Thiết lập cấp độ log từ biến môi trường
    logLevel := os.Getenv("LOG_LEVEL")
    switch logLevel {
    case "DEBUG":
        DefaultLogger.Level = DEBUG
    case "INFO":
        DefaultLogger.Level = INFO
    case "WARN":
        DefaultLogger.Level = WARN
    case "ERROR":
        DefaultLogger.Level = ERROR
    case "FATAL":
        DefaultLogger.Level = FATAL
    }
}