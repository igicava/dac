package server

import (
	"net/http"
	"log"
    
	"dac/orchestrator/http/server/handler"
)

func Run() {
    r := http.NewServeMux()
    
    r.HandleFunc("/api/v1/calculate", handler.AddExpression)
    r.HandleFunc("/api/v1/expressions", handler.GetExpressions)
    r.HandleFunc("/api/v1/expressions/{id}", handler.GetExpressionByID)

    log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}