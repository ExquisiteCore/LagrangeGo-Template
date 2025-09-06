package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/LagrangeDev/LagrangeGo/utils"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

// LogLevel 日志级别
type LogLevel int

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// LogConfig 日志配置
type LogConfig struct {
	Level        LogLevel
	EnableFile   bool
	EnableColor  bool
	LogDir       string
	LogFile      string
	MaxSize      int64  // MB
	MaxBackups   int
	MaxAge       int    // days
	Format       string // json, text
}

// DefaultLogConfig 默认日志配置
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:       InfoLevel,
		EnableFile:  true,
		EnableColor: true,
		LogDir:      "logs",
		LogFile:     "bot.log",
		MaxSize:     10,
		MaxBackups:  5,
		MaxAge:      30,
		Format:      "text",
	}
}

// Logger 统一日志接口
type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// BotLogger 机器人日志器
type BotLogger struct {
	*logrus.Logger
	config *LogConfig
}

// NewBotLogger 创建新的机器人日志器
func NewBotLogger(config *LogConfig) *BotLogger {
	if config == nil {
		config = DefaultLogConfig()
	}
	
	logger := logrus.New()
	
	// 设置日志级别 (修复等级反向问题)
	logger.SetLevel(convertToLogrusLevel(config.Level))
	
	// 设置格式化器
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		if config.EnableColor {
			logger.SetFormatter(&ColoredFormatter{})
		} else {
			logger.SetFormatter(&PlainFormatter{})
		}
	}
	
	// 设置输出 (修复彩色输出和文件输出冲突)
	var outputs []io.Writer
	
	// 总是添加控制台输出
	if config.EnableColor {
		outputs = append(outputs, colorable.NewColorableStdout())
	} else {
		outputs = append(outputs, os.Stdout)
	}
	
	// 如果启用文件输出，添加文件输出
	if config.EnableFile {
		// 创建日志目录
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			panic(fmt.Sprintf("创建日志目录失败: %v", err))
		}
		
		// 创建日志文件
		logFile := filepath.Join(config.LogDir, config.LogFile)
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("打开日志文件失败: %v", err))
		}
		outputs = append(outputs, file)
	}
	
	// 使用多重输出
	logger.SetOutput(io.MultiWriter(outputs...))
	
	return &BotLogger{
		Logger: logger,
		config: config,
	}
}

// WithField 添加字段
func (l *BotLogger) WithField(key string, value interface{}) Logger {
	return &BotLogger{
		Logger: l.Logger.WithField(key, value).Logger,
		config: l.config,
	}
}

// WithFields 添加多个字段
func (l *BotLogger) WithFields(fields map[string]interface{}) Logger {
	return &BotLogger{
		Logger: l.Logger.WithFields(fields).Logger,
		config: l.config,
	}
}

// ProtocolLogger 协议日志器
type ProtocolLogger struct {
	logger Logger
}

// NewProtocolLogger 创建协议日志器
func NewProtocolLogger(logger Logger) *ProtocolLogger {
	return &ProtocolLogger{logger: logger}
}

const fromProtocol = "Lgr -> "

func (p *ProtocolLogger) Info(format string, arg ...any) {
	p.logger.Infof(fromProtocol+format, arg...)
}

func (p *ProtocolLogger) Warning(format string, arg ...any) {
	p.logger.Warnf(fromProtocol+format, arg...)
}

func (p *ProtocolLogger) Debug(format string, arg ...any) {
	p.logger.Debugf(fromProtocol+format, arg...)
}

func (p *ProtocolLogger) Error(format string, arg ...any) {
	p.logger.Errorf(fromProtocol+format, arg...)
}

func (p *ProtocolLogger) Dump(data []byte, format string, arg ...any) {
	message := fmt.Sprintf(format, arg...)
	dumpDir := "dump"
	
	if _, err := os.Stat(dumpDir); err != nil {
		err = os.MkdirAll(dumpDir, 0o755)
		if err != nil {
			p.logger.Errorf("出现错误 %v. 详细信息转储失败", message)
			return
		}
	}
	
	dumpFile := path.Join(dumpDir, fmt.Sprintf("%v.dump", time.Now().Unix()))
	p.logger.Errorf("出现错误 %v. 详细信息已转储至文件 %v 请连同日志提交给开发者处理", message, dumpFile)
	_ = os.WriteFile(dumpFile, data, 0o644)
}

// ColoredFormatter 彩色格式化器
type ColoredFormatter struct{}

const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorYellow = "\x1b[33m"
	colorGreen  = "\x1b[32m"
	colorBlue   = "\x1b[34m"
	colorWhite  = "\x1b[37m"
	colorCyan   = "\x1b[36m"
)

func (f *ColoredFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	
	var levelColor string
	switch entry.Level {
	case logrus.TraceLevel:
		levelColor = colorCyan
	case logrus.DebugLevel:
		levelColor = colorBlue
	case logrus.InfoLevel:
		levelColor = colorGreen
	case logrus.WarnLevel:
		levelColor = colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = colorRed
	default:
		levelColor = colorWhite
	}
	
	// 格式化字段
	var fieldsStr string
	if len(entry.Data) > 0 {
		fields := make([]string, 0, len(entry.Data))
		for key, value := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", key, value))
		}
		fieldsStr = fmt.Sprintf(" [%s]", strings.Join(fields, " "))
	}
	
	return utils.S2B(fmt.Sprintf("[%s] [%s%s%s]%s: %s\n",
		timestamp, levelColor, strings.ToUpper(entry.Level.String()), colorReset, fieldsStr, entry.Message)), nil
}

// PlainFormatter 普通格式化器 (不带颜色)
type PlainFormatter struct{}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	
	// 格式化字段
	var fieldsStr string
	if len(entry.Data) > 0 {
		fields := make([]string, 0, len(entry.Data))
		for key, value := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", key, value))
		}
		fieldsStr = fmt.Sprintf(" [%s]", strings.Join(fields, " "))
	}
	
	return utils.S2B(fmt.Sprintf("[%s] [%s]%s: %s\n",
		timestamp, strings.ToUpper(entry.Level.String()), fieldsStr, entry.Message)), nil
}

// 全局日志器
var (
	globalLogger Logger
	protocolLogger *ProtocolLogger
)

// Init 初始化日志系统
func Init() {
	config := DefaultLogConfig()
	globalLogger = NewBotLogger(config)
	protocolLogger = NewProtocolLogger(globalLogger)
}

// InitWithConfig 使用配置初始化日志系统
func InitWithConfig(config *LogConfig) {
	globalLogger = NewBotLogger(config)
	protocolLogger = NewProtocolLogger(globalLogger)
}

// GetLogger 获取全局日志器
func GetLogger() Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

// GetProtocolLogger 获取协议日志器
func GetProtocolLogger() *ProtocolLogger {
	if protocolLogger == nil {
		Init()
	}
	return protocolLogger
}

// 便捷方法
func Trace(args ...interface{}) {
	GetLogger().Trace(args...)
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	GetLogger().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

func WithField(key string, value interface{}) Logger {
	return GetLogger().WithField(key, value)
}

func WithFields(fields map[string]interface{}) Logger {
	return GetLogger().WithFields(fields)
}

// convertToLogrusLevel 转换日志等级
func convertToLogrusLevel(level LogLevel) logrus.Level {
	switch level {
	case TraceLevel:
		return logrus.TraceLevel
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	case PanicLevel:
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
