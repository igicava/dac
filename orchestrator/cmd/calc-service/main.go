package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"dac/orchestrator/http/server"
	pb "dac/proto"
	grpc_server "dac/orchestrator/grpc/server"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	LoadEnv()
	fmt.Println("Server is start on port 8080")
	go server.Run()
	fmt.Println("gRPC server runing...")
	host := "localhost"
	port := "8081"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части CalcService
	calcServiceServer := grpc_server.NewServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterCalcServiceServer(grpcServer, calcServiceServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

// Загрузка переменных среды. Для предотвращения ошибок на этом этапе запускайте проект так как написано в README
func LoadEnv() {
	envPath := filepath.Join(".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}
