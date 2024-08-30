package pkg

import (
	"fmt"
	"time"
)

const (
	// 起始时间戳（2020-01-01）
	twepoch          int64 = 1577836800000
	workerIdBits     uint8 = 5
	datacenterIdBits uint8 = 5
	sequenceBits     uint8 = 12

	maxWorkerId     int64 = -1 ^ (-1 << workerIdBits)
	maxDatacenterId int64 = -1 ^ (-1 << datacenterIdBits)
	sequenceMask    int64 = -1 ^ (-1 << sequenceBits)

	workerIdShift      uint8 = sequenceBits
	datacenterIdShift  uint8 = sequenceBits + workerIdBits
	timestampLeftShift uint8 = sequenceBits + workerIdBits + datacenterIdBits
)

type Snowflake struct {
	lastTimestamp int64
	workerId      int64
	datacenterId  int64
	sequence      int64
}

func NewSnowflake(workerId, datacenterId int64) *Snowflake {
	if workerId < 0 || workerId > maxWorkerId {
		panic(fmt.Sprintf("worker Id can't be greater than %d or less than 0", maxWorkerId))
	}
	if datacenterId < 0 || datacenterId > maxDatacenterId {
		panic(fmt.Sprintf("datacenter Id can't be greater than %d or less than 0", maxDatacenterId))
	}
	return &Snowflake{
		workerId:      workerId,
		datacenterId:  datacenterId,
		lastTimestamp: -1,
		sequence:      0,
	}
}

func (s *Snowflake) NextId() int64 {
	timestamp := time.Now().UnixNano() / 1e6
	if timestamp < s.lastTimestamp {
		panic(fmt.Sprintf("Clock moved backwards. Refusing to generate id for %d milliseconds", s.lastTimestamp-timestamp))
	}
	if s.lastTimestamp == timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = timestamp
	return ((timestamp - twepoch) << timestampLeftShift) |
		(s.datacenterId << datacenterIdShift) |
		(s.workerId << workerIdShift) |
		s.sequence
}
