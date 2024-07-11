package server

import (
	"context"
	"encoding/json"
	"log"

	"dac/orchestrator/models"
	pb "dac/proto"
)

type Server struct {
	pb.CalcServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GETtask(ctx context.Context, in *pb.GETRequest) (*pb.GETResponse, error) {
	task := <-models.Tasks
	byteTask, err := json.Marshal(task)
	if err != nil {
		log.Println("Error Marshal Task")
	}
	return &pb.GETResponse{
		Result: byteTask,
	}, nil
}

func (s *Server) POSTtask(ctx context.Context, in *pb.POSTRequest) (*pb.POSTResponse, error) {
	var taskResult struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	err := json.Unmarshal([]byte(in.JsonTASK), &taskResult)
	if err != nil {
		log.Println("Error Unmarshaling JSON")
	}

	models.Results[taskResult.ID] <- taskResult

	models.Mu.Lock()
	defer models.Mu.Unlock()

	err = models.UpdateResult(taskResult.ID, taskResult.Result)
	if err != nil {
		log.Printf("Error grpc POSTtask 49 : %s", err)
	}

	return &pb.POSTResponse{}, nil
}
