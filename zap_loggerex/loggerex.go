package zap_loggerex

import (
	"fmt"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/blueturbo-ad/go-utils/config_manage"

	"github.com/VarusHsu/lumberjack"
)

/* LoggerWrapper 日志包装器
 * ZapLogger 用于表示zap日志
 * LoggerManager 日志管理器 双缓存管理器
 * Logger 日志

 */

type LoggerWrapper struct {
	ZapLogger *zap.Logger
}

type Logger struct {
	Logger  *LoggerWrapper
	version string
}

type LoggerEx struct {
	LoggersMap map[string]*Logger
}

// 双缓存管理器
type LoggerManager struct {
	Config  *config_manage.ZapLoggerConfig
	current *LoggerEx
	next    *LoggerEx
	rwMutex sync.RWMutex
}

var (
	instance    *LoggerManager
	once        sync.Once
	EmptyString = ""
)

func GetSingleton() *LoggerManager {
	once.Do(func() {
		instance = new(LoggerManager)
		instance.current = new(LoggerEx)
		instance.next = new(LoggerEx)
	})
	return instance
}

func GetLogger() *LoggerManager {
	return GetSingleton()
}

// 函数用于内存更新etcd配置
func (l *LoggerManager) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "logger":
		var e = new(config_manage.ZapLoggerConfig)
		err = e.LoadMemoryZapConfig([]byte(value), env)
		if err != nil {
			return err
		}
		return l.UpdateLogger(e)
	default:
		return fmt.Errorf("unknown UpdateFromEtcd: %s", key)
	}
}

func (l *LoggerManager) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.ZapLoggerConfig)
	err = e.LoadZapConfig(confPath, env)
	if err != nil {
		return err
	}

	return l.UpdateLogger(e)
}

func (l *LoggerManager) UpdateLogger(config *config_manage.ZapLoggerConfig) error {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()
	var loger = new(LoggerEx) //生成新的数据
	for _, value := range config.Loggers {
		zapLogger := newZapLogger(&value)
		if zapLogger == nil {
			return nil
		}

		if loger.LoggersMap == nil {
			loger.LoggersMap = make(map[string]*Logger)
		}

		loger.LoggersMap[value.Name] = &Logger{&LoggerWrapper{zapLogger}, config.Version}
	}
	l.next = loger
	l.current, l.next = l.next, l.current

	return nil
}

func newZapLogger(conf *config_manage.LoggerConfig) *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.WarnLevel && lev >= zapcore.Level(conf.Level)
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.WarnLevel && lev >= zap.DebugLevel && lev >= zapcore.Level(conf.Level)
	})

	prodEncoder := createLogOutputEncoderConfig()
	prodEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	lowWriteSyncer := createWriteSyncer(conf, true)
	highWriteSyncer := createWriteSyncer(conf, false)

	highCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), highWriteSyncer, highPriority)
	lowCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), lowWriteSyncer, lowPriority)

	return zap.New(zapcore.NewTee(highCore, lowCore), zap.AddCaller(), zap.AddCallerSkip(2))
}

func createWriteSyncer(conf *config_manage.LoggerConfig, isinfo bool) zapcore.WriteSyncer {
	var info string
	if isinfo {
		info = conf.Info
	} else {
		info = conf.Error
	}
	if len(info) == 0 {
		log.Printf("LoggerEx logger path length is 0")
	}
	var hookFunc func(string) = nil

	lumberJackLogger := &lumberjack.Logger{
		Filename:   info,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
		Hook:       hookFunc,
	}
	if conf.Async {
		syncWriter := &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(lumberJackLogger),
			Size:          4096,
			FlushInterval: 1 * time.Second,
		}
		return syncWriter
	}

	return zapcore.AddSync(lumberJackLogger)
}

func createLogOutputEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:          "ts",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller_line",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       "\n",
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.EpochTimeEncoder,
		ConsoleSeparator: "\n",
		EncodeCaller:     zapcore.ShortCallerEncoder,
	}
}

func (l *LoggerManager) getLogger(name string) (*LoggerWrapper, error) {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()
	if l.current == nil {
		return nil, fmt.Errorf("LoggerEx index is -1")
	}
	if l.current.LoggersMap == nil {
		return nil, fmt.Errorf("LoggerEx LoggersMap is nil")
	}
	if l.current.LoggersMap[name] == nil {
		return nil, fmt.Errorf("LoggerEx LoggersMap[%s] is nil", name)
	}
	if l.current.LoggersMap[name].Logger == nil {
		return nil, fmt.Errorf("LoggerEx LoggersMap[%s].Logger is nil", name)
	}

	return l.current.LoggersMap[name].Logger, nil
}

func (l *LoggerManager) Debug(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx Debug getLogger error:%v", err)
		return err
	}
	logger.Debug(format, a...)
	return nil
}

func (l *LoggerManager) Info(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx Info getLogger error:%v", err)
		return err
	}
	logger.Info(format, a...)
	return nil
}

func (l *LoggerManager) Warn(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx Warn getLogger error:%v", err)
		return err
	}
	logger.Warn(format, a...)
	return nil
}

func (l *LoggerManager) Error(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx Error getLogger error:%v", err)
		return err
	}
	logger.Error(format, a...)
	return nil
}

func (l *LoggerManager) DPanic(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx DPanic getLogger error:%v", err)
		return err
	}
	logger.DPanic(format, a...)
	return nil
}

func (l *LoggerManager) Panic(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx Panic getLogger error:%v", err)
		return err
	}
	logger.Panic(format, a...)
	return nil
}

func (l *LoggerManager) Fatal(name string, format string, a ...any) error {
	logger, err := l.getLogger(name)
	if err != nil {
		log.Printf("LoggerEx Fatal getLogger error:%v", err)
		return err
	}
	logger.Fatal(format, a...)
	return nil
}

func (l *LoggerWrapper) Debug(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.DebugLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.Debug(msg)
}

func (l *LoggerWrapper) Info(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.InfoLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.Info(msg)
}

func (l *LoggerWrapper) Warn(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.WarnLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.Warn(msg)
}

func (l *LoggerWrapper) Error(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.WarnLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.Error(msg)
}

func (l *LoggerWrapper) DPanic(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.DPanicLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.DPanic(msg)
}

func (l *LoggerWrapper) Panic(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.PanicLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.Panic(msg)
}

func (l *LoggerWrapper) Fatal(format string, fields ...any) {
	checkedEntry := l.ZapLogger.Check(zapcore.FatalLevel, EmptyString)
	if checkedEntry == nil {
		return
	}
	msg := fmt.Sprintf(format, fields...)
	l.ZapLogger.Fatal(msg)
}
