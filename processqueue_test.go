package main_test

import (
	"time"
	"github.com/garyburd/redigo/redis"
	. "github.com/octoblu/redis-interval-work-queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProcessQueue", func() {
	var redisConn redis.Conn
	var sut ProcessQueue

	BeforeEach(func() {
		var err error

		redisConn, err = redis.Dial("tcp", ":6379")
		Expect(err).To(BeNil())

		sut = NewProcessQueue()
	})

	AfterEach(func(){
		redisConn.Close()
	})

	Describe("Process", func(){
		Context("When there are two jobs in the circular queue", func(){
			BeforeEach(func(){
				var err error
				_,err = redisConn.Do("DEL", "linear-job-queue")
				_,err = redisConn.Do("DEL", "circular-job-queue")
				_,err = redisConn.Do("DEL", "[namespace]-foolishly-ignore-warning")
				_,err = redisConn.Do("DEL", "[namespace]-nah-i-got-this")
				_,err = redisConn.Do("RPUSH", "circular-job-queue", "foolishly-ignore-warning", "nah-i-got-this")
				err = sut.Process()
				Expect(err).To(BeNil())
			})

			It("should cycle the queue", func(){
				resp,err := redisConn.Do("LINDEX", "circular-job-queue", 0)
				Expect(err).To(BeNil())

				record := string(resp.([]uint8))
				Expect(record).To(Equal("nah-i-got-this"))
			})

			It("should add to the other queue", func(){
				resp,err := redisConn.Do("LINDEX", "linear-job-queue", 0)
				Expect(err).To(BeNil())

				record := string(resp.([]uint8))
				Expect(record).To(Equal("nah-i-got-this"))
			})

			Context("when Process is called a second time", func(){
				BeforeEach(func(){
					err := sut.Process()
					Expect(err).To(BeNil())
				})

				It("should add to the other queue", func(){
					resp,err := redisConn.Do("LINDEX", "linear-job-queue", 0)
					Expect(err).To(BeNil())
					Expect(resp).NotTo(BeNil())

					record := string(resp.([]uint8))
					Expect(record).To(Equal("foolishly-ignore-warning"))
				})
			});
		})

		Context("when there is one job in the circular queue", func(){
			BeforeEach(func(){
				var err error
				_,err = redisConn.Do("DEL", "linear-job-queue")
				Expect(err).To(BeNil())
				_,err = redisConn.Do("DEL", "circular-job-queue")
				Expect(err).To(BeNil())
				_,err = redisConn.Do("RPUSH", "circular-job-queue", "parachute-failure")
				Expect(err).To(BeNil())
			})

			Context("When the job has already run this second", func(){
				BeforeEach(func(){
					then := int64(time.Now().Unix()	+ 1)
					_,err := redisConn.Do("SET", "[namespace]-parachute-failure", then)
					Expect(err).To(BeNil())

					err = sut.Process()
					Expect(err).To(BeNil())
				})

				AfterEach(func(){
					_,err := redisConn.Do("DEL", "[namespace]-parachute-failure")
					Expect(err).To(BeNil())
				})

				It("should not push the job into the linear queue", func(){
					resp,err := redisConn.Do("LLEN", "linear-job-queue")
					Expect(err).To(BeNil())
					Expect(resp).To(Equal(int64(0)))
				})
			})

			Context("When the job ran in the previous second", func(){
				BeforeEach(func(){
					then := int64(time.Now().Unix())
					_,err := redisConn.Do("SET", "[namespace]-parachute-failure", then)
					Expect(err).To(BeNil())

					err = sut.Process()
					Expect(err).To(BeNil())
				})

				AfterEach(func(){
					_,err := redisConn.Do("DEL", "[namespace]-parachute-failure")
					Expect(err).To(BeNil())
				})

				It("should push the job into the linear queue", func(){
					resp,err := redisConn.Do("LLEN", "linear-job-queue")
					Expect(err).To(BeNil())
					Expect(resp).To(Equal(int64(1)))
				})
			})
		})
	})
})
