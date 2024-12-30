package requestidmanager

import (
	"sync"
	"time"
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64 //时间戳
	workerID  int64 //机器ID
	sequence  int64 //序列号
}

/*
雪花算法
符号位（1位）：

总是 0，因为 ID 是正数。
时间戳（41位）：

通常是从某个固定时间开始的毫秒数，这部分提供了 69 年的时间使用范围。
数据中心 ID（5位）：

用于区分不同的数据中心，最多支持 32 个。
机器 ID（5位）：

用于区分同一数据中心内的不同机器，最多支持 32 台。
序列号（12位）：

用于在同一毫秒内生成多个 ID，支持每毫秒 4096 个不同的 ID。
*/
const (
	workerIDBits   = 5                           // 机器 ID 占用的位数
	sequenceBits   = 12                          // 序列号占用的位数
	workerIDShift  = sequenceBits                // 机器 ID 左移位数
	timestampShift = sequenceBits + workerIDBits // 时间戳左移位数
	sequenceMask   = -1 ^ (-1 << sequenceBits)   // 序列号掩码
	epoch          = int64(1609459200000)        // 自定义起始时间戳，2021-01-01 00:00:00 UTC
)

func NewSnowflake(workerID int64) *Snowflake { // 创建一个新的雪花算法 workerID用于标识生成 ID 的节点
	return &Snowflake{
		workerID: workerID,
	}
}

func (s *Snowflake) NextID() int64 { // 生成一个新的 ID
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now
	id := ((now - epoch) << timestampShift) | (s.workerID << workerIDShift) | s.sequence
	return id
}
