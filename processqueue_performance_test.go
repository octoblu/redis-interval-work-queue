package main_test

import (
  "github.com/garyburd/redigo/redis"
	. "github.com/octoblu/redis-interval-work-queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProcessqueuePerformance", func() {
	var redisConn redis.Conn

	BeforeEach(func(){
		var err error
		redisConn,err = redis.Dial("tcp", ":6379")
		Expect(err).To(BeNil())

		redisConn.Do("DEL", "circular-job-queue")
		redisConn.Do("DEL", "linear-job-queue")
		redisConn.Do("LPUSH", "circular-job-queue", "happenstance", "huh")
	})

	AfterEach(func(){
		redisConn.Close()
	})

	Measure("one time", func(b Benchmarker) {
		queue := NewProcessQueue()
		runtime := b.Time("runtime", func() {
			queue.Process()
		})

		Expect(runtime.Seconds()).To(BeNumerically("<", 0.2), "queue.Process()'s execution time is too darn high!'")
	}, 10)

	Measure("one hundred times", func(b Benchmarker) {
		queue := NewProcessQueue()
		runtime := b.Time("runtime", func() {
			oneHundred := 100
			for i := 0; i < oneHundred; i++ {
				queue.Process()
			}
		})

		Expect(runtime.Seconds()).To(BeNumerically("<", 1), "queue.Process()'s execution time is too darn high!'")
	}, 10)

	Measure("one thousand times", func(b Benchmarker) {
		queue := NewProcessQueue()
		runtime := b.Time("runtime", func() {
			oneThousand := 1000
			for i := 0; i < oneThousand; i++ {
				queue.Process()
			}
		})

		Expect(runtime.Seconds()).To(BeNumerically("<", 1), "queue.Process()'s execution time is too darn high!'")
	}, 10)

	Measure("ten thousand times", func(b Benchmarker) {
		queue := NewProcessQueue()
		runtime := b.Time("runtime", func() {
			tenThousand := 10 * 1000
			for i := 0; i < tenThousand; i++ {
				queue.Process()
			}
		})

		Expect(runtime.Seconds()).To(BeNumerically("<", 1), "queue.Process()'s execution time is too darn high!'")
	}, 10)
})
