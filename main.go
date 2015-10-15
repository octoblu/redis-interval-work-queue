package main

import (
  "log"
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "redis-interval-work-queue"
  app.Usage = "Run the work queue"
  app.Action = processQueue
  app.Run(os.Args)
}

func processQueue(context *cli.Context) {
  queue := NewProcessQueue()
  for {
    if err := queue.Process(); err != nil {
      log.Fatalf("Error occured: %v", err.Error())
    }
  }
}
