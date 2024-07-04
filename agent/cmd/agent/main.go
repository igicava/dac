package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"dac/agent/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	LoadEnv()
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Printf("Invalid COMPUTING_POWER: %v", err)
		computingPower = 4
	}

	var wg sync.WaitGroup

	fmt.Println("Agents is runing")
	// Запуск горутин
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
					time.Sleep(2 * time.Second)
				} 
			}
		}()
	}
	wg.Wait()
}

// Загрузка переменных среды. Для предотвращения ошибок на этом этапе запускайте проект так как написано в README
func LoadEnv() {
	envPath := filepath.Join(".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}
