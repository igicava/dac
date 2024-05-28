package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"dac/agent/internal/app"
)

func main() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER: %v", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				task := app.GetTask()
				if task != nil {
					log.Printf("Received task: %+v", task)
					result := app.ComputeTask(*task)
					log.Printf("Computed result for task %s: %f", task.ID, result)
					app.SendResult(task.ID, result)
				} else {
					time.Sleep(3 * time.Second) // Интервал ожидания для предотвращения частого опроса
				}
			}
		}()
	}
	wg.Wait()
}
