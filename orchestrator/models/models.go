package models

import (
	"sync"
)

var Expressions = make(map[string]Expression)
var Tasks = make(chan Task)
var Mu sync.Mutex

type Expression struct {
    ID         string  `json:"id"`
    Expression string  `json:"expression"`
    Status     string  `json:"status"`
    Result     float64 `json:"result"`
}

type Task struct {
    ID            string  `json:"id"`
    Arg1          float64 `json:"arg1"`
    Arg2          float64 `json:"arg2"`
    Operation     string  `json:"operation"`
    OperationTime int     `json:"operation_time"`
}