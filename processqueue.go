package main

import (
  "github.com/garyburd/redigo/redis"
	"github.com/octoblu/circularqueue"
	"github.com/octoblu/claimablejob"
)

// A ProcessQueue processes the circular work queue and dump jobs into the linear work queue
type ProcessQueue interface {
	Process() error
}

// A RedisProcessQueue processes the circular work queue and dump jobs into the linear work queue
type RedisProcessQueue struct {
	_conn redis.Conn
}

// NewProcessQueue constructs a new Redis Process Queue instance
func NewProcessQueue() *RedisProcessQueue {
	return new(RedisProcessQueue)
}

// Process processes the circular work queue and dump jobs into the linear work queue
func (redisProcessQueue *RedisProcessQueue) Process() error {
	conn,err := redisProcessQueue.conn()
	if err != nil {
		return err
	}

	queue := circularqueue.New(conn)
	job,err := queue.Pop()
	if err != nil {
		return err
	}

  claimableJob := claimablejob.NewFromJob(job, conn)

  if claimed, err := claimableJob.Claim(); err != nil {
    return err
  } else if !claimed {
    return nil
  }

  claimableJob.PushKeyIntoQueue("linear-job-queue")
	return nil
}

func (redisProcessQueue *RedisProcessQueue) conn() (redis.Conn,error) {
	if redisProcessQueue._conn != nil {
		return redisProcessQueue._conn, nil
	}

	conn,err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	redisProcessQueue._conn = conn
	return conn, nil
}
