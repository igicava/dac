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

    r.HandleFunc("/api/v1/register", handler.RegisterNewUser)
    r.HandleFunc("/api/v1/login", handler.LoginUser)

    log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}