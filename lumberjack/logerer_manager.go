package lumberjack

import (
	"strings"
	"sync"
	"time"
)

// 定义一个结构体用于存储日志目录信息
type LogInfo struct {
	CurDate  string
	Filename string
}

// LoggerDateManager 用于管理日志信息的映射
type LoggerDateManager struct {
	Logger map[string]LogInfo
	mu     sync.RWMutex // 读写锁
}

// 单例实例和 sync.Once
var (
	date_instance *LoggerDateManager
	date_once     sync.Once
)

// GetInstance 返回 LoggerDateManager 的单例实例
func GetInstance() *LoggerDateManager {
	date_once.Do(func() {
		date_instance = &LoggerDateManager{
			Logger: make(map[string]LogInfo),
		}
	})
	return date_instance
}

// UpdateLogDateInfo 日志信息
func (manager *LoggerDateManager) UpdateLogDateInfo(filepath string) {
	now := time.Now().Format("2006-01-02")

	// 如果找到了日志信息，
	manager.mu.RLock() // 加读锁
	value, exists := manager.Logger[filepath]
	manager.mu.RUnlock() // 释放读锁

	if exists { // 如果找到了日志信息，检查日期是否相同
		if now != value.CurDate { // 如果日期不同，更新日志信息
			filename := strings.ReplaceAll(filepath, "{DATE}", now)
			value.CurDate = now
			value.Filename = filename

			manager.mu.Lock() // 加写锁
			manager.Logger[filepath] = value
			manager.mu.Unlock() // 释放写锁
		}
	} else { // 如果没有找到日志信息，添加新的日志信息
		filename := strings.ReplaceAll(filepath, "{DATE}", now)
		manager.mu.Lock() // 加写锁
		manager.Logger[filepath] = LogInfo{
			CurDate:  now,
			Filename: filename,
		}
		manager.mu.Unlock() // 释放写锁
	}
}

// GetLogInfo 获取日志信息
func (manager *LoggerDateManager) GetLogDateInfo(filepath string) string {
	now := time.Now().Format("2006-01-02")

	// 如果找到了日志信息，
	manager.mu.RLock() // 加读锁
	value, exists := manager.Logger[filepath]
	manager.mu.RUnlock() // 释放读锁

	if exists { // 如果找到了日志信息，检查日期是否相同
		if now != value.CurDate { // 如果日期不同，更新日志信息
			filename := strings.ReplaceAll(filepath, "{DATE}", now)
			value.CurDate = now
			value.Filename = filename

			manager.mu.Lock() // 加写锁
			manager.Logger[filepath] = value
			manager.mu.Unlock() // 释放写锁

			return filename
		} else {
			return value.Filename
		}
	} else { // 如果没有找到日志信息，添加新的日志信息
		filename := strings.ReplaceAll(filepath, "{DATE}", now)

		manager.mu.Lock()         // 加写锁
		defer manager.mu.Unlock() // 释放写锁

		value, exists := manager.Logger[filepath]
		if exists {
			return value.Filename
		}

		manager.Logger[filepath] = LogInfo{
			CurDate:  now,
			Filename: filename,
		}
		return filename
	}
}
