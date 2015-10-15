package main

import (
  "time"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"github.com/octoblu/circularqueue"
)

// A ProcessQueue processes the circular work queue and dump jobs into the linear work queue
type ProcessQueue interface {
	Process() error
}

// A RedisProcessQueue processes the circular work queue and dump jobs into the linear work queue
type RedisProcessQueue struct {
}

// Process processes the circular work queue and dump jobs into the linear work queue
func (redisProcessQueue *RedisProcessQueue) Process() error {
	queue := circularqueue.New()
	job,err := queue.Pop()
	if err != nil {
		return err
	}

	redisConn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	result,err := redisConn.Do("GET", "[namespace]-parachute-failure")
	if err != nil {
		return err
	}

	nextTick := parseNextTick(result)
  now := time.Now().Unix()
	if now < nextTick {
		return nil
	}

	redisConn.Do("RPUSH", "linear-job-queue", job.GetKey())

	return nil
}

// NewProcessQueue constructs a new Redis Process Queue instance
func NewProcessQueue() *RedisProcessQueue {
	return new(RedisProcessQueue)
}

func parseNextTick(redisResult interface{}) int64 {
  strNextTick,ok := redisResult.([]uint8)
  if !ok {
    return 0
  }

  nextTick,err := strconv.ParseInt(string(strNextTick), 10, 64)
  if err != nil {
    return 0
  }

  return nextTick
}
