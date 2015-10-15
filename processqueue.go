package main

import (
	"github.com/octoblu/circularqueue"
	"github.com/octoblu/claimablejob"
)

// A ProcessQueue processes the circular work queue and dump jobs into the linear work queue
type ProcessQueue interface {
	Process() error
}

// A RedisProcessQueue processes the circular work queue and dump jobs into the linear work queue
type RedisProcessQueue struct {
}

// NewProcessQueue constructs a new Redis Process Queue instance
func NewProcessQueue() *RedisProcessQueue {
	return new(RedisProcessQueue)
}

// Process processes the circular work queue and dump jobs into the linear work queue
func (redisProcessQueue *RedisProcessQueue) Process() error {
	queue := circularqueue.New()
	job,err := queue.Pop()
	if err != nil {
		return err
	}

  claimableJob := claimablejob.NewFromJob(job)

  if claimed, err := claimableJob.Claim(); err != nil {
    return err
  } else if !claimed {
    return nil
  }

  claimableJob.PushKeyIntoQueue("linear-job-queue")
	return nil
}
